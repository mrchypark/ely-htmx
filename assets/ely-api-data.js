import { register, emit } from './ely-api.js';
import './ely-patterns.js';

const visibleRows = root => [...root.querySelectorAll('tbody tr, .ely-data__list > li')];
const changed = (root, detail = {}) => root.dispatchEvent(new CustomEvent('ely:change', {
  bubbles: true, detail: { id: root.id, component: root.dataset.elyComponent, ...detail },
}));

register('data', root => {
  if (root.dataset.elyDataMounted) return;
  root.dataset.elyDataMounted = 'true';
  root.querySelectorAll('[data-ely-left], [data-ely-width]').forEach(element => {
    const percent = value => Math.max(0, Math.min(100, Number(value) || 0)) + '%';
    if (element.dataset.elyLeft != null) element.style.marginLeft = percent(element.dataset.elyLeft);
    if (element.dataset.elyWidth != null) element.style.width = percent(element.dataset.elyWidth);
  });
  const comparison = root.querySelector('[data-ely-before-after]');
  if (comparison) {
    const range = comparison.querySelector('[data-ely-range]');
    const after = comparison.querySelector('[data-ely-after-image]');
    comparison.style.display = 'grid';
    comparison.querySelectorAll('img').forEach(image => { image.style.gridArea = '1 / 1'; image.style.maxWidth = '100%'; });
    if (after && range) after.style.clipPath = `inset(0 0 0 ${range.value}%)`;
  }
  const swipe = root.querySelector('[data-ely-swipe]');
  if (swipe) {
    const dismiss = () => { swipe.querySelector('ul').hidden = true; swipe.querySelector('[data-ely-swipe-dismiss]').hidden = true; swipe.querySelector('[data-ely-swipe-undo]').hidden = false; emit(root,'dismiss-item',{value:swipe.dataset.value}); };
    let swipeStart;
    swipe.addEventListener('pointerdown',event=>{swipeStart=event.clientX;});
    swipe.addEventListener('pointerup',event=>{if(swipeStart!=null && event.clientX-swipeStart < -60)dismiss();swipeStart=undefined;});
    swipe.querySelector('[data-ely-swipe-dismiss]').addEventListener('click',dismiss);
    swipe.querySelector('[data-ely-swipe-undo]').addEventListener('click',()=>{swipe.querySelector('ul').hidden=false;swipe.querySelector('[data-ely-swipe-dismiss]').hidden=false;swipe.querySelector('[data-ely-swipe-undo]').hidden=true;changed(root,{action:'restore-item'});});
  }
  root.addEventListener('click', event => {
    const select = event.target.closest('[data-ely-select]');
    if (select && !select.disabled) {
      root.querySelectorAll('[data-ely-select][aria-selected="true"]').forEach(item => item.setAttribute('aria-selected', 'false'));
      select.setAttribute('aria-selected', 'true');
      changed(root, { value: select.dataset.value || select.value });
    }
    const move = event.target.closest('[data-ely-move]');
    if (move && !move.disabled) {
      const item = move.closest('[data-ely-sort-item]');
      const next = item?.nextElementSibling;
      if (item && next) item.parentElement.insertBefore(next, item);
      changed(root, { action: 'reorder' });
    }
    const pin = event.target.closest('[data-ely-pin]');
    if (pin) {
      const cells = [...root.querySelectorAll('[data-ely-pinning] tr > :first-child')];
      const active = pin.getAttribute('aria-pressed') !== 'true';
      cells.forEach(cell => { cell.style.position = active ? 'sticky' : ''; cell.style.left = active ? '0' : ''; });
      pin.setAttribute('aria-pressed', String(active));
      changed(root, { action: 'pin-column', pinned: active });
    }
    const dismiss = event.target.closest('[data-ely-dismiss]');
    if (dismiss && !dismiss.disabled) {
      const alert = dismiss.closest('[data-ely-alert]');
      if (alert) alert.hidden = true;
      changed(root, { action: 'acknowledge' });
    }
    if (event.target.closest('[data-ely-carousel-next], [data-ely-carousel-prev]')) {
      const position = root.querySelector('[data-ely-carousel-position]');
      const items = [...root.querySelectorAll('[data-ely-carousel-item]')];
      if (position && items.length) {
        const current = items.findIndex(item => item.hasAttribute('data-ely-current'));
        const next = (Math.max(0, current) + (event.target.closest('[data-ely-carousel-next]') ? 1 : -1) + items.length) % items.length;
        items.forEach((item, index) => { item.hidden = index !== next; item.toggleAttribute('data-ely-current', index === next); });
        position.value = String(next + 1);
      }
      changed(root, { action: 'navigate', position: position?.value });
    }
    const drawing = event.target.closest('[data-ely-drawing]');
    if (drawing) {
      root.querySelectorAll('[data-ely-drawing]').forEach(tool => tool.setAttribute('aria-pressed', String(tool === drawing)));
      changed(root, { action: 'drawing-tool', value: drawing.value });
    }
    const range = event.target.closest('[data-ely-range-select]');
    if (range) emit(root, 'range', { value: range.value });
    if (event.target.closest('[data-ely-realtime]')) {
      emit(root, 'sample-request', { samples: Number(root.querySelector('[data-ely-samples]')?.value || 0) });
    }
    if (event.target.closest('[data-ely-sort]')) {
      const body = root.querySelector('tbody');
      if (body) [...body.rows].sort((a, b) => a.textContent.localeCompare(b.textContent)).forEach(row => body.append(row));
      changed(root, { action: 'sort' });
    }
  });
  root.addEventListener('input', event => {
    if (event.target.matches('[data-ely-range]')) {
      const out = root.querySelector('[data-ely-range-value]');
      if (out) out.value = event.target.value;
      changed(root, { value: event.target.value });
      const after = root.querySelector('[data-ely-after-image]');
      if (after) after.style.clipPath = `inset(0 0 0 ${event.target.value}%)`;
    }
    if (event.target.matches('[data-ely-filter]')) {
      const query = event.target.value.trim().toLocaleLowerCase();
      visibleRows(root).forEach(row => { row.hidden = Boolean(query) && !row.textContent.toLocaleLowerCase().includes(query); });
      changed(root, { value: event.target.value });
    }
    if (event.target.matches('[data-ely-column]')) {
      const column = [...root.querySelectorAll('.ely-data__table th')].findIndex(cell => cell.textContent.trim() === event.target.value);
      if (column >= 0) root.querySelectorAll('.ely-data__table tr').forEach(row => row.cells[column].hidden = !event.target.checked);
      changed(root, { action: 'column', value: event.target.value, checked: event.target.checked });
    }
    if (event.target.matches('[data-ely-pivot-measure]')) {
      root.querySelectorAll('[data-ely-pivot-rows]').forEach(panel => { panel.hidden = panel.dataset.elyPivotRows !== event.target.value; });
      changed(root, { action: 'pivot-measure', measure: event.target.value });
    }
  });
  const applySamples = event => {
    if (event.target !== root) return;
    const detail = event.detail || {};
    const line = root.querySelector('[data-ely-chart] .ely-data__line, .ely-data__chart .ely-data__line');
    const samples = root.querySelector('[data-ely-samples]');
    if (typeof detail.points === 'string' && line) line.setAttribute('points', detail.points);
    if (detail.samples != null && samples) samples.value = String(detail.samples);
    changed(root, { action: 'samples', points: detail.points, samples: detail.samples });
  };
  root.addEventListener('ely:data', applySamples);
  root.addEventListener('ely:samples', applySamples);
  root.addEventListener('ely:data', event => {
    if (event.target !== root || !event.detail?.measure) return;
    root.querySelectorAll('[data-ely-pivot-rows]').forEach(panel => { panel.hidden = panel.dataset.elyPivotRows !== event.detail.measure; });
  });
  const drawing = root.querySelector('svg[data-ely-drawing-surface]');
  if (drawing) {
    let points = [];
    const toPoint = event => {
      const box = drawing.getBoundingClientRect();
      return [Math.round((event.clientX - box.left) * 320 / box.width), Math.round((event.clientY - box.top) * 120 / box.height)];
    };
    const finish = () => {
      if (points.length < 2) return;
      const line = document.createElementNS('http://www.w3.org/2000/svg', 'polyline');
      line.classList.add('ely-data__drawing-path');
      line.setAttribute('points', points.map(point => point.join(',')).join(' '));
      drawing.append(line);
      changed(root, { action: 'draw', paths: [...drawing.querySelectorAll('.ely-data__drawing-path')].map(path => path.getAttribute('points')) });
      points = [];
    };
    drawing.addEventListener('pointerdown', event => { if (root.querySelector('[data-ely-drawing][aria-pressed="true"]')?.value === 'erase') { drawing.querySelector('.ely-data__drawing-path:last-of-type')?.remove(); changed(root, { action: 'erase' }); return; } points = [toPoint(event)]; if(event.isTrusted)drawing.setPointerCapture(event.pointerId); });
    drawing.addEventListener('pointermove', event => { if (points.length) points.push(toPoint(event)); });
    drawing.addEventListener('pointerup', finish);
  }
});
