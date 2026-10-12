// Gallery host adapter. Only the latest host result updates the editor preview.
export async function previewMarkdown(detail){
 const form=new URLSearchParams({source:detail.value||detail.source||''});
 const response=await fetch('/components/markdown-preview',{method:'POST',body:form});
 if(!response.ok)throw new Error('Markdown preview request failed.');
 const document=new DOMParser().parseFromString(await response.text(),'text/html');
 const rendered=document.querySelector('.ely-workspace__markdown');
 if(!rendered)throw new Error('Markdown renderer returned no preview.');
 return {nodes:[...rendered.childNodes]};
}
export function bindMarkdownPreview(root){
 const update=event=>{if(event.target===root&&event.detail.action==='preview-markdown')root.querySelector('[data-workspace-markdown-preview]').replaceChildren(...event.detail.result.nodes);};
 root.addEventListener('ely:result',update);
 const cleanup=()=>root.removeEventListener('ely:result',update);
 root.addEventListener('ely:cleanup',cleanup,{once:true});
 return cleanup;
}
