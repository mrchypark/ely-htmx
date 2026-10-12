import {register, emit} from './ely-api.js';
import {evaluateExpression, validateCron} from './ely-expressions.js';

const changed = (root, detail) => root.dispatchEvent(new CustomEvent('ely:change', {bubbles:true, detail:{...detail,id:root.id,component:root.dataset.elyComponent}}));
const notify = input => { input.dispatchEvent(new Event('input',{bubbles:true})); input.dispatchEvent(new Event('change',{bubbles:true})); };

register('forms', root => {
  const q = selector => root.querySelector(selector);
  const all = selector => [...root.querySelectorAll(selector)];
  let regexWorker, regexTimer, inlinePrevious;
  const cron = () => {
    const input = q('[data-ely-cron]'); if (!input) return;
    try { const value=validateCron(input.value); input.setCustomValidity(''); q('[data-ely-cron-status]').textContent=value; }
    catch(error){input.setCustomValidity(error.message); q('[data-ely-cron-status]').textContent=error.message;}
  };
  const gradient = () => {
    const host=q('[data-ely-gradient]'); if(!host)return;
    const angle=Math.max(0,Math.min(360,Number(q('[data-ely-gradient-angle]').value)||0));
    const value=`linear-gradient(${angle}deg, ${q('[data-ely-gradient-start]').value}, ${q('[data-ely-gradient-end]').value})`;
    q('[data-ely-gradient-value]').value=value;q('[data-ely-gradient-preview]').style.background=value;
  };
  const rating = () => {
    const host=q('[data-ely-rating]');if(!host)return;const value=Number(host.querySelector('input').value);
    host.querySelectorAll('[data-ely-rate]').forEach(button=>button.setAttribute('aria-pressed',String(Number(button.value)<=value)));
  };
  const checklist = () => {
    const progress=q('[data-ely-checklist-progress]');if(!progress)return;const boxes=all('.ely-api-checklist input[type=checkbox]'),count=boxes.filter(box=>box.checked).length;
    progress.value=boxes.length?count/boxes.length*100:0;q('[data-ely-checklist-count]').textContent=`${count} / ${boxes.length}`;
  };
  const range = () => {const input=q('[data-ely-range]'),end=q('[data-ely-range-end]');if(input&&end){if(Number(input.value)>Number(end.value))end.value=input.value;input.max=end.value;end.min=input.value;}if(input)q('[data-ely-range-value]').value=end?`${input.value}–${end.value}`:input.value;const dial=q('[data-ely-dial]');if(dial){dial.value=q('[data-ely-number]').value;dial.setAttribute('aria-valuetext',dial.value);}};
  const dateRange = () => {const start=q('[data-ely-date-start]'),end=q('[data-ely-date-end]');if(start&&end)end.setCustomValidity(start.value&&end.value&&end.value<start.value?'End must not be before start.':'');};
  const code = () => {const input=q('[data-ely-code-input]');if(!input||!root.dataset.elyComponent.endsWith('/JsonInput'))return;try{if(input.value)JSON.parse(input.value);input.setCustomValidity('');q('[data-ely-code-status]').textContent='';}catch(error){input.setCustomValidity('Invalid JSON.');q('[data-ely-code-status]').textContent='Invalid JSON.';}};
  const initialForms = new Map(all('[data-ely-form]').map(form=>[form,JSON.stringify([...new FormData(form)])]));
  const dirty = () => initialForms.forEach((initial,form)=>{const output=form.querySelector('[data-ely-dirty]');if(output)output.textContent=JSON.stringify([...new FormData(form)])===initial?'':'Unsaved changes';});
  const inlineClose = save => {const editor=q('[data-ely-inline-editor]'),input=q('[data-ely-inline-value]'),button=q('[data-ely-inline-open]');if(!editor||!input)return;if(!save)input.value=inlinePrevious;else{button.textContent=input.value;changed(root,{name:input.name,value:input.value});emit(root,'save',{value:input.value,name:input.name});}editor.hidden=true;button.hidden=false;button.focus();};
  const evaluate = () => {
    const input=q('[data-ely-expression]'),output=q('[data-ely-evaluation]');if(!input)return;
    if(root.dataset.elyKind==='ExpressionInput'){
      const variables=Object.fromEntries(all('[data-ely-variable]').map(node=>[node.dataset.elyVariable,node.textContent]));
      try{output.value=String(evaluateExpression(input.value,variables));input.setCustomValidity('');}catch(error){output.value=error.message;input.setCustomValidity(error.message);}return;
    }
    regexWorker?.terminate();clearTimeout(regexTimer);regexWorker=new Worker(new URL('./ely-api-regex-worker.js',import.meta.url));
    const worker=regexWorker;output.value='Evaluating';
    const finish=result=>{if(worker!==regexWorker)return;clearTimeout(regexTimer);worker.terminate();regexWorker=null;input.setCustomValidity(result.error||'');output.value=result.error||JSON.stringify(result.matches);};
    worker.onmessage=event=>finish(event.data);worker.onerror=()=>finish({error:'Regex worker failed.'});
    worker.postMessage({pattern:input.value,flags:q('[data-ely-regex-flags]').value,text:q('[data-ely-regex-text]').value});
    regexTimer=setTimeout(()=>finish({error:'Regex evaluation exceeded the time limit.'}),250);
  };
  const transfer = direction => {
    const source=q(direction==='right'?'[data-ely-transfer-source]':'[data-ely-transfer-selected]'),destination=q(direction==='right'?'[data-ely-transfer-selected]':'[data-ely-transfer-source]');
    [...source.selectedOptions].filter(option=>!option.disabled).forEach(option=>{option.selected=false;destination.append(option);});
    const template=q('[data-ely-transfer-template]');q('[data-ely-transfer-values]').replaceChildren(...[...q('[data-ely-transfer-selected]').options].map(option=>{const input=template.content.firstElementChild.cloneNode(true);input.value=option.value;return input;}));
    changed(root,{values:[...q('[data-ely-transfer-selected]').options].map(option=>option.value)});
  };
  const cascade = button => {
    const node=button.closest('[data-ely-cascade-node]'),group=node.parentElement;
    group.querySelectorAll(':scope > [data-ely-cascade-node] > [data-ely-cascade-group]').forEach(child=>child.hidden=true);
    group.querySelectorAll(':scope > [data-ely-cascade-node] > button').forEach(choice=>choice.setAttribute('aria-expanded','false'));
    const children=node.querySelector(':scope > [data-ely-cascade-group]');
    if(children){children.hidden=false;button.setAttribute('aria-expanded','true');return;}
    q('[data-ely-cascade-value]').value=button.value;
    const path=[];let cursor=node;while(cursor){path.unshift(cursor.querySelector(':scope > button')?.dataset.elyLabel||'');cursor=cursor.parentElement.closest('[data-ely-cascade-node]');}
    q('[data-ely-cascade-path]').value=path.join(' / ');changed(root,{value:button.value,path});emit(root,'select',{value:button.value,path});
  };
  const wizard=q('[data-ely-wizard]'),steps=all('[data-ely-step]'),snapshots=new Map(),wizardDisabled=Boolean(wizard&&steps.length&&steps.every(item=>item.disabled));let step=0;
  const showStep = index => {const previous=q('[data-ely-wizard-previous]'),next=q('[data-ely-wizard-next]'),finish=q('[data-ely-wizard-finish]');if(!steps.length){previous.disabled=next.disabled=finish.disabled=true;next.hidden=finish.hidden=false;q('[data-ely-step-status]').value='0 / 0';return;}step=index;steps.forEach((item,i)=>{item.hidden=i!==step;item.disabled=wizardDisabled||i!==step;});q('[data-ely-step-status]').value=`${step+1} / ${steps.length}`;previous.disabled=wizardDisabled||step===0;next.disabled=finish.disabled=wizardDisabled;next.hidden=step===steps.length-1;finish.hidden=step!==steps.length-1;};
  const stepValid = () => steps.length>0&&[...steps[step].querySelectorAll('input,textarea,select')].every(input=>input.disabled||input.reportValidity());
  const snapshot = () => {if(wizard&&steps.length)snapshots.set(step,[...new FormData(wizard)]);};
  const settingsJSON = () => {const input=q('[data-ely-settings-json]'),status=q('[data-ely-settings-json-status]');if(!input)return;try{const data=JSON.parse(input.value);input.setCustomValidity('');status.value='';return data;}catch{input.setCustomValidity('Invalid JSON.');status.value='Invalid JSON.';}};
  const click = event => {
    const button=event.target.closest('button');if(!button||button.disabled||button.closest('[data-ely-family]')!==root)return;
    if(button.matches('[data-ely-step]')){const input=q('[data-ely-number]');try{button.dataset.elyStep==='up'?input.stepUp():input.stepDown();notify(input);}catch{}return;}
    if(button.matches('[data-ely-inline-open]')){const input=q('[data-ely-inline-value]');inlinePrevious=input.value;q('[data-ely-inline-editor]').hidden=false;button.hidden=true;input.focus();return;}
    if(button.matches('[data-ely-inline-save]')){if(q('[data-ely-inline-value]').reportValidity())inlineClose(true);return;}
    if(button.matches('[data-ely-inline-cancel]')){inlineClose(false);return;}
    if(button.matches('[data-ely-evaluate]')){evaluate();return;}
    if(button.matches('[data-ely-transfer]')){transfer(button.dataset.elyTransfer);return;}
    if(button.matches('[data-ely-cascade-option]')){cascade(button);return;}
    if(button.matches('[data-ely-rate]')){q('[data-ely-rating] input').value=button.value;rating();changed(root,{value:button.value});return;}
    if(button.matches('[data-ely-chip]')){const value=q('[data-ely-chip-value]');if(root.dataset.elyComponent.endsWith('/FilterChips')){const values=new Set(value.value.split(',').filter(Boolean));values.has(button.value)?values.delete(button.value):values.add(button.value);value.value=[...values].join(',');all('[data-ely-chip]').forEach(item=>item.setAttribute('aria-pressed',String(values.has(item.value))));}else{value.value=button.value;all('[data-ely-chip]').forEach(item=>item.setAttribute('aria-pressed',String(item===button)));}changed(root,{value:value.value});return;}
    if(button.matches('[data-ely-swatch]')){const input=q('[data-ely-color]');input.value=button.value;notify(input);return;}
    if(button.matches('[data-ely-date-offset]')){const input=q('[data-ely-relative-date]'),base=q('[data-ely-base-date]').value;const date=base?new Date(`${base}T12:00:00`):new Date();date.setDate(date.getDate()+Number(button.value));input.value=`${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`;notify(input);return;}
    if(button.matches('[data-ely-mention-value]')){const input=q('[data-ely-mention]');input.value=input.value.replace(/[@#][\p{L}\p{N}_-]*$/u,`${root.dataset.elyComponent.endsWith('/HashtagInput')?'#':'@'}${button.value} `);q('[data-ely-mention-list]').hidden=true;notify(input);input.focus();return;}
    if(button.matches('[data-ely-account]')){q('[data-ely-account-value]').value=button.value;all('[data-ely-account]').forEach(item=>item.setAttribute('aria-pressed',String(item===button)));emit(root,'switch-account',{value:button.value});return;}
    if(button.matches('[data-ely-coachmark-close]')){q('[data-ely-coachmark]').hidden=true;emit(root,'dismiss');return;}
    if(button.matches('[data-ely-whats-new-open]')){q('[data-ely-whats-new]').showModal();return;}
    if(button.matches('[data-ely-whats-new-close]')){q('[data-ely-whats-new]').close();emit(root,'dismiss');return;}
    if(button.matches('[data-ely-export-settings],[data-ely-import-settings]')){const data=settingsJSON();if(data===undefined)return;emit(root,button.matches('[data-ely-export-settings]')?'export-settings':'import-settings',{data});return;}
    if(button.matches('[data-ely-array-add]')){const node=q('[data-ely-field-row]').content.firstElementChild.cloneNode(true);q('[data-ely-field-array]').append(node);node.querySelector('input').focus();dirty();return;}
    if(button.matches('[data-ely-array-remove]')){button.closest('[data-ely-array-row]').remove();dirty();return;}
    if(button.matches('[data-ely-wizard-previous]')){snapshot();showStep(Math.max(0,step-1));return;}
    if(button.matches('[data-ely-wizard-next]')){if(stepValid()){snapshot();showStep(Math.min(steps.length-1,step+1));}return;}
    if(button.matches('[data-ely-wizard-finish]')){if(!stepValid())return;snapshot();const values=[...snapshots.values()].flat(),formData=new FormData();values.forEach(([name,value])=>formData.append(name,value));emit(root,'finish',{values,formData});}
  };
  const input = event => {
    const field=event.target;if(field.closest('[data-ely-family]')!==root)return;
    if(field.matches('[data-ely-autoresize]')){field.style.height='auto';field.style.height=`${Math.min(600,field.scrollHeight)}px`;}
    if(field.matches('[data-ely-mask]')&&field.dataset.elyMask){let index=0;const digits=field.value.replace(/\D/g,'');field.value=[...field.dataset.elyMask].map(char=>char==='#'?(digits[index++]||''):(index<digits.length?char:'')).join('');}
    if(field.matches('[data-ely-mention]')){const match=field.value.match(/[@#]([\p{L}\p{N}_-]*)$/u);q('[data-ely-mention-list]').hidden=!match;all('[data-ely-mention-value]').forEach(button=>button.hidden=Boolean(match)&&!button.value.toLocaleLowerCase().startsWith(match[1].toLocaleLowerCase()));}
    if(field.matches('[data-ely-settings-search]'))all('[data-ely-settings-section]').forEach(section=>section.hidden=!section.textContent.toLocaleLowerCase().includes(field.value.toLocaleLowerCase()));
    if(field.matches('[data-ely-picker]'))q('[data-ely-picker-preview]').value=field.selectedOptions[0]?.textContent||'';
    if(field.matches('[data-ely-dial]')){q('[data-ely-number]').value=field.value;notify(q('[data-ely-number]'));}
    settingsJSON();cron();gradient();rating();range();dateRange();code();dirty();checklist();
  };
  const key = event => {if(event.target.matches('[data-ely-inline-value]')){if(event.key==='Escape'){event.preventDefault();inlineClose(false);}if(event.key==='Enter'){event.preventDefault();if(event.target.reportValidity())inlineClose(true);}}};
  let scrub;
  const down = event => {if(!event.target.matches('[data-ely-scrub]')||q('[data-ely-number]').disabled)return;scrub={x:event.clientX,value:Number(q('[data-ely-number]').value)||0};if(event.isTrusted)event.target.setPointerCapture(event.pointerId);};
  const move = event => {if(!scrub)return;const field=q('[data-ely-number]'),value=scrub.value+Math.round((event.clientX-scrub.x)/8)*(Number(field.step)||1);field.value=String(Math.max(field.min===''?-Infinity:Number(field.min),Math.min(field.max===''?Infinity:Number(field.max),value)));notify(field);};
  const up = () => {scrub=null;};
  root.addEventListener('click',click);root.addEventListener('input',input);root.addEventListener('change',input);root.addEventListener('keydown',key);root.addEventListener('pointerdown',down);root.addEventListener('pointermove',move);root.addEventListener('pointerup',up);root.addEventListener('pointercancel',up);root.addEventListener('reset',()=>queueMicrotask(dirty));
  cron();gradient();rating();range();dateRange();code();settingsJSON();checklist();if(wizard)showStep(0);const whatsNew=q('[data-ely-whats-new]');if(whatsNew?.dataset.elyWhatsNewInitial)whatsNew.showModal?.();
  return()=>{regexWorker?.terminate();clearTimeout(regexTimer);root.removeEventListener('click',click);root.removeEventListener('input',input);root.removeEventListener('change',input);root.removeEventListener('keydown',key);root.removeEventListener('pointerdown',down);root.removeEventListener('pointermove',move);root.removeEventListener('pointerup',up);root.removeEventListener('pointercancel',up);};
});
