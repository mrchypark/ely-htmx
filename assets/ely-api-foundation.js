import {register, emit} from './ely-api.js';

const selectText = root => {
  const value = root.querySelector('[data-foundation-copy]');
  if (!value) return '';
  const selection = getSelection();
  const range = document.createRange();
  range.selectNodeContents(value);
  selection?.removeAllRanges();
  selection?.addRange(range);
  return value.textContent;
};

const format = node => {
  const value = node.getAttribute('datetime') || node.textContent.trim();
  const number = Number(value);
  if (node.dataset.foundationFormat === 'number' && Number.isFinite(number)) node.textContent = new Intl.NumberFormat().format(number);
  if (node.dataset.foundationFormat === 'currency' && Number.isFinite(number)) node.textContent = new Intl.NumberFormat(undefined, {style: 'currency', currency: node.dataset.foundationCurrency || 'USD'}).format(number);
  if (node.dataset.foundationFormat === 'percent' && Number.isFinite(number)) node.textContent = new Intl.NumberFormat(undefined, {style: 'percent', maximumFractionDigits: 2}).format(number);
  if (node.dataset.foundationFormat === 'filesize' && Number.isFinite(number)) { const units = ['B', 'KB', 'MB', 'GB', 'TB']; let unit = 0, size = number; while (Math.abs(size) >= 1024 && unit < units.length - 1) { size /= 1024; unit++; } node.textContent = `${new Intl.NumberFormat(undefined, {maximumFractionDigits: 1}).format(size)} ${units[unit]}`; }
  if (node.dataset.foundationFormat === 'duration' && Number.isFinite(number)) { const seconds = Math.max(0, Math.floor(number)); node.textContent = [seconds / 3600, (seconds % 3600) / 60, seconds % 60].map((part, index) => String(Math.floor(part)).padStart(index ? 2 : 1, '0')).join(':'); }
  if (node.dataset.foundationFormat === 'plural' && Number.isFinite(number)) node.textContent = `${number} ${node.dataset.foundationUnit || 'item'}${number === 1 ? '' : 's'}`;
  const date = new Date(value);
  if (node.dataset.foundationFormat === 'datetime' && !Number.isNaN(date)) node.textContent = new Intl.DateTimeFormat(undefined, {dateStyle: 'medium', timeStyle: 'short'}).format(date);
  if (node.dataset.foundationFormat === 'relative' && !Number.isNaN(date)) { const seconds = Math.round((date - Date.now()) / 1000); const units = [['year', 31536000], ['month', 2592000], ['day', 86400], ['hour', 3600], ['minute', 60], ['second', 1]]; const [unit, amount] = units.find(([, size]) => Math.abs(seconds) >= size) || units.at(-1); node.textContent = new Intl.RelativeTimeFormat(undefined, {numeric: 'auto'}).format(Math.round(seconds / amount), unit); }
};

const prefersReducedMotion = () => globalThis.matchMedia?.('(prefers-reduced-motion: reduce)').matches;
const animationDuration = node => Math.max(0, Math.min(10000, Number(node.dataset.foundationDuration) || 300));

const drawCanvas = canvas => {
  const context = canvas.getContext?.('2d');
  const values = (canvas.dataset.foundationCanvas.match(/-?\d+(?:\.\d+)?/g) || []).slice(0, 400).map(Number).filter(Number.isFinite);
  if (!context || values.length < 4) return;
  context.strokeStyle = 'currentColor';
  context.beginPath();
  for (let index = 0; index + 1 < values.length; index += 2) {
    if (index === 0) context.moveTo(values[index], values[index + 1]);
    else context.lineTo(values[index], values[index + 1]);
  }
  context.stroke();
};

const animateNumber = (node, cleanups) => {
  const end = Number(node.dataset.foundationNumber), start = Number(node.dataset.foundationFrom);
  if (!Number.isFinite(end)) return;
  const from = Number.isFinite(start) ? start : 0;
  if (node.getAttribute('aria-disabled') === 'true' || prefersReducedMotion() || animationDuration(node) === 0) { node.textContent = String(end); return; }
  const started = performance.now(), total = animationDuration(node); let frame = 0;
  const tick = now => { const progress = Math.min(1, (now - started) / total); node.textContent = String(Math.round(from + (end - from) * progress)); if (progress < 1) frame = requestAnimationFrame(tick); };
  frame = requestAnimationFrame(tick); cleanups.push(() => cancelAnimationFrame(frame));
};

const typewrite = (node, cleanups) => {
  const text = node.dataset.foundationTypewriter || '';
  if (node.getAttribute('aria-disabled') === 'true' || prefersReducedMotion() || animationDuration(node) === 0) { node.textContent = text; return; }
  const started = performance.now(), total = animationDuration(node); let frame = 0;
  const tick = now => { const progress = Math.min(1, (now - started) / total); node.textContent = text.slice(0, Math.floor(text.length * progress)); if (progress < 1) frame = requestAnimationFrame(tick); else node.textContent = text; };
  frame = requestAnimationFrame(tick); cleanups.push(() => cancelAnimationFrame(frame));
};

register('foundation', root => {
  if (root.dataset.foundationMounted) return;
  root.dataset.foundationMounted = 'true';
  const cleanups = [];
  root.querySelectorAll('[data-foundation-canvas]').forEach(drawCanvas);
  root.querySelectorAll('[data-foundation-number]').forEach(node => animateNumber(node, cleanups));
  root.querySelectorAll('[data-foundation-typewriter]').forEach(node => typewrite(node, cleanups));
  root.querySelectorAll('[data-foundation-format]').forEach(format);
  root.querySelectorAll('[data-foundation-measure]').forEach(node => {
    if (!globalThis.ResizeObserver) return;
    const draw = () => { node.textContent = `${Math.round(node.getBoundingClientRect().width)} px`; };
    draw();
    const observer = new ResizeObserver(draw); observer.observe(node);
    cleanups.push(() => observer.disconnect());
  });
  const onChange = event => {
    const show = event.target.closest('[data-foundation-show]');
    if (show) { const target = root.querySelector('[data-foundation-conditional]'); if (target) target.hidden = !show.checked; emit(root, 'change', {name: show.name, value: show.checked}); }
    const details = event.target.closest('details');
    if (details) emit(root, 'change', {open: details.open});
    const range = event.target.closest('[data-foundation-split] input[type=range]');
    if (range) { range.closest('[data-foundation-split]').style.setProperty('--foundation-split', `${range.value}%`); emit(root, 'change', {value: range.value}); }
  };
  const onClick = event => {
    const menu = event.target.closest('[data-foundation-menu]');
    if (menu) { const panel = menu.parentElement.querySelector('[data-foundation-menu-panel]'); const open = panel?.hidden; if (panel) panel.hidden = !open; menu.setAttribute('aria-expanded', String(open)); emit(root, 'change', {open}); return; }
    if (!event.target.closest('.ely-foundation__menu')) root.querySelectorAll('[data-foundation-menu-panel]').forEach(panel => { panel.hidden = true; });
    const open = event.target.closest('[data-foundation-open]');
    if (open) { const scope = root.querySelector(`#${CSS.escape(open.dataset.foundationOpen)}`); if (scope) { scope.hidden = false; open.setAttribute('aria-expanded', 'true'); scope.focus(); } return; }
    const close = event.target.closest('[data-foundation-close]');
    if (close) { const scope = close.closest('[data-foundation-scope]'); if (scope) { scope.hidden = true; root.querySelector(`[data-foundation-open="${CSS.escape(scope.id)}"]`)?.setAttribute('aria-expanded', 'false'); } return; }
    const dialogOpen = event.target.closest('[data-foundation-dialog-open]');
    if (dialogOpen) { root.querySelector(`#${CSS.escape(dialogOpen.dataset.foundationDialogOpen)}`)?.showModal?.(); return; }
    if (event.target.closest('[data-foundation-dialog-close]')) { event.target.closest('dialog')?.close(); return; }
    if (event.target.closest('[data-foundation-press]')) { const output = root.querySelector('[data-foundation-output]'); if (output) output.value = output.textContent = 'Pressed'; emit(root, 'press'); return; }
    if (event.target.closest('[data-foundation-copy-button]')) { emit(root, 'copy', {value: selectText(root)}); return; }
    if (event.target.closest('[data-foundation-scroll-top]')) { root.querySelector('[data-foundation-scroll]')?.scrollTo({top: 0, behavior: 'smooth'}); emit(root, 'scroll-top'); return; }
    const item = event.target.closest('[data-foundation-master] nav button');
    if (item) { const master = item.closest('[data-foundation-master]'); master.querySelectorAll('nav button').forEach(button => button.setAttribute('aria-current', String(button === item))); master.querySelector('[data-foundation-detail]').textContent = item.textContent; emit(root, 'change', {value: item.textContent}); }
  };
  const onKeydown = event => {
    const keys = event.target.closest('[data-foundation-keys]');
    if (keys && (event.key === 'Enter' || event.key === 'Escape')) { keys.querySelector('output').textContent = event.key; emit(root, 'key', {key: event.key}); }
    const scope = event.target.closest('[data-foundation-scope]');
    if (scope && event.key === 'Tab') {
      const focusable = [...scope.querySelectorAll('button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [href]')];
      const first = focusable[0], last = focusable[focusable.length - 1];
      if (first && last && ((!event.shiftKey && event.target === last) || (event.shiftKey && event.target === first))) {
        event.preventDefault();
        (event.shiftKey ? last : first).focus();
      }
    }
    if (event.key === 'Escape') root.querySelectorAll('[data-foundation-scope]').forEach(scope => { if (!scope.hidden) { scope.hidden = true; root.querySelector(`[data-foundation-open="${CSS.escape(scope.id)}"]`)?.focus(); } });
    const dropzone = event.target.closest('[data-foundation-dropzone]');
    if (dropzone && (event.key === 'Enter' || event.key === ' ')) { event.preventDefault(); dropzone.dataset.active = String(dropzone.dataset.active !== 'true'); emit(root, 'dropzone-focus', {active: dropzone.dataset.active === 'true'}); }
  };
  const onDrag = event => {
    const zone = event.target.closest?.('[data-foundation-dropzone]');
    if (!zone) return;
    event.preventDefault();
    if (event.type === 'dragenter') zone.dataset.active = 'true';
    if (event.type === 'dragleave') delete zone.dataset.active;
    if (event.type === 'drop') { delete zone.dataset.active; emit(root, 'drop', {value: event.dataTransfer?.getData('text/plain') || ''}); }
  };
  root.addEventListener('change', onChange);
  root.addEventListener('click', onClick);
  root.addEventListener('keydown', onKeydown);
  root.addEventListener('dragenter', onDrag);
  root.addEventListener('dragleave', onDrag);
  root.addEventListener('dragover', onDrag);
  root.addEventListener('drop', onDrag);
  const stop = () => { cleanups.splice(0).forEach(cleanup => cleanup()); };
  root.addEventListener('ely:cleanup', stop, {once: true});
  return () => {
    stop();
    root.removeEventListener('ely:cleanup', stop);
    root.removeEventListener('change', onChange);
    root.removeEventListener('click', onClick);
    root.removeEventListener('keydown', onKeydown);
    root.removeEventListener('dragenter', onDrag);
    root.removeEventListener('dragleave', onDrag);
    root.removeEventListener('dragover', onDrag);
    root.removeEventListener('drop', onDrag);
  };
});
