import {mount,emit} from '/assets/ely-api.js';
import {assert,waitFor} from './advanced-checks.js';
export async function checkRuntime(){
 const root=document.createElement('section');root.id='audit-runtime-host';root.dataset.elyFamily='foundation';root.dataset.elyComponent='primitives/Box';root.dataset.elyKind='Box';document.body.append(root);
 const results=[],errors=[];root.addEventListener('ely:result',event=>results.push(event.detail));root.addEventListener('ely:error',event=>errors.push(event.detail));let cleanup;
 try{
  cleanup=mount(root,{actions:{save:async detail=>({value:detail.value})}});emit(root,'save',{value:'caller-value'});await waitFor(()=>results.length===1);assert(results[0].result.value==='caller-value','Host result omitted caller data');cleanup();
  cleanup=mount(root,{onAction:async()=>{throw new Error('Expected provider failure');}});emit(root,'fail');await waitFor(()=>errors.length===1);assert(errors[0].error.message==='Expected provider failure','Host failure did not reach the instance');cleanup();
  const pending=[];cleanup=mount(root,{onAction:detail=>new Promise(resolve=>pending.push({detail,resolve}))});emit(root,'update',{value:'old'});emit(root,'update',{value:'new'});pending[1].resolve('new-result');await waitFor(()=>results.length===2);pending[0].resolve('old-result');await Promise.resolve();await Promise.resolve();assert(results.length===2&&results.at(-1).result==='new-result','Stale host result replaced current data');
  cleanup();emit(root,'update',{value:'after-cleanup'});assert(pending.length===2,'Unbound host still received actions');
  root.dataset.elyDisabled='true';let disabledCalls=0;cleanup=mount(root,{onAction:()=>disabledCalls++});emit(root,'disabled-save');assert(disabledCalls===0,'Disabled instance invoked host action');
  return {name:'runtime/HostCallbacks',status:'passed',assertions:5};
 }finally{cleanup?.();root.remove();}
}
