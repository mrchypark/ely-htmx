import {mount,emit} from '/assets/ely-api.js';
import {assert,waitFor} from './advanced-checks.js';
import {actionSummary} from './action-summary.js';
function checkControlSizes(){
 const fixture=document.createElement('section');document.body.append(fixture);
 try{
  for(const composite of [false,true]){
   if(composite)fixture.dataset.uiComponent='sizing-fixture';
   for(const name of ['input','select','button','textarea']){
    const template=document.querySelector(`[data-ui-control-template="${name}"]`);
    assert(template,`Missing shared ${name} template`);
    const node=template.content.firstElementChild.cloneNode(true);fixture.append(node);
    assert(getComputedStyle(node).fontSize==='13px',`${name} typography differs between standalone and composite usage`);
    if(name!=='textarea')assert(node.getBoundingClientRect().height===28,`${name} must use the 28px default height`);
    if(name==='input')assert(getComputedStyle(node).paddingTop==='3px'&&getComputedStyle(node).paddingBottom==='3px','Input padding must leave room for its text');
    if(name==='button')for(const [size,height] of Object.entries({xs:20,sm:24,lg:32,icon:28,'icon-xs':20,'icon-sm':24,'icon-lg':32})){
     const variant=node.cloneNode(true);variant.dataset.size=size;fixture.append(variant);const rect=variant.getBoundingClientRect();
     assert(rect.height===height&&(!size.startsWith('icon')||rect.width===height),`Button ${size} must be ${height}px; got ${rect.width} × ${rect.height}px`);
     variant.remove();
    }
    node.remove();
   }
  }
 }finally{fixture.remove();}
}
async function checkRemounting(){
 const cases=[
  ['foundation','primitives/Pressable','ely:action','<button type="button" data-foundation-press>Press</button><output data-foundation-output></output>'],
  ['data','lists/SelectableList','ely:change','<ul data-ely-selectable><li><button type="button" data-ely-select value="first">First</button></li></ul>'],
  ['native','media/Playlist','ely:change','<button type="button" data-native-select="first">First</button>'],
  ['interaction','buttons/Toggle','ely:change','<button type="button" data-interaction-toggle aria-pressed="false">Toggle</button>'],
  ['workspace','files/FileExplorer','ely:action','<button type="button" data-ws-file-view="grid">Grid</button>'],
  ['communication','chat/MessageList','ely:action','<button type="button" data-ely-select data-ely-value="first">First</button>'],
  ['forms','account/AccountSwitcher','ely:action','<button type="button" data-ely-account value="first">First</button><button type="button" data-ely-account-menu>Accounts</button><div data-ely-account-list hidden></div><input type="hidden" data-ely-account-value>'],
 ];
 for(const [family,component,event,markup] of cases){
  const root=document.createElement('section');root.id='audit-remount-'+family;root.dataset.elyFamily=family;root.dataset.elyComponent=component;root.innerHTML=markup;
  let count=0,cleanup;root.addEventListener(event,()=>count++);
  try{
   document.body.append(root);cleanup=mount(root);root.querySelector('button').click();assert(count===1,`${family} initial mount lost or duplicated its handler`);
   const detached=new Promise(resolve=>root.addEventListener('ely:cleanup',resolve,{once:true}));root.remove();await detached;
   count=0;root.querySelector('button').click();assert(count===0,`${family} retained handlers after removal`);
   document.body.append(root);cleanup=mount(root);root.querySelector('button').click();assert(count===1,`${family} remount lost or duplicated its handler`);
  }finally{cleanup?.();root.remove();}
 }
}
export async function checkRuntime(){
 checkControlSizes();
 await checkRemounting();
 const credentials={component:'account/LoginForm',action:'authenticate',values:[['email','person@example.test'],['password','private-example']]};
 const summary=actionSummary(credentials);assert(summary.values[1][1]==='[redacted]'&&summary.values[0][1]==='person@example.test'&&credentials.values[1][1]==='private-example','Gallery output exposed credentials or changed the host payload');
 assert(actionSummary({component:'forms/PasswordInput',value:'private-example'}).value==='[redacted]','Gallery output exposed a password field action');
 assert(actionSummary({component:'account/APITokenList',values:[['token','example-key']]}).values[0][1]==='[redacted]','Gallery output exposed an account token');
 assert(actionSummary({component:'theme/ThemeEditor',values:[['token','color-paper']]}).values[0][1]==='color-paper','Gallery output hid a public theme token identifier');
 const root=document.createElement('section');root.id='audit-runtime-host';root.dataset.elyFamily='foundation';root.dataset.elyComponent='primitives/Box';root.dataset.elyKind='Box';document.body.append(root);
 const results=[],errors=[];root.addEventListener('ely:result',event=>results.push(event.detail));root.addEventListener('ely:error',event=>errors.push(event.detail));let cleanup;
 try{
  cleanup=mount(root,{actions:{save:async detail=>({value:detail.value})}});emit(root,'save',{value:'caller-value'});await waitFor(()=>results.length===1);assert(results[0].result.value==='caller-value','Host result omitted caller data');cleanup();
  cleanup=mount(root,{onAction:async()=>{throw new Error('Expected provider failure');}});emit(root,'fail');await waitFor(()=>errors.length===1);assert(errors[0].error.message==='Expected provider failure','Host failure did not reach the instance');cleanup();
  const pending=[];cleanup=mount(root,{onAction:detail=>new Promise(resolve=>pending.push({detail,resolve}))});emit(root,'update',{value:'old'});emit(root,'update',{value:'new'});pending[1].resolve('new-result');await waitFor(()=>results.length===2);pending[0].resolve('old-result');await Promise.resolve();await Promise.resolve();assert(results.length===2&&results.at(-1).result==='new-result','Stale host result replaced current data');
  cleanup();emit(root,'update',{value:'after-cleanup'});assert(pending.length===2,'Unbound host still received actions');
  root.dataset.elyDisabled='true';let disabledCalls=0;cleanup=mount(root,{onAction:()=>disabledCalls++});emit(root,'disabled-save');assert(disabledCalls===0,'Disabled instance invoked host action');
  const control=document.createElement('div');control.tabIndex=0;root.append(control);let localEvents=0;
  for(const type of ['keydown','input','change','submit']){control.addEventListener(type,()=>localEvents++);control.dispatchEvent(new Event(type,{bubbles:true,cancelable:true}));}
  assert(localEvents===0,'Disabled instance still changed through keyboard or form events');
  const tab=new KeyboardEvent('keydown',{key:'Tab',bubbles:true,cancelable:true});control.dispatchEvent(tab);assert(!tab.defaultPrevented,'Disabled instance trapped keyboard focus');
  return {name:'runtime/HostCallbacks',status:'passed',assertions:32};
 }finally{cleanup?.();root.remove();}
}
