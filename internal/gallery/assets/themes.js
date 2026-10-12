import {createControl} from '/assets/ui-controls.js';

const api = window.ElyTheme;
const preset = document.querySelector('#theme-preset'), name = document.querySelector('#theme-name'), mode = document.querySelector('#theme-mode');
const json = document.querySelector('#theme-json'), message = document.querySelector('[data-theme-message]'), check = document.querySelector('[data-theme-check]'), apply = document.querySelector('[data-theme-apply]');
const labels = {'color-paper':'페이지 배경', 'color-ink':'본문', 'color-muted':'보조 글자', 'color-line':'구분선', 'color-accent':'강조', surface:'표면', border:'입력 경계', 'accent-hover':'강조 호버', 'accent-tint':'강조 배경', 'on-accent':'강조 위 글자', focus:'포커스', draft:'대기', 'draft-tint':'대기 배경', approved:'완료', 'approved-tint':'완료 배경', warning:'주의', 'warning-tint':'주의 배경', error:'오류', 'error-tint':'오류 배경'};
let draft;
const notify = text => { message.textContent = text; };
function update() {
  try {
    for (const preview of document.querySelectorAll('[data-theme-preview]')) api.preview(preview,draft,preview.dataset.themePreview);
    const problems = api.issues(draft), invalid = !name.checkValidity() || [...document.querySelectorAll('[data-theme-hex]')].some(input => !input.checkValidity());
    check.dataset.invalid = String(problems.length > 0 || invalid);
    check.textContent = invalid ? '이름과 HEX 색상 형식을 확인하세요.' : problems.length ? problems.map(p=>`${p.mode} · ${labels[p.foreground]} / ${labels[p.background]}: ${p.ratio.toFixed(2)}:1 (최소 ${p.minimum}:1)`).join('\n') : '라이트·다크의 기본 대비 28개 확인';
    apply.disabled = Boolean(problems.length || invalid);
    json.value = JSON.stringify(draft,null,2);
  } catch (error) { apply.disabled = true; check.dataset.invalid = 'true'; check.textContent = error.message; }
}
function renderFields() {
  document.querySelector('[data-theme-fields]').replaceChildren();
  document.querySelector('[data-theme-status-fields]').replaceChildren();
  for (const [index,key] of api.tokens.entries()) {
    const row=document.createElement('div');row.className='theme-field';
    const label=document.createElement('label');label.htmlFor='theme-color-'+key;label.append(document.createTextNode(labels[key]));
    const token=document.createElement('code');token.textContent=key;label.append(token);
    const hex=createControl('input',{id:label.htmlFor,type:'text',value:draft[mode.value][key],pattern:'#[0-9A-Fa-f]{6}',maxlength:7,required:true,'aria-label':labels[key]+' HEX','data-theme-hex':key});
    const picker=createControl('input',{type:'color',value:draft[mode.value][key],'aria-label':labels[key]+' 색상 선택'});
    hex.addEventListener('input',()=>{if(hex.checkValidity()){draft[mode.value][key]=hex.value.toUpperCase();picker.value=hex.value;}update();});
    picker.addEventListener('input',()=>{draft[mode.value][key]=picker.value.toUpperCase();hex.value=draft[mode.value][key];update();});
    row.append(label,hex,picker);document.querySelector(index<11?'[data-theme-fields]':'[data-theme-status-fields]').append(row);
  }
  update();
}
function choose(id) { draft=api.get(id);draft.id='custom';name.value=draft.name;renderFields(); }
const active=api.current();
if(active.id==='custom'){const option=document.createElement('option');option.value='custom';option.textContent='저장한 제안';preset.append(option);preset.value='custom';}
else if([...preset.options].some(option=>option.value===active.id))preset.value=active.id;
choose(preset.value);
preset.addEventListener('change',()=>{choose(preset.value);notify('제안을 미리 보고 있습니다. 적용 버튼을 누르면 다른 화면에도 사용합니다.');});
mode.addEventListener('change',renderFields);
name.addEventListener('input',()=>{draft.name=name.value;update();});
apply.addEventListener('click',()=>{try {api.register(draft);api.apply(draft.id);notify('이 제안을 적용했습니다. 다른 화면과 새로고침 후에도 같은 설정을 사용합니다.');}catch(error){notify(error.message);}});
document.querySelector('[data-theme-reset]').addEventListener('click',()=>{api.apply('blue');preset.value='blue';choose('blue');notify('기본 Blue 테마를 복원했습니다.');});
document.querySelector('[data-theme-import]').addEventListener('click',()=>{
  try {const imported=JSON.parse(json.value);api.issues(imported);draft=imported;draft.id='custom';name.value=draft.name;renderFields();notify('JSON을 새 제안으로 불러왔습니다. 대비 확인 후 적용하세요.');}
  catch(error){notify('불러오기 실패: '+(error instanceof SyntaxError ? '올바른 JSON 형식인지 확인하세요.' : error.message));}
});
document.querySelector('[data-theme-copy]').addEventListener('click',async()=>{
  try {await navigator.clipboard.writeText(json.value);notify('테마 JSON을 복사했습니다.');}
  catch {json.focus();json.select();notify('복사를 사용할 수 없습니다. 선택한 JSON을 직접 복사하세요.');}
});
