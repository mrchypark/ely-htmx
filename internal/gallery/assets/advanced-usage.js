import {mount} from '/assets/ely-api.js';
import {previewMarkdown,bindMarkdownPreview} from './advanced-markdown.js';
import {actionSummary} from './action-summary.js';
document.querySelectorAll('[data-ely-family]').forEach(root=>{if(root.dataset.elyKind==='MarkdownEditor')bindMarkdownPreview(root);root.addEventListener('ely:error',event=>{const status=document.querySelector('[data-api-host-status]');if(status)status.value=event.detail.error.message;});mount(root,{onAction:(detail,instance)=>{
 const status=document.querySelector('[data-api-host-status]');if(status)status.value=JSON.stringify(actionSummary(detail));
 if(detail.action==='preview-markdown')return previewMarkdown(detail,instance);
}});});
