(() => {
  const byID = id => document.getElementById((id || '').replace(/^#/, ''));
  let toastTimer;
  const toast = message => {
    let status = byID('demo-toast');
    if (!status) {
      status = document.createElement('div'); status.id = 'demo-toast';
      status.className = 'demo-toast'; status.setAttribute('role','status');
      document.body.append(status);
    }
    status.textContent = message; status.hidden = false;
    clearTimeout(toastTimer); toastTimer = setTimeout(() => { status.hidden = true; }, 4500);
  };
  let tabGroupCount=0;
  const tab = button => {
    const core=button.hasAttribute('data-core-tab');
    const group=button.closest(core ? '[data-core-tabs]' : '[data-tabs]');
    if (!group) return;
    group.id ||= `demo-tabs-${++tabGroupCount}`;
    const triggers=[...group.querySelectorAll(core ? '[data-core-tab]' : '[data-tab]')];
    const panels=[...group.querySelectorAll(core ? '[data-core-panel]' : '[data-panel]')];
    const value=core ? button.dataset.coreTab : button.dataset.tab;
    triggers.forEach((item,index) => {
      const active=item===button;
      const itemValue=core ? item.dataset.coreTab : item.dataset.tab;
      const panel=panels.find(p => (core ? p.dataset.corePanel : p.dataset.panel)===itemValue);
      item.id ||= `${group.id}-tab-${index}`;
      item.setAttribute('role','tab'); item.setAttribute('aria-selected',String(active));
      item.tabIndex=active ? 0 : -1; item.classList.toggle('is-active',active);
      if (panel) {
        panel.id ||= `${group.id}-panel-${index}`;
        panel.setAttribute('role','tabpanel'); panel.setAttribute('aria-labelledby',item.id);
        item.setAttribute('aria-controls',panel.id); panel.tabIndex=0;
      }
    });
    panels.forEach(panel => { panel.hidden=(core ? panel.dataset.corePanel : panel.dataset.panel)!==value; });
  };
  document.addEventListener('click', async event => {
    const b = event.target.closest('button, [data-dialog-open], [data-toast]');
    if (!b || b.disabled || b.getAttribute('aria-disabled') === 'true') return;
    if (b.dataset.dialogOpen) byID(b.dataset.dialogOpen)?.showModal();
    if (b.hasAttribute('data-dialog-close')) b.closest('dialog')?.close();
    if (b.dataset.clear) { const input = byID(b.dataset.clear); if (input) { if ('value' in input) input.value=''; else input.textContent=''; input.dispatchEvent(new Event('input',{bubbles:true})); input.focus(); } }
    if (b.dataset.toggle) { const target = byID(b.dataset.toggle); if (target) { target.hidden=!target.hidden; b.setAttribute('aria-expanded',String(!target.hidden)); } }
    if (b.dataset.tab || b.dataset.coreTab) tab(b);
    if (b.dataset.toast) toast(b.dataset.toast);
    if (b.dataset.copy) {
      const source=byID(b.dataset.copy);
      if (source) { try { await navigator.clipboard.writeText(source.value ?? source.textContent); toast('복사했습니다.'); } catch { toast('복사하지 못했습니다. 텍스트를 선택해서 복사하세요.'); } }
    }
    if (b.hasAttribute('data-press')) b.setAttribute('aria-pressed',String(b.getAttribute('aria-pressed') !== 'true'));
  });
  document.addEventListener('keydown', event => {
    const b=event.target.closest('[data-tab],[data-core-tab]'); if (!b) return;
    const core=b.hasAttribute('data-core-tab');
    const tabs=[...b.closest(core ? '[data-core-tabs]' : '[data-tabs]').querySelectorAll(core ? '[data-core-tab]' : '[data-tab]')];
    let index=tabs.indexOf(b);
    if (event.key==='ArrowRight' || event.key==='ArrowDown') index=(index+1)%tabs.length;
    else if (event.key==='ArrowLeft' || event.key==='ArrowUp') index=(index+tabs.length-1)%tabs.length;
    else if (event.key==='Home') index=0;
    else if (event.key==='End') index=tabs.length-1;
    else return;
    event.preventDefault(); tabs[index].focus(); tab(tabs[index]);
  });
  document.addEventListener('input', event => {
    if (!event.target.id) return;
    document.querySelectorAll('[data-value-for]').forEach(out => {
      if (out.dataset.valueFor===event.target.id) out.textContent=event.target.value;
    });
  });
  const init = () => {
    document.querySelectorAll('[data-initial-indeterminate]').forEach(input => { input.indeterminate=true; });
    document.querySelectorAll('[data-tabs],[data-core-tabs]').forEach(group => {
      const active=group.querySelector('[role="tab"][aria-selected="true"]') || group.querySelector('[data-tab],[data-core-tab]');
      if (active) tab(active);
    });
  };
  init(); document.addEventListener('htmx:after:swap',init);
})();
