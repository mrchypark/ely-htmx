import {mount} from '/assets/ely-api.js';
import {previewMarkdown,bindMarkdownPreview} from './advanced-markdown.js';
document.querySelectorAll('[data-ely-family]').forEach(root=>{if(root.dataset.elyKind==='MarkdownEditor')bindMarkdownPreview(root);root.addEventListener('ely:error',event=>{document.querySelector('[data-api-host-status]').value=event.detail.error.message;});mount(root,{onAction:(detail,instance)=>{
 document.querySelector('[data-api-host-status]').value=JSON.stringify({component:detail.component,action:detail.action,value:detail.value,values:detail.values?.map(([name,value])=>[name,value instanceof File?value.name:value])});
 if(detail.action==='preview-markdown')return previewMarkdown(detail,instance);
}});});
