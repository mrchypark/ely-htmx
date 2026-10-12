import {register, emit} from './ely-api.js';

register('workspace', root => {
  const change = detail => root.dispatchEvent(new CustomEvent('ely:change', {bubbles: true, detail: {id: root.id, component: root.dataset.elyComponent, ...detail}}));
  let drag;
  const canvasDisabled = () => root.querySelector('[data-workspace-canvas]')?.dataset.wsDisabled === 'true';
  const chapter = root.dataset.elyComponent.split('/')[0];
  let canvasMode = 'pan', pan = {x:0,y:0}, drawing, panning, selectedRange;
  const stage = root.querySelector('[data-workspace-stage]');
  const transformStage = () => { if (stage) { const scale = Number.parseInt(root.querySelector('output[data-workspace-zoom]').textContent,10) / 100; stage.style.transform = 'translate(' + pan.x + 'px,' + pan.y + 'px) scale(' + scale + ')'; stage.style.transformOrigin = '0 0'; } };
  const stagePoint = event => { const matrix = stage.getScreenCTM(); if (!matrix) return null; const point = new DOMPoint(event.clientX,event.clientY).matrixTransform(matrix.inverse()); return {x:Math.round(point.x*100)/100,y:Math.round(point.y*100)/100}; };
  const whitespace = root.querySelector('[data-ws-whitespace]');
  if (whitespace) whitespace.textContent = whitespace.textContent.replace(/ /g,'·').replace(/\t/g,'→').replace(/\n/g,'↵\n');
  const graphNodes = [...root.querySelectorAll('[data-workspace-graph] g[data-parents]')];
  const graphByID = new Map(graphNodes.map(node => [node.dataset.value,node]));
  graphNodes.forEach(node => node.dataset.parents.trim().split(/\s+/).filter(Boolean).forEach(id => {
    const parent = graphByID.get(id);
    if (!parent) return;
    const values = [node.dataset.x,node.dataset.y,parent.dataset.x,parent.dataset.y].map(Number);
    if (!values.every(Number.isFinite)) return;
    const path = document.createElementNS('http://www.w3.org/2000/svg','path');path.dataset.wsParentEdge=id;path.setAttribute('d','M'+values[0]+' '+values[1]+' L'+values[2]+' '+values[3]);node.closest('svg').prepend(path);
  }));
  const markdown = root.querySelector('[data-workspace-markdown-renderer] .ely-workspace__markdown');
  if (markdown) {
    const nav = root.querySelector('[data-ws-markdown-nav]');
    markdown.querySelectorAll('h1,h2,h3,h4,h5,h6').forEach((heading,index) => {
      heading.id = root.id + '-heading-' + (index+1);
      if (nav) {
        const link = document.createElement('a');
        link.href = '#' + heading.id;
        link.textContent = heading.textContent;
        link.dataset.level = heading.tagName.slice(1);
        nav.append(link);
      }
    });
  }
  const reading = root.querySelector('[data-ws-reading]');
  const readingProgress = () => {
    const extent = reading.scrollHeight - reading.clientHeight;
    root.querySelector('[data-ws-reading-progress]').value = extent > 0 ? Math.min(100,Math.max(0,reading.scrollTop / extent * 100)) : 100;
  };
  if (reading) { reading.addEventListener('scroll',readingProgress); readingProgress(); }
  // Build a fresh tree: formatting survives, caller attributes and active content do not.
  const safeRich = node => {
    const result = document.createElement('div');
    const allowed = new Set(['P','DIV','BR','STRONG','EM','B','I','U','S','UL','OL','LI','H1','H2','H3','H4','H5','H6','BLOCKQUOTE','CODE','PRE']);
    let count = 0;
    const copy = (source, destination, depth = 0) => {
      if (depth > 100) throw new Error('Content nesting exceeds 100 levels.');
      for (const child of source.childNodes) {
        if (++count > 10000) throw new Error('Content exceeds 10000 nodes.');
        if (child.nodeType === Node.TEXT_NODE) destination.append(document.createTextNode(child.textContent));
        else if (child.nodeType === Node.ELEMENT_NODE && !['SCRIPT','STYLE','IFRAME','OBJECT','EMBED'].includes(child.tagName)) {
          if (allowed.has(child.tagName)) { const clean = document.createElement(child.tagName.toLowerCase()); copy(child,clean,depth+1); destination.append(clean); }
          else copy(child,destination,depth+1);
        }
      }
    };
    copy(node,result);
    return result.innerHTML;
  };
  const codeText = editor => {
    let count = 0;
    const read = (node, depth = 0) => {
      if (++count > 10000 || depth > 100) throw new Error('Code content exceeds the supported DOM size.');
      if (node.nodeType === Node.TEXT_NODE) return node.textContent;
      if (node.nodeType === Node.ELEMENT_NODE && node.getAttribute('contenteditable') === 'false') return '';
      if (node.nodeName === 'BR') return '\n';
      const text = [...node.childNodes].map(child => read(child,depth+1)).join('');
      return ['DIV','P','LI'].includes(node.nodeName) && !text.endsWith('\n') ? text+'\n' : text;
    };
    const nodes = [...editor.childNodes], value = nodes.map(node => read(node)).join('');
    return ['DIV','P','LI'].includes(nodes.at(-1)?.nodeName) ? value.replace(/\n$/,'') : value;
  };
  const syncEditor = editor => {
    const form = editor.closest('form[data-ws-edit-form]');
    if (!form) return editor.textContent;
    const code = form.hasAttribute('data-workspace-code');
    const value = code ? codeText(editor) : safeRich(editor);
    form.querySelector('input[data-ws-content-value]').value = value;
    return value;
  };
  root.querySelectorAll('[data-workspace-rich]').forEach(syncEditor);
  const editError = (editor,message) => { editor.closest('form[data-ws-edit-form]').querySelector('[data-ws-edit-error]').textContent = message; editor.setAttribute('aria-invalid','true'); };
  const validateEditor = event => {
    const form = event.target;
    if (!(form instanceof HTMLFormElement) || !form.matches('[data-ws-edit-form]')) return;
    const editor = form.querySelector('[data-workspace-rich]');
    try { syncEditor(editor); } catch (error) { event.preventDefault(); event.stopImmediatePropagation(); editError(editor,error.message); editor.focus(); return; }
    const current = form.hasAttribute('data-workspace-code') ? form.querySelector('[data-ws-content-value]').value : editor.textContent;
    const missing = form.dataset.wsRequired === 'true' && !current.trim();
    if (missing || editor.contentEditable === 'false') {
      event.preventDefault(); event.stopImmediatePropagation();
      if (missing) { editError(editor,'Enter content before submitting.'); editor.focus(); }
    }
  };
  // Capture runs before the base form dispatcher reads FormData.
  root.addEventListener('submit',validateEditor,true);
  const refreshWords = node => { const output = node.parentElement.querySelector('[data-workspace-word-count]'); if (output) output.textContent = node.textContent.trim().split(/\s+/).filter(Boolean).length + ' words'; };
  root.querySelectorAll('[data-workspace-rich]').forEach(refreshWords);
  const trackSelection = () => {
    const selection = document.getSelection();
    if (!selection?.rangeCount) return;
    const range = selection.getRangeAt(0), editor = root.querySelector('[data-workspace-rich]');
    if (!editor || !editor.contains(range.commonAncestorContainer)) return;
    selectedRange = range.cloneRange();
    if (chapter === 'editor') {
      const before = document.createRange();before.selectNodeContents(editor);before.setEnd(range.endContainer,range.endOffset);
      const fragment = before.cloneContents();fragment.querySelectorAll('[aria-hidden="true"]').forEach(node=>node.remove());
      const prefix = codeText(fragment), parts = prefix.split('\n'), position = {line:parts.length,column:parts.at(-1).length+1};
      const output=root.querySelector('[data-workspace-cursor]');if(output)output.textContent='Ln '+position.line+', Col '+position.column;
      change({name:'cursor',value:position});
    }
  };
  document.addEventListener('selectionchange',trackSelection);
  root.addEventListener('pointerdown',event=>{if(event.target.closest('[data-ws-format]'))event.preventDefault();});
  const search = input => {
    const rows = [...root.querySelectorAll('[data-ws-search-row]')];
    const term = input.value.toLocaleLowerCase();
    rows.forEach(row => { row.hidden = !row.textContent.toLocaleLowerCase().includes(term); });
    root.querySelector('[data-ws-search-count]').textContent = rows.filter(row => !row.hidden).length + ' results';
  };
  root.querySelectorAll('input[data-ws-search]').forEach(search);
  const todoCount = () => { const boxes = [...root.querySelectorAll('input[data-ws-todo]')]; const output = root.querySelector('[data-ws-todo-count]'); if (output) output.textContent = boxes.filter(box => box.checked).length + ' / ' + boxes.length + ' completed'; };
  todoCount();
  const redisSearch = input => {
    const term = input.value.toLocaleLowerCase();
    root.querySelectorAll('[data-ws-redis-entry]').forEach(entry => { entry.hidden = Boolean(term) && !entry.dataset.key.toLocaleLowerCase().includes(term); });
  };
  root.querySelectorAll('input[data-ws-redis-search]').forEach(redisSearch);
  const hostOutput = event => {
    const form = root.querySelector('form[data-ws-graphql],form[data-ws-regex]');
    const output = form?.querySelector('[data-ws-host-result]');
    if (!form || !output || event.detail?.action !== form.dataset.elySubmit) return;
    let text;
    if (event.type === 'ely:error') text = event.detail.error?.message || String(event.detail.error || 'Host operation failed.');
    else {
      if (!Object.hasOwn(event.detail, 'result')) return;
      try { text = typeof event.detail.result === 'string' ? event.detail.result : JSON.stringify(event.detail.result); } catch { text = 'Host result could not be displayed.'; }
    }
    output.textContent = String(text).slice(0, 10000);
    output.dataset.wsResultState = event.type === 'ely:error' ? 'error' : 'result';
    output.hidden = false;
  };
  root.addEventListener('ely:result',hostOutput);
  root.addEventListener('ely:error',hostOutput);
  root.addEventListener('click', event => {
    const format = event.target.closest('[data-ws-format]');
    if (format && !format.disabled) {
      const editor=root.querySelector('[data-workspace-rich]'), range=selectedRange;
      if (!range || range.collapsed || !editor.contains(range.commonAncestorContainer)) editError(editor,'Select text before formatting.');
      else if (['strong','em','code'].includes(format.dataset.wsFormat)) { const tag=document.createElement(format.dataset.wsFormat);tag.append(range.extractContents());range.insertNode(tag);const next=document.createRange();next.selectNodeContents(tag);const selection=document.getSelection();selection.removeAllRanges();selection.addRange(next);selectedRange=next.cloneRange();editor.dispatchEvent(new Event('input',{bubbles:true})); }
    }
    const block = event.target.closest('[data-ws-insert-block]');
    if (block && !block.disabled) { const editor=root.querySelector('[data-workspace-rich]');const tag=['p','h2','blockquote','pre'].includes(block.dataset.wsInsertBlock)?block.dataset.wsInsertBlock:'p';const node=document.createElement(tag);node.textContent=block.dataset.blockText;editor.append(node);editor.dispatchEvent(new Event('input',{bubbles:true}));change({name:'block',value:{id:block.dataset.blockId,type:tag,text:node.textContent}}); }
    const completion=event.target.closest('[data-ws-complete]');
    if(completion && !completion.disabled){const source=root.querySelector('textarea[data-ws-completion-source]');source.setRangeText(completion.dataset.insert,source.selectionStart,source.selectionEnd,'end');root.querySelector('[data-ws-completion-status]').textContent='Inserted '+completion.textContent;emit(root,'select-completion',{candidate:completion.dataset.wsComplete,insertText:completion.dataset.insert,source:source.value});change({name:'source',value:source.value});}
    const focus=event.target.closest('[data-ws-focus-toggle]');
    if(focus && !focus.disabled){const active=focus.getAttribute('aria-pressed')!=='true';focus.setAttribute('aria-pressed',String(active));root.querySelector('[data-ws-focus]').dataset.focused=String(active);root.querySelector('[data-ws-focus-chrome]').hidden=active;if(active)root.querySelector('[data-ws-focus-content]').focus();change({name:'focus',value:active});}
    const mode=event.target.closest('[data-ws-canvas-mode]');
    if(mode && !mode.disabled){canvasMode=mode.dataset.wsCanvasMode;root.querySelectorAll('[data-ws-canvas-mode]').forEach(button=>button.setAttribute('aria-pressed',String(button===mode)));change({name:'canvasMode',value:canvasMode});}
    const cancel = event.target.closest('[data-ws-rename-cancel]');
    if (cancel && !cancel.disabled) { cancel.form.reset(); root.querySelector('[data-ws-status]').textContent = 'Edit canceled'; }
    const page = event.target.closest('[data-ws-page]');
    if (page && !page.disabled) root.querySelectorAll('[data-ws-page-content]').forEach(node => { node.hidden = node.dataset.wsPageContent !== page.dataset.wsPage; });
    const revision = event.target.closest('[data-ws-revision]');
    if (revision && !revision.disabled) { root.querySelectorAll('[data-ws-revision-content]').forEach(node => { node.hidden = node.dataset.wsRevisionContent !== revision.dataset.wsRevision; }); change({name: 'revision', value: revision.textContent.trim()}); }
    const select = event.target.closest('[data-workspace-select]');
    if (select && root.contains(select)) { root.querySelectorAll('[data-workspace-select]').forEach(node => node.setAttribute('aria-pressed', String(node === select))); change({name: 'selection', value: select.dataset.value || select.textContent.trim()}); }
    const history = event.target.closest('[data-workspace-history]');
    if (history && root.contains(history)) { const command = root.querySelector('[data-workspace-terminal] input[name="command"]'); if (command) { command.value = history.value; command.focus(); change({name: 'command', value: command.value}); } }
    const tree = event.target.closest('[data-workspace-tree]');
    if (tree && root.contains(tree)) { const child = tree.parentElement.querySelector(':scope > ul'); if (child) { const expanded = tree.getAttribute('aria-expanded') !== 'true'; tree.setAttribute('aria-expanded', String(expanded)); child.hidden = !expanded; } change({name: 'file', value: tree.dataset.value}); }
    const redisKey = event.target.closest('[data-ws-redis-key]');
    if (redisKey && !redisKey.disabled && root.contains(redisKey)) { root.querySelectorAll('[data-ws-redis-key]').forEach(node => node.setAttribute('aria-pressed', String(node === redisKey))); const preview = root.querySelector('[data-ws-redis-preview]'); if (preview) preview.textContent = redisKey.dataset.preview || ''; emit(root, 'select-redis-key', {key: redisKey.dataset.key, preview: redisKey.dataset.preview || ''}); change({name: 'redis-key', value: redisKey.dataset.key}); }
    const zoom = event.target.closest('button[data-workspace-zoom]');
    if (zoom && !zoom.disabled && root.contains(zoom)) { const output = zoom.parentElement.querySelector('output[data-workspace-zoom]'); const value = Math.max(50, Math.min(200, Number.parseInt(output.textContent, 10) + (zoom.dataset.workspaceZoom === 'in' ? 10 : -10))); output.textContent = value + '%'; transformStage(); change({name: 'zoom', value}); }
  });
  root.addEventListener('input', event => {
    const target = event.target;
    if(target.matches('input[data-ws-candidate-filter]'))root.querySelectorAll('[data-ws-candidate]').forEach(row=>{row.hidden=!row.textContent.toLocaleLowerCase().includes(target.value.toLocaleLowerCase());});
    if (target.matches('input[data-ws-search]')) search(target);
    if (target.matches('input[data-ws-redis-search]')) { redisSearch(target); change({name: target.name, value: target.value}); }
    if (target.matches('input[data-ws-color]')) { root.querySelector('[data-ws-color-preview]').setAttribute('fill', target.value); root.querySelector('[data-ws-color-label]').textContent = target.value; }
    const editor = target.closest('[data-workspace-rich]');
    if (editor) { try { const value = syncEditor(editor); refreshWords(editor); editor.removeAttribute('aria-invalid'); editor.closest('form[data-ws-edit-form]').querySelector('[data-ws-edit-error]').textContent = ''; change({name: editor.dataset.name, value}); } catch (error) { editError(editor,error.message); } }
    if (target.matches('textarea[data-workspace-markdown]')) { emit(root, 'preview-markdown', {source: target.value}); change({name: target.name, value: target.value}); }
    if (target.matches('input[data-workspace-find],input[data-workspace-filter]')) { const term = target.value.toLowerCase(); root.querySelectorAll('.ely-workspace__items li,.ely-workspace__tree li').forEach(item => { item.hidden = Boolean(term) && !item.textContent.toLowerCase().includes(term); }); const matches = [...root.querySelectorAll('.ely-workspace__code code')].filter(node => node.textContent.toLowerCase().includes(term)).length; const output = root.querySelector('[data-workspace-match-count]'); if (output) output.textContent = term ? matches + ' matches' : ''; change({name: target.name, value: target.value}); }
    if (target.matches('input[type="file"]')) { const output = root.querySelector('[data-workspace-file-count]'); if (output) output.textContent = target.files.length + ' file(s) selected'; change({name: target.name, files: [...target.files].map(file => ({name: file.name, size: file.size, type: file.type}))}); }
  });
  root.addEventListener('change', event => {
    const target = event.target;
    if (target.disabled) return;
    if (chapter==='git' && target.matches('input[name="resolution"]')) { const source = target.closest('fieldset').querySelectorAll('pre');const merged=root.querySelector('textarea[name="merged"]');merged.value=source[target.value==='theirs'?1:0].textContent;change({name:'merged',value:merged.value}); }
    if (target.matches('input[data-ws-breakpoint]')) emit(root, 'toggle-breakpoint', {breakpoint: target.value, file: target.dataset.file, line: Number(target.dataset.line), enabled: target.checked});
    if (target.matches('input[data-ws-todo]')) todoCount();
    if (target.matches('select[data-ws-branch]')) { root.querySelector('[data-ws-branch-label]').textContent = target.selectedOptions[0]?.textContent || ''; emit(root, 'switch-branch', {branch: target.value}); }
    if (target.matches('select[data-ws-tool]')) { root.querySelector('[data-ws-tool-label]').textContent = target.value; root.dataset.workspaceTool = target.value; }
    if (target.matches('input[data-ws-layer]')) { root.querySelectorAll('[data-workspace-shape]').forEach(shape => { if (shape.dataset.value === target.value) shape.style.display = target.checked ? '' : 'none'; }); change({name: 'layer', value: target.value, visible: target.checked}); }
    if (target.matches('input[type="file"]')) { const output = root.querySelector('[data-workspace-file-count]'); if (output) output.textContent = target.files.length + ' file(s) selected'; }
  });
  root.addEventListener('pointerdown', event => { const shape = event.target.closest('[data-workspace-shape]'); if (!shape || canvasMode==='draw' || canvasDisabled() || !root.contains(shape) || event.button !== 0) return; const point = shape.getAttribute('transform').match(/-?\d+(?:\.\d+)?/g) || [0, 0]; drag = {shape, x: event.clientX, y: event.clientY, left: Number(point[0]), top: Number(point[1])}; if (event.isTrusted) { try { shape.setPointerCapture?.(event.pointerId); } catch { drag = undefined; return; } } change({name: 'selection', value: shape.dataset.value}); });
  root.addEventListener('pointermove', event => { if (!drag) return; const x = drag.left + event.clientX - drag.x, y = drag.top + event.clientY - drag.y; drag.shape.setAttribute('transform', 'translate(' + x + ' ' + y + ')'); change({name: 'position', value: {id: drag.shape.dataset.value, x, y}}); });
  root.addEventListener('pointerup', () => { drag = undefined; });
  root.addEventListener('pointerdown',event=>{
    if(!stage || !stage.contains(event.target) || canvasDisabled() || event.button!==0) return;
    if(canvasMode==='draw') { const point=stagePoint(event);if(!point)return;const line=document.createElementNS('http://www.w3.org/2000/svg','polyline');line.setAttribute('fill','none');line.setAttribute('stroke','currentColor');line.setAttribute('stroke-width','2');root.querySelector('[data-ws-strokes]').append(line);drawing={line,points:[point]};line.setAttribute('points',point.x+','+point.y); }
    else if(!event.target.closest('[data-workspace-shape]'))panning={x:event.clientX,y:event.clientY,left:pan.x,top:pan.y};
    if((drawing || panning) && event.isTrusted){try{stage.setPointerCapture(event.pointerId);}catch{/* A canceled native pointer may no longer be capturable. */}}
  });
  root.addEventListener('pointermove',event=>{
    if(panning){pan={x:panning.left+event.clientX-panning.x,y:panning.top+event.clientY-panning.y};transformStage();change({name:'pan',value:{...pan}});}
    if(drawing && drawing.points.length<5000){const point=stagePoint(event);if(point){drawing.points.push(point);drawing.line.setAttribute('points',drawing.points.map(p=>p.x+','+p.y).join(' '));}}
  });
  root.addEventListener('pointerup',()=>{
    if(drawing){emit(root,'draw-stroke',{points:drawing.points,color:'currentColor',width:2});drawing=undefined;}panning=undefined;
  });
  root.addEventListener('pointercancel',()=>{drag=undefined;panning=undefined;if(drawing)drawing.line.remove();drawing=undefined;});
  root.addEventListener('keydown', event => { const shape = event.target.closest('[data-workspace-shape]'); if (shape && !canvasDisabled() && (event.key === 'Enter' || event.key === ' ')) { event.preventDefault(); emit(root, 'select-shape', {shape: shape.dataset.value}); } });
  root.addEventListener('keydown',event=>{if(event.target.matches('[data-workspace-graph] g[data-workspace-select]') && ['Enter',' '].includes(event.key)){event.preventDefault();event.target.dispatchEvent(new MouseEvent('click',{bubbles:true}));}});
  return () => { document.removeEventListener('selectionchange',trackSelection); root.removeEventListener('submit',validateEditor,true); root.removeEventListener('ely:result',hostOutput); root.removeEventListener('ely:error',hostOutput); if (reading) reading.removeEventListener('scroll',readingProgress); };
});
