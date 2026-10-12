import {registerChecks, assert, waitFor} from './advanced-checks.js';

const action = (actions, name) => assert(actions.some(entry => entry.action === name), 'missing ' + name + ' action');
const one = (root, selector) => { const node = root.querySelector(selector); assert(Boolean(node), 'missing expected control: ' + selector); return node; };
const checkRequiredEditor = (root, actions, editor) => {
  const form = one(root,'form[data-ws-edit-form]');
  if (form.dataset.wsRequired !== 'true') return;
  editor.textContent = ''; editor.dispatchEvent(new Event('input',{bubbles:true}));
  const count = actions.length;
  one(form,'button[type="submit"]').click();
  assert(actions.length === count, 'empty required editor emitted a host action');
  assert(editor.getAttribute('aria-invalid') === 'true' && one(form,'[data-ws-edit-error]').textContent === 'Enter content before submitting.', 'required editor error is absent');
};
const code = async ({root,actions,twin}) => {
  const editor = one(root, '[data-workspace-rich]');
  const twinBefore = twin?.querySelector('[data-workspace-rich]')?.textContent;
  const line=document.createElement('div'),gutter=document.createElement('span'),source=document.createElement('code');gutter.contentEditable='false';gutter.textContent='1';source.textContent='first line';line.append(gutter,source);
  const added=document.createElement('div');added.textContent='new browser line';editor.replaceChildren(line,added);editor.dispatchEvent(new Event('input',{bubbles:true}));submitData(root,actions,editor.dataset.name,'first line\nnew browser line');
  if(editor.closest('form').dataset.wsRequired==='true'){source.textContent='';editor.replaceChildren(line);const count=actions.length;one(root,'button[type="submit"]').click();assert(actions.length===count,'line numbers satisfied required code input');}
  editor.textContent = 'changed caller code'; editor.dispatchEvent(new Event('input', {bubbles: true}));
  submitData(root,actions,editor.dataset.name,'changed caller code');
  if(twin)assert(one(twin,'[data-workspace-rich]').textContent === twinBefore,'code edit changed twin');
  checkRequiredEditor(root,actions,editor);
};
const find = async ({root}) => {
  const input = one(root, 'input[data-workspace-find]'); input.value = 'caller'; input.dispatchEvent(new Event('input', {bubbles: true}));
  assert(one(root, '[data-workspace-match-count]').textContent==='2 matches', 'find did not count caller lines');
};
const diff = async ({root, actions}) => {
  assert(root.querySelector('.ely-workspace__diff-add') && root.querySelector('.ely-workspace__diff-remove'), 'diff hunks are absent');
  const button = one(root, '[data-ely-action]'); button.click(); action(actions, button.dataset.elyAction);
};
const conflict = async ({root, actions}) => {
  const theirs = one(root, 'input[value="theirs"]'); theirs.click(); assert(theirs.checked, 'conflict choice did not change');
  const merged = one(root, 'textarea[name="merged"]');assert(merged.value==='caller theirs','conflict choice did not populate merged result');merged.value = 'resolved caller value';
  assert(submitData(root,actions,'merged','resolved caller value').get('resolution')==='theirs','conflict lost resolution');
};
const graph = async ({root}) => {
  assert(one(root,'[data-ws-parent-edge="commit-parent"]').getAttribute('d')==='M260 120 L80 80','commit parents did not render edges');
  const node = one(root, 'svg g[data-workspace-select]'); node.dispatchEvent(new MouseEvent('click', {bubbles: true}));
  assert(node.getAttribute('aria-pressed') === 'true', 'graph node was not selected');
};
const blame = async ({root}) => {
  const row = one(root, 'tbody tr[data-workspace-select]'); row.dispatchEvent(new MouseEvent('click', {bubbles: true}));
  assert(row.getAttribute('aria-pressed') === 'true', 'blame row was not selected');
};
const terminal = async ({root, actions}) => {
  const history = one(root, '[data-workspace-history]'); history.click();
  const command = one(root, 'input[name="command"]'); assert(command.value === history.value, 'history did not populate command');
  const form = one(root, 'form'); form.dispatchEvent(new Event('submit', {bubbles: true, cancelable: true})); action(actions, form.dataset.elySubmit);
};
const rich = async ({root,actions,twin}) => {
  const editor = one(root, '[data-workspace-rich]'); editor.textContent = 'one two three'; editor.dispatchEvent(new Event('input', {bubbles: true}));
  assert(one(root, '[data-workspace-word-count]').textContent === '3 words', 'word count did not reflect edit');
  submitData(root,actions,editor.dataset.name,'one two three');
  const range=document.createRange();range.selectNodeContents(editor);const selection=document.getSelection();selection.removeAllRanges();selection.addRange(range);document.dispatchEvent(new Event('selectionchange'));
  one(root,'[data-ws-format="strong"]').click();
  assert(one(editor,'strong').textContent==='one two three','native range formatting failed');
  submitData(root,actions,editor.dataset.name,'<strong>one two three</strong>');
  const twinBefore=twin?.querySelector('[data-workspace-rich]')?.textContent;
  editor.replaceChildren();
  const bold=document.createElement('strong');bold.textContent='Caller format';bold.setAttribute('onclick','hostile');
  const link=document.createElement('a');link.href='javascript:hostile';link.textContent=' link';
  editor.append(bold,link);editor.dispatchEvent(new Event('input',{bubbles:true}));
  submitData(root,actions,editor.dataset.name,'<strong>Caller format</strong> link');
  if(twin)assert(one(twin,'[data-workspace-rich]').textContent===twinBefore,'rich edit changed twin');
  checkRequiredEditor(root,actions,editor);
  if(root.dataset.elyKind==='CommentThread')checkComment(root,actions);
};
const markdown = async ({root, actions}) => {
  const source = one(root, 'textarea[data-workspace-markdown]'); source.value = '# caller'; source.dispatchEvent(new Event('input', {bubbles: true}));
  action(actions, 'preview-markdown');
  await waitFor(() => root.querySelector('[data-workspace-markdown-preview] h1')?.textContent === 'caller');
  assert(one(root, '[data-workspace-markdown-preview] h1').textContent === 'caller', 'host preview did not render caller heading');
};
const markdownRenderer = async ({root,twin}) => {
  assert(!root.querySelector('textarea,[contenteditable="true"]'), 'renderer must be read-only');
  const markdown=one(root,'.ely-workspace__markdown'), heading=one(markdown,'h1');
  assert(heading.textContent==='Caller markdown','caller markdown heading missing');
  const headings=[...markdown.querySelectorAll('h1,h2,h3,h4,h5,h6')];
  assert(headings.every(h=>h.id.startsWith(root.id+'-heading-')),'heading IDs must be local');
  if(twin)assert(headings.every(h=>!twin.querySelector('[id="'+h.id+'"]')),'heading IDs overlap twin');
  const variant=root.dataset.elyComponent.split('/')[1];
  if(variant==='TableOfContents' || variant==='DocumentOutline') {
    const links=[...one(root,'[data-ws-markdown-nav]').querySelectorAll('a')];
    assert(links.length===2 && links.length===headings.length,'heading navigation incomplete');
    assert(links.every((link,i)=>link.hash==='#'+headings[i].id && link.textContent===headings[i].textContent),'heading navigation targets wrong');
    links[1].click(); await waitFor(()=>location.hash==='#'+headings[1].id);
    assert(location.hash===links[1].hash,'heading link failed to navigate');
  }
  if(variant==='ReadingProgress') {
    const area=one(root,'[data-ws-reading]'), progress=one(root,'progress[data-ws-reading-progress]');
    assert(area.scrollHeight>area.clientHeight && progress.value<100,'reading fixture must overflow');
    area.scrollTop=area.scrollHeight; area.dispatchEvent(new Event('scroll'));
    assert(progress.value===100 && area.scrollTop>0,'reading progress did not reach end');
    area.scrollTop=0;area.dispatchEvent(new Event('scroll'));assert(progress.value===0,'reading progress did not return to start');
  }
};
const files = async ({root}) => {
  const toggle = one(root, '[data-workspace-tree]'), branch=one(toggle.parentElement,':scope > ul'); toggle.click();
  assert(toggle.getAttribute('aria-expanded') === 'false' && branch.hidden, 'tree did not collapse');
  toggle.click(); assert(!branch.hidden && toggle.getAttribute('aria-expanded')==='true','tree did not expand');
};
const data = async ({root}) => {
  const row = one(root, 'tr[data-workspace-select]');
  row.dispatchEvent(new MouseEvent('click', {bubbles: true})); assert(row.getAttribute('aria-pressed') === 'true', 'data row did not select');
};
const filesDevtools = async ({root,actions,expectError}) => {
  const kind=root.dataset.elyKind;
  switch(kind) {
    case 'FileIcon': { const icon=one(root,'[data-ws-file-icon]'); assert(icon.dataset.extension==='md' && icon.textContent==='MD' && root.textContent.includes('docs/brief.md'),'file icon or caller metadata missing'); break; }
    case 'FileItem': { const card=one(root,'[data-ws-file-card]'); assert(card.dataset.path==='docs/brief.md' && one(card,'.ely-workspace__file-meta').textContent.includes('4 KB') && card.textContent.includes('2026-10-11'),'file item metadata missing'); break; }
    case 'FileGrid': { const cards=[...root.querySelectorAll('[data-ws-file-card]')]; assert(cards.length===2 && cards[0].dataset.path==='docs/brief.md' && cards[1].textContent.includes('budget.csv'),'file grid did not render caller cards'); break; }
    case 'FileContextMenu': { assert(one(root,'[data-ws-file-path]').textContent==='docs/brief.md','context path missing'); hostButton(root,actions,'rename-file','docs/brief.md'); break; }
    case 'RecentFiles': { const item=one(root,'[data-ws-recent-file]'); assert(item.textContent.includes('brief.md') && one(item,'time').textContent==='2026-10-11' && item.textContent.includes('4 KB'),'recent file metadata missing'); break; }
    case 'Favorites': { const box=one(root,'input[data-ws-favorite]'); assert(box.value==='docs/brief.md' && box.checked,'caller favorite missing'); box.click(); assert(!box.checked,'favorite checkbox did not update locally'); break; }
    case 'DuplicateFinder': { const group=one(root,'[data-ws-duplicate-group]'); assert(group.dataset.hash==='sha256:caller-duplicate' && group.querySelectorAll('li').length===2,'duplicate group missing caller files'); hostButton(root,actions,'delete-duplicate','docs/brief-copy.md'); assert(group.isConnected && group.textContent.includes('archive/brief.md'),'duplicate was removed without host result'); break; }
    case 'SchemaTree': { const parent=one(root,'[data-schema-path="public"] details'); assert(!parent.open,'schema starts collapsed'); one(parent,'summary').click(); const leaf=one(parent,'[data-schema-path="public.projects.id"]'); assert(parent.open && leaf.textContent.includes('uuid'),'recursive schema child missing'); break; }
    case 'TableStructureEditor': { const name=one(root,'input[name="columns[0][name]"]'); name.value='project_uuid'; const nullable=one(root,'input[name="columns[0][nullable]"]'); nullable.click(); const fd=submitData(root,actions,'columns[0][name]','project_uuid'); assert(fd.get('columns[0][type]')==='uuid' && fd.has('columns[0][nullable]')===nullable.checked,'table structure payload is wrong');nullable.click();const unchecked=submitData(root,actions,'columns[0][name]','project_uuid');assert(!nullable.checked&&!unchecked.has('columns[0][nullable]'),'unchecked nullable remained in submitted fields'); break; }
    case 'RedisKeyBrowser': { const input=one(root,'input[data-ws-redis-search]'), entries=[...root.querySelectorAll('[data-ws-redis-entry]')]; input.value='project'; input.dispatchEvent(new Event('input',{bubbles:true})); assert(entries[0].hidden && !entries[1].hidden,'redis key search did not filter'); const key=one(entries[1],'[data-ws-redis-key]'); key.click(); assert(one(root,'[data-ws-redis-preview]').textContent==='{"name":"Caller project"}' && actions.findLast(a=>a.action==='select-redis-key')?.key==='project:7','redis selection payload or preview wrong'); hostButton(root,actions,'delete-redis-key','project:7'); assert(!entries[1].hidden && entries[1].isConnected,'redis key was deleted without host'); break; }
    case 'DocumentViewer': { const doc=one(root,'[data-ws-devtools-document]'); assert(!doc.querySelector('textarea,[contenteditable="true"]') && one(doc,'[data-ws-document-json]').textContent==='{"_id":"project-7","active":true}' && doc.textContent.includes('project-7'),'document viewer is not caller-driven read-only JSON'); break; }
    case 'ContainerList': { const row=one(root,'tr[data-ws-container]'); assert(row.dataset.containerId==='a1b2c3d4' && row.textContent.includes('caller/api:1.2') && row.textContent.includes('8080:8080'),'container metadata missing'); hostButton(root,actions,'stop-container','a1b2c3d4'); assert(row.textContent.includes('running'),'container status changed without host'); break; }
    case 'GraphQLExplorer': { const query=one(root,'textarea[name="query"]'), variables=one(root,'textarea[name="variables"]'); query.value='query { project { id } }'; variables.value='{"id":"changed"}'; const fd=submitData(root,actions,'query','query { project { id } }'); assert(fd.get('variables')==='{"id":"changed"}','GraphQL variables missing from host payload'); root.dispatchEvent(new CustomEvent('ely:result',{bubbles:true,detail:{action:one(root,'form').dataset.elySubmit,result:{data:{project:{id:'changed'}}}}})); const output=one(root,'[data-ws-host-result]'); assert(!output.hidden && output.dataset.wsResultState==='result' && output.textContent==='{"data":{"project":{"id":"changed"}}}','GraphQL host result did not render safely'); break; }
    case 'RegexTester': { const pattern=one(root,'input[name="pattern"]'), text=one(root,'textarea[name="text"]'); pattern.value='caller-[0-9]+'; text.value='caller-77'; const fd=submitData(root,actions,'pattern','caller-[0-9]+'); assert(fd.get('text')==='caller-77' && !root.querySelector('[data-workspace-rich]'),'regex host payload is incomplete or became generic code'); expectError('Pattern rejected by host');root.dispatchEvent(new CustomEvent('ely:error',{bubbles:true,detail:{action:one(root,'form').dataset.elySubmit,error:new Error('Pattern rejected by host')}})); const output=one(root,'[data-ws-host-result]'); assert(!output.hidden && output.dataset.wsResultState==='error' && output.textContent==='Pattern rejected by host','regex host error did not render'); break; }
    default: throw new Error('No files/devtools assertion for '+kind);
  }
};
const canvas = async ({root,actions}) => {
  const zoom = one(root, '[data-workspace-zoom="in"]'); zoom.click(); assert(one(root, 'output[data-workspace-zoom]').textContent === '110%', 'canvas zoom did not change');
  assert(Math.abs(new DOMMatrix(getComputedStyle(one(root,'[data-workspace-stage]')).transform).a - 1.1) < 0.001, 'visible canvas scale did not change');
  const shape = one(root, '[data-workspace-shape]'); shape.dispatchEvent(new KeyboardEvent('keydown', {key: 'Enter', bubbles: true}));
  assert(actions.findLast(a=>a.action==='select-shape')?.shape==='shape-a','canvas selection lost shape');
  const stage=one(root,'[data-workspace-stage]');
  const pointer=(type,x,y)=>stage.dispatchEvent(new PointerEvent(type,{button:0,clientX:x,clientY:y,bubbles:true,pointerId:9}));
  let rect=stage.getBoundingClientRect();
  one(root,'[data-ws-canvas-mode="pan"]').click();pointer('pointerdown',rect.left+20,rect.top+20);pointer('pointermove',rect.left+50,rect.top+40);pointer('pointerup',rect.left+50,rect.top+40);
  const matrix=new DOMMatrix(getComputedStyle(stage).transform);assert(matrix.e===30 && matrix.f===20,'visible pan failed');
  if(root.dataset.elyKind==='Whiteboard'){
    const matrix=stage.getScreenCTM(),start=new DOMPoint(100,100).matrixTransform(matrix),end=new DOMPoint(180,180).matrixTransform(matrix);
    one(root,'[data-ws-canvas-mode="draw"]').click();pointer('pointerdown',start.x,start.y);pointer('pointermove',end.x,end.y);pointer('pointerup',end.x,end.y);
    const stroke=one(root,'[data-ws-strokes] polyline'),entry=actions.findLast(a=>a.action==='draw-stroke');
    assert(stroke.getAttribute('points').split(' ').length===2 && entry?.points.length===2 && entry.points[0].x!==entry.points[1].x && entry.points.every(point=>point.x>0&&point.x<800&&point.y>0&&point.y<420),'visible freehand geometry or payload absent');
  }
};
const request = async ({root, actions}) => {
  const method=one(root,'select[name="method"]');
  assert(method.value==='POST' && [...method.options].some(option=>option.value==='GET'),'caller method not selected');
  method.value='GET';method.dispatchEvent(new Event('change',{bubbles:true}));
  const endpoint = one(root, 'input[name="endpoint"]'); endpoint.value = '/caller/changed';
  const body=one(root,'textarea[name="body"]');body.value='{"caller":"changed"}';
  const fd=submitData(root,actions,'endpoint','/caller/changed');
  assert(fd.get('method')==='GET' && fd.get('body')==='{"caller":"changed"}','request payload is stale or incomplete');
};
const fileInput = async ({root}) => { const input = one(root, 'input[type="file"]'); assert(input.multiple, 'file input must accept a FileList'); assert(input.name === 'files', 'file input has no host payload name'); };
const submitData = (root, actions, name, expected) => {
  const form=[...root.querySelectorAll('form')].find(form=>form.elements.namedItem(name));assert(form,'missing form field '+name);
  one(form,'button[type="submit"]').click();
  const entry=actions.findLast(a=>a.action===form.dataset.elySubmit);
  assert(entry?.formData instanceof FormData && entry.formData.get(name)===expected,'wrong submitted '+name);
  return entry.formData;
};
const hostButton = (root,actions,name,value) => {
  const button=[...root.querySelectorAll('[data-ely-action="'+name+'"]')].find(button=>button.value===value);assert(button,'missing '+name+' for '+value);button.click();
  assert(actions.findLast(a=>a.action===name)?.value===value,'wrong '+name+' payload');
};
const checkComment = (root,actions) => {
  const comment=one(root,'[data-ws-comment]');
  assert(comment.textContent.includes('Caller reviewer') && comment.textContent.includes('Explain caller change'),'caller comment metadata missing');
  hostButton(comment,actions,'resolve-comment','comment-8');
  const resolved=actions.findLast(a=>a.action==='resolve-comment');
  assert(resolved.formData.get('comment')==='comment-8' && one(comment,'[data-ws-comment-status]').textContent==='Open','resolve lost payload or invented success');
  one(comment,'input[name="reply"]').value='Caller reply';
  assert(submitData(comment,actions,'reply','Caller reply').get('comment')==='comment-8','reply lost parent comment');
};
const semantic = async ({root,actions,twin}) => {
  const kind=root.dataset.elyKind,variant=root.dataset.elyComponent.split('/')[1];
  switch(kind){
    case 'CompletionMenu':{
      const filter=one(root,'input[data-ws-candidate-filter]'),rows=[...root.querySelectorAll('[data-ws-candidate]')];assert(rows.length===2,'caller completion candidates absent');
      filter.value='readCaller';filter.dispatchEvent(new Event('input',{bubbles:true}));assert(rows[0].hidden&&!rows[1].hidden,'completion filtering failed');
      filter.value='';filter.dispatchEvent(new Event('input',{bubbles:true}));
      const source=one(root,'textarea[data-ws-completion-source]');source.setSelectionRange(0,3);one(root,'[data-ws-complete="complete-print"]').click();
      const a=actions.findLast(a=>a.action==='select-completion');assert(source.value==='printCaller()' && a?.source===source.value && a?.candidate==='complete-print','completion insertion/payload failed');break;
    }
    case 'CodeActionMenu':
      hostButton(root,actions,'apply-code-action','action-extract');{const fd=actions.findLast(a=>a.action==='apply-code-action').formData;assert(fd.get('path')==='src/worker.go'&&fd.get('source')==='callerSource'&&fd.get('command')==='refactor.extract','code action payload incomplete');}break;
    case 'ProblemsPanel':
      assert(root.textContent.includes('warning')&&root.textContent.includes('Caller value is unused'),'diagnostic missing');hostButton(root,actions,'open-location','diag-1');{const fd=actions.findLast(a=>a.action==='open-location').formData;assert(fd.get('file')==='src/worker.go'&&fd.get('line')==='27'&&fd.get('column')==='8','diagnostic location incomplete');}break;
    case 'WhitespaceRenderer':
      assert(one(root,'pre[data-ws-whitespace]').textContent==='a·b→c↵\n'&&!root.querySelector('textarea,[contenteditable="true"]'),'whitespace markers or readonly semantics missing');break;
    case 'DebugConsole':
      one(root,'input[name="expression"]').value='request.count';assert(submitData(root,actions,'expression','request.count').get('session')==='session-12','console lost session');assert(one(root,'[data-ws-console-output]').textContent==='Caller debugger output','console invented result');break;
    case 'SettingsEditor':{
      const size=one(root,'input[name="fontSize"]');size.value='2';const count=actions.length;one(root,'button[type="submit"]').click();assert(actions.length===count&&!size.validity.valid,'setting bounds not enforced');size.value='18';assert(submitData(root,actions,'fontSize','18').get('language')==='go','setting payload incomplete');break;
    }
    case 'WelcomePage':
      assert(one(root,'nav[aria-label="Projects"] a').getAttribute('href')==='#caller-project'&&root.textContent.includes('Caller workspace'),'caller project link absent');hostButton(root,actions,'run-task','task-check');break;
    case 'CursorPosition':{
      const editor=one(root,'[data-workspace-rich]');assert(editor.textContent==='first\nsecond','cursor caller content missing');const range=document.createRange();range.setStart(editor.firstChild,8);range.collapse(true);const selection=document.getSelection();selection.removeAllRanges();selection.addRange(range);document.dispatchEvent(new Event('selectionchange'));assert(one(root,'[data-workspace-cursor]').textContent==='Ln 2, Col 3','cursor selection was not tracked');await code({root,actions,twin});break;
    }
    case 'PullRequestCard':
      assert(root.textContent.includes('#42 Caller review')&&root.textContent.includes('feature/worker')&&root.textContent.includes('Reviewer one'),'PR metadata incomplete');hostButton(root,actions,'request-review','pr-42');assert(actions.findLast(a=>a.action==='request-review').formData.get('pullRequest')==='pr-42'&&one(root,'[data-ws-pr-state]').textContent.includes('Open'),'PR payload missing or fictional status');if(variant==='ReviewComment')checkComment(root,actions);break;
    case 'Block types':{
      one(root,'[data-ws-insert-block="h2"]').click();assert(one(root,'[data-workspace-rich] h2').textContent==='Caller inserted heading','block was not inserted');const editor=one(root,'[data-workspace-rich]');submitData(root,actions,editor.dataset.name,'Caller body<h2>Caller inserted heading</h2>');checkRequiredEditor(root,actions,editor);break;
    }
    case 'Footnote':
      if(variant==='Footnote'){const link=one(root,'[data-ws-footnote-link]');assert(link.hash==='#'+root.id+'-note'&&one(root,'[id="'+root.id+'-note"]').textContent==='Caller footnote definition','footnote target wrong');link.click();await waitFor(()=>location.hash===link.hash);}
      else{const details=one(root,'details[data-ws-glossary]');one(details,'summary').click();assert(details.open&&one(details,'p').textContent==='Caller footnote definition','glossary definition disclosure failed');}break;
    case 'FocusMode':
      one(root,'[data-ws-focus-toggle]').click();assert(one(root,'[data-ws-focus-chrome]').hidden&&one(root,'[data-ws-focus-toggle]').getAttribute('aria-pressed')==='true'&&document.activeElement===one(root,'[data-ws-focus-content]'),'focus failed');one(root,'[data-ws-focus-toggle]').click();assert(!one(root,'[data-ws-focus-chrome]').hidden,'focus did not exit');break;
    case 'PrintPreview':{
      const pages=[...root.querySelectorAll('[data-ws-print-page]')];assert(pages.length===2&&pages[1].textContent.includes('Caller page two')&&!root.querySelector('[contenteditable="true"]'),'readonly print pages missing');const fd=submitData(root,actions,'document','doc-21');assert(fd.get('paper')==='A4'&&fd.get('orientation')==='portrait'&&fd.getAll('page').join('|')==='Caller page one|Caller page two','print payload incomplete');break;
    }
    case 'PageTree':
      await files({root});assert(one(root,'nav[aria-label="Page tree"]').textContent.includes('Caller lift page'),'recursive caller pages absent');break;
    case 'PageCover':{
      const image=one(root,'img[data-ws-cover-image]');assert(image.getAttribute('src')==='/assets/sample.svg'&&image.alt==='Caller cover image'&&one(root,'dd').textContent==='Draft','cover data missing');await waitFor(()=>image.complete);assert(image.naturalWidth>0,'cover image failed to display');break;
    }
    case 'Backlinks':
      assert(one(root,'nav[aria-label="Backlinks"] a').getAttribute('href')==='#caller-source'&&root.textContent.includes('Caller linked passage'),'backlink source missing');if(variant==='GraphView'){const node=one(root,'svg g[data-workspace-select]');node.dispatchEvent(new MouseEvent('click',{bubbles:true}));assert(node.getAttribute('aria-pressed')==='true','backlink graph selection failed');}break;
    case 'TemplatePicker':
      hostButton(root,actions,'select-template','tpl-brief');assert(actions.findLast(a=>a.action==='select-template').formData.get('content')==='Caller template content','template content payload missing');break;
    case 'PresenceAvatars':
      assert(one(root,'[role="img"]').getAttribute('aria-label')==='Caller writer'&&root.textContent.includes('Editing')&&one(root,'mark').textContent==='Caller selected phrase','presence semantics missing');hostButton(root,actions,'follow-user','user-7');break;
    case 'CommentBubble':case 'CommentSidebar':case 'Annotation':
      checkComment(root,actions);if(kind==='Annotation')assert(one(root,'mark').textContent==='Caller passage','annotation quote missing');break;
    case 'SuggestionMode':
      assert(one(root,'del').textContent==='Old caller text'&&one(root,'ins').textContent==='New caller text','suggestion change missing');for(const name of ['accept-suggestion','reject-suggestion']){hostButton(root,actions,name,'suggestion-9');const fd=actions.findLast(a=>a.action===name).formData;assert(fd.get('original')==='Old caller text'&&fd.get('replacement')==='New caller text','suggestion payload incomplete');}assert(one(root,'[data-ws-suggestion-status]').textContent==='Pending','suggestion invented provider result');break;
    case 'ActivityLog':
      assert(one(root,'ol[aria-label="Activity timeline"] time').dateTime==='2026-10-11T10:40:00+09:00'&&one(root,'ol li a').getAttribute('href')==='#caller-notes'&&root.textContent.includes('Caller writer')&&root.textContent.includes('edited'),'caller timeline semantics missing');break;
    default:throw new Error('Missing semantic test for '+kind);
  }
};
const domain = async ({root,actions}) => {
  const kind=root.dataset.elyKind, variant=root.dataset.elyComponent.split('/')[1];
  switch(kind) {
    case 'RenameInline': {
      const input=one(root,'input[data-ws-rename-value]'); input.value='changed.md'; one(root,'[data-ws-rename-cancel]').click();
      assert(input.value==='brief.md' && one(root,'[data-ws-status]').textContent==='Edit canceled','cancel did not restore name');
      input.value='revised.md'; assert(submitData(root,actions,'name','revised.md').get('path')==='docs/brief.md','rename lost source path'); break;
    }
    case 'FileSearch': case 'SearchPanel': {
      const input=one(root,'input[data-ws-search]'), rows=[...root.querySelectorAll('[data-ws-search-row]')];
      assert(rows.length===2,'caller search rows missing'); input.value='worker'; input.dispatchEvent(new Event('input',{bubbles:true}));
      assert(!rows[0].hidden && rows[1].hidden && one(root,'[data-ws-search-count]').textContent==='1 results','path filter failed');
      input.value='no-such-file'; input.dispatchEvent(new Event('input',{bubbles:true})); assert(rows.every(r=>r.hidden),'empty filter failed'); break;
    }
    case 'FilePreview':
      assert(one(root,'[data-ws-preview] pre').textContent==='Read-only caller file <script>content</script>','file content not escaped'); assert(!root.querySelector('script,[contenteditable="true"],textarea'),'file preview editable or executable'); break;
    case 'HexViewer': case 'MemoryViewer': {
      const rows=[...root.querySelectorAll('[data-ws-hex] tbody tr')]; assert(rows.length===4,'byte count wrong'); assert(rows[0].cells[0].textContent==='00000010','offset missing'); assert(rows.map(r=>r.cells[1].textContent).join(' ')==='00 41 7f ff','hex encoding wrong'); assert(!root.querySelector('[contenteditable="true"],textarea'),'hex must be read-only'); break;
    }
    case 'FileOperationProgress': { const p=one(root,'progress[data-ws-progress]'); assert(p.value===37 && p.max===80 && root.textContent.includes('Running'),'operation progress wrong'); break; }
    case 'TransferQueue': case 'DownloadManager': {
      const p=one(root,'progress[data-ws-progress]'); assert(p.value===32 && p.max===80,'transfer progress wrong'); hostButton(root,actions,kind==='TransferQueue'?'resume-transfer':'cancel-transfer','transfer-7'); assert(root.textContent.includes('Paused') && root.textContent.includes('report.pdf'),'provider state changed without host'); break;
    }
    case 'StorageUsage': { const m=one(root,'meter[data-ws-storage]'); assert(m.value===240 && m.max===500,'storage meter wrong'); assert(one(root,'dl dd').textContent==='180 GB','storage category missing'); break; }
    case 'DebugToolbar': for(const a of ['debug-run','debug-pause','debug-step','debug-stop'])hostButton(root,actions,a,'session-12'); assert(one(root,'[data-ws-debug-state]').textContent==='Paused','debugger claimed provider result'); break;
    case 'BreakpointList': {
      const box=one(root,'input[data-ws-breakpoint]'); box.click(); const a=actions.findLast(a=>a.action==='toggle-breakpoint');
      assert(!box.checked && a?.enabled===false && a?.file==='src/worker.go' && a?.line===27 && a?.breakpoint==='bp-4','breakpoint payload wrong');
      const frame=one(root,'ol[aria-label="Call stack"] button'); frame.click(); assert(frame.getAttribute('aria-pressed')==='true' && frame.textContent.includes('HandleRequest'),'stack selection failed'); break;
    }
    case 'VariablesPanel': {
      const detail=one(root,'details'); assert(one(detail,'dd code').textContent==='4' && root.textContent.includes('attempts'),'recursive variable missing');
      one(detail,'summary').click(); assert(!detail.open,'variable disclosure failed');
      if(variant==='ConsoleREPL'){one(root,'input[name="expression"]').value='request.attempts';submitData(root,actions,'expression','request.attempts');} break;
    }
    case 'PerformanceProfiler': {
      const bars=[...root.querySelectorAll('meter[data-ws-duration]')]; assert(bars.length===2 && bars[0].value===40 && bars[1].value===120 && bars.every(b=>b.max===200),'duration bars wrong'); assert(root.textContent.includes('ParseRequest'),'sample label missing'); break;
    }
    case 'NetworkInspector': {
      const details=one(root,'details[data-ws-network]'); assert(one(details,'summary').textContent.includes('POST /caller/notes') && details.textContent.includes('202') && details.textContent.includes('18 ms'),'network metadata missing');
      one(details,'summary').click(); assert(details.open && one(details,'pre').textContent==='{"accepted":true}','response disclosure failed'); break;
    }
    case 'ProcessList': {
      const cells=one(root,'tbody tr').cells; assert(cells[0].textContent==='4242' && cells[2].textContent==='12.5' && cells[3].textContent==='48 MB','process data missing'); hostButton(root,actions,'terminate-process','4242'); assert(root.textContent.includes('caller-worker'),'terminated without host'); break;
    }
    case 'AnsiText': case 'LogViewer': {
      const rows=[...root.querySelectorAll('[data-ws-search-row]')]; assert(rows.length===2 && rows[0].textContent.includes('10:32'),'caller log missing');
      const input=one(root,'input[data-ws-search]'); input.value='error'; input.dispatchEvent(new Event('input',{bubbles:true})); assert(rows[0].hidden && !rows[1].hidden && rows[1].textContent.includes('caller retry'),'log filter failed'); assert(!root.querySelector('[contenteditable="true"]'),'logs must be read-only'); break;
    }
    case 'BranchSelector': {
      const input=one(root,'select[data-ws-branch]'); input.value='feature/worker'; input.dispatchEvent(new Event('change',{bubbles:true})); assert(one(root,'[data-ws-branch-label]').textContent==='feature/worker' && actions.findLast(a=>a.action==='switch-branch')?.branch==='feature/worker','branch payload wrong'); break;
    }
    case 'ChangesList':
      assert(root.textContent.includes('Modified') && root.textContent.includes('docs/brief.md'),'caller changes missing'); hostButton(root,actions,'stage-change','src/worker.go');
      if(variant==='CommitInput'){one(root,'textarea[name="message"]').value='Update worker';submitData(root,actions,'message','Update worker');} break;
    case 'TodoItem': {
      const box=one(root,'input[data-ws-todo]');box.click(); assert(!box.checked && one(root,'[data-ws-todo-count]').textContent==='0 / 2 completed','todo count wrong'); break;
    }
    case 'DocumentViewer': case 'OfficePreview': case 'EpubReader': {
      assert(!root.querySelector('[contenteditable="true"],textarea'),'document editable'); one(root,'button[data-ws-page="1"]').click(); const page=one(root,'[data-ws-page-content="1"]');
      assert(!page.hidden && one(page,'pre').textContent==='Caller document details' && one(root,'[data-ws-page-content="0"]').hidden,'chapter navigation failed');
      if(kind==='OfficePreview')assert(one(root,'tbody td').textContent==='Caller' && one(root,'tbody tr').cells[1].textContent==='1200','office cells wrong'); break;
    }
    case 'VersionHistory': one(root,'button[data-ws-revision="1"]').click(); assert(!one(root,'[data-ws-revision-content="1"]').hidden && one(root,'[data-ws-revision-content="1"]').textContent==='Updated caller revision' && one(root,'[data-ws-revision-content="0"]').hidden,'revision navigation failed'); break;
    case 'TrashBin': hostButton(root,actions,'restore-file','deleted-8'); assert(root.textContent.includes('docs/retired.md'),'restore claimed success without provider'); break;
    case 'ToolPalette': {
      const input=one(root,'select[data-ws-tool]');input.value='rectangle';input.dispatchEvent(new Event('change',{bubbles:true})); assert(root.dataset.workspaceTool==='rectangle' && one(root,'[data-ws-tool-label]').textContent==='rectangle','tool selection failed'); break;
    }
    case 'LayerPanel': {
      const input=one(root,'input[data-ws-layer]');input.click();assert(!input.checked && one(root,'[data-workspace-shape][data-value="shape-a"]').style.display==='none','visibility failed'); break;
    }
    case 'ColorPanel': {
      const input=one(root,'input[type="color"][data-ws-color]');input.value='#ff9900';input.dispatchEvent(new Event('input',{bubbles:true}));assert(one(root,'[data-ws-color-preview]').getAttribute('fill')==='#ff9900' && one(root,'[data-ws-color-label]').textContent==='#ff9900','color preview failed');break;
    }
    case 'ConstraintsEditor': {
      one(root,'input[name="width"]').value='640';const fd=submitData(root,actions,'width','640');assert(fd.get('height')==='180' && fd.get('gap')==='16' && fd.has('autoLayout'),'constraints payload incomplete');if(variant==='ExportPanel')hostButton(root,actions,'export-selection','');break;
    }
    case 'NodeEditor': {
      const node=one(root,'[data-workspace-shape][data-value="node-source"]');assert(node.getAttribute('transform')==='translate(40 70)' && one(node,'rect').getAttribute('width')==='130' && one(root,'[data-workspace-edge]').getAttribute('d')==='M170 100 L300 100','node geometry wrong');node.dispatchEvent(new KeyboardEvent('keydown',{key:'Enter',bubbles:true}));assert(actions.findLast(a=>a.action==='select-shape')?.shape==='node-source','node selection missing');
      node.dispatchEvent(new PointerEvent('pointerdown',{button:0,clientX:40,clientY:70,bubbles:true,pointerId:7}));
      node.dispatchEvent(new PointerEvent('pointermove',{clientX:60,clientY:90,bubbles:true,pointerId:7}));
      node.dispatchEvent(new PointerEvent('pointerup',{bubbles:true,pointerId:7}));
      assert(node.getAttribute('transform')==='translate(60 90)','node movement failed');break;
    }
    case 'ArchiveViewer': assert(one(root,'[data-ws-document] pre').textContent==='Caller rich text content' && !root.querySelector('[contenteditable="true"]'),'archive manifest must be read-only');break;
    default: throw new Error('No domain assertion for '+kind);
  }
};
const domainChecks=Object.fromEntries([
  'files/RenameInline','files/FileSearch','editor/SearchPanel','files/FilePreview','files/HexViewer','debug/MemoryViewer',
  'files/FileOperationProgress','files/TransferQueue','files/DownloadManager','files/StorageUsage','debug/DebugToolbar','debug/BreakpointList','debug/VariablesPanel','debug/PerformanceProfiler','debug/NetworkInspector','terminal/ProcessList','terminal/AnsiText','terminal/LogViewer','git/BranchSelector','git/ChangesList','documents/TodoItem','documents/DocumentViewer','documents/OfficePreview','documents/EpubReader','documents/VersionHistory','documents/TrashBin','canvas/ToolPalette','canvas/LayerPanel','canvas/ColorPanel','canvas/ConstraintsEditor','canvas/NodeEditor','files/ArchiveViewer'
].map(key=>[key,domain]));
const make = (names, check) => Object.fromEntries(names.map(name => [name, check]));
registerChecks('workspace', {
  ...make(["editor/CompletionMenu","editor/CodeActionMenu","editor/ProblemsPanel","editor/WhitespaceRenderer","editor/DebugConsole","editor/SettingsEditor","editor/WelcomePage","editor/CursorPosition","git/PullRequestCard","documents/Block types","documents/Footnote","documents/FocusMode","documents/PrintPreview","documents/PageTree","documents/PageCover","documents/Backlinks","documents/TemplatePicker","collab/PresenceAvatars","collab/CommentBubble","collab/CommentSidebar","collab/Annotation","collab/SuggestionMode","collab/ActivityLog"], semantic),
  ...make(['editor/CodeEditor','devtools/QueryEditor'], code),
  ...make(['editor/InlineDiff'], diff),
  ...make(['editor/FindWidget'], find), ...make(['git/DiffViewer'], diff), ...make(['git/ThreeWayMerge'], conflict), ...make(['git/CommitGraph'], graph), ...make(['git/BlameView'], blame),
  ...make(['terminal/Terminal','terminal/CommandBlock'], terminal), ...make(['documents/RichTextEditor','documents/BlockEditor','collab/CommentThread'], rich), ...make(['documents/MarkdownEditor'], markdown), ...make(['documents/MarkdownRenderer'], markdownRenderer),
  ...make(['files/FileExplorer','files/FolderTree','files/FileList'], files), ...make(['files/DragDropFiles'], fileInput), ...make(['canvas/InfiniteCanvas','canvas/FlowChart','canvas/MindMap','canvas/Whiteboard'], canvas),
  ...make(['collab/ShareDialog','devtools/ConnectionManager','devtools/ApiRequestBuilder','devtools/WebSocketConsole'], request),
  ...make(['files/FileIcon','files/FileItem','files/FileGrid','files/FileContextMenu','files/RecentFiles','files/Favorites','files/DuplicateFinder','devtools/SchemaTree','devtools/TableStructureEditor','devtools/RedisKeyBrowser','devtools/DocumentViewer','devtools/ContainerList','devtools/GraphQLExplorer','devtools/RegexTester'], filesDevtools)
  ,...domainChecks
});
