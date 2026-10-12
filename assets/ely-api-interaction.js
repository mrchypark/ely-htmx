import {register, emit} from './ely-api.js';

const closest = (root, target, selector) => target instanceof Element && target.closest(selector)?.closest('[data-ely-family]') === root ? target.closest(selector) : null;
const change = (root, detail = {}) => root.dispatchEvent(new CustomEvent('ely:change', {bubbles: true, detail: {id: root.id, component: root.dataset.elyComponent, ...detail}}));
const itemValue = item => item.dataset.interactionValue || item.value || item.textContent.trim();
const menu = (root, id, open) => {
  const panel = root.querySelector(`#${CSS.escape(id)}`);
  const trigger = root.querySelector(`[data-interaction-menu="${CSS.escape(id)}"]`);
  if (!panel || !trigger) return;
  panel.hidden = !open;
  trigger.setAttribute('aria-expanded', String(open));
  if (open) panel.querySelector('button:not(:disabled)')?.focus();
};

register('interaction', root => {
  const timers = new Set();
  const lottieCleanups = [];
  const later = (callback, delay) => { const timer = setTimeout(() => { timers.delete(timer); callback(); }, delay); timers.add(timer); return timer; };
  const close = () => root.querySelectorAll('[data-interaction-menu]').forEach(trigger => menu(root, trigger.dataset.interactionMenu, false));
  const click = event => {
    const button = closest(root, event.target, 'button, [data-interaction-draggable]');
    if (!button || button.disabled) return;
    const target = event.target instanceof Element ? event.target : null;
    const toggle = target?.closest('[data-interaction-toggle]');
    if (toggle) {
      const pressed = toggle.getAttribute('aria-pressed') !== 'true';
      toggle.setAttribute('aria-pressed', String(pressed));
      change(root, {name: toggle.name || '', value: toggle.value, pressed});
      return;
    }
    const choice = target?.closest('[data-interaction-choice]');
    if (choice) {
      choice.closest('[role="radiogroup"]')?.querySelectorAll('[data-interaction-choice]').forEach(item => item.setAttribute('aria-checked', String(item === choice)));
      change(root, {name: choice.name || '', value: choice.value, selected: true});
      return;
    }
    const select = target?.closest('[data-interaction-select]');
    if (select) {
      select.closest('.ely-interaction__items, [role="group"]')?.querySelectorAll('[data-interaction-select]').forEach(item => item.setAttribute('aria-current', item === select ? 'true' : 'false'));
      change(root, {name: select.name || '', value: select.value, selected: true});
      return;
    }
    const tab = target?.closest('[data-interaction-tab]');
    if (tab) {
      const tabs = [...(tab.closest('[data-interaction-tabs]')?.querySelectorAll('[data-interaction-tab]') || [])];
      tabs.forEach((item, index) => {
        const selected = item === tab;
        item.setAttribute('aria-selected', String(selected));
        item.tabIndex = selected ? 0 : -1;
        const panelID = item.getAttribute('aria-controls');
        const panel = panelID ? root.querySelector(`#${CSS.escape(panelID)}`) : null;
        if (panel) panel.hidden = !selected;
        if (selected) change(root, {name: 'tab', value: item.value || String(index), selected: true});
      });
      return;
    }
    const trigger = target?.closest('[data-interaction-menu]');
    if (trigger) { const panel = root.querySelector(`#${CSS.escape(trigger.dataset.interactionMenu)}`); menu(root, trigger.dataset.interactionMenu, Boolean(panel?.hidden)); return; }
    const opener = target?.closest('[data-interaction-dialog]');
    if (opener) { const dialog = root.querySelector(`#${CSS.escape(opener.dataset.interactionDialog)}`); if (dialog?.showModal) dialog.showModal(); else if (dialog) dialog.hidden = false; return; }
    if (target?.closest('[data-interaction-close]')) { const dialog = target.closest('dialog'); if (dialog?.open) dialog.close(); else target.closest('[role="dialog"], aside')?.setAttribute('hidden', ''); return; }
    const dismiss = target?.closest('[data-interaction-dismiss]');
    if (dismiss) { dismiss.closest('.ely-interaction')?.setAttribute('hidden', ''); change(root, {name: 'dismissed', value: 'true'}); return; }
    if (target?.closest('[data-interaction-more]')) { const output = root.querySelector('[data-interaction-more-count]'); if (output) { output.value = String(Number(output.value || output.textContent || 0) + 1); output.textContent = output.value; change(root, {name: 'count', value: output.value}); } return; }
    if (target?.closest('[data-interaction-loading]')) { button.setAttribute('aria-busy', String(button.getAttribute('aria-busy') !== 'true')); return; }
    if (target?.closest('[data-interaction-refresh]')) { button.setAttribute('aria-busy', 'true'); later(() => button.setAttribute('aria-busy', 'false'), 250); return; }
    if (target?.closest('[data-interaction-ripple]')) { button.classList.add('is-rippling'); later(() => button.classList.remove('is-rippling'), 300); return; }
    if (target?.closest('[data-interaction-resize]')) { const box = target.closest('[data-interaction-resizable]'); box?.classList.toggle('is-expanded'); change(root, {name: 'size', value: box?.classList.contains('is-expanded') ? 'expanded' : 'default'}); return; }
    if (target?.closest('[data-interaction-undo], [data-interaction-redo]')) { const history = root.querySelector('[data-interaction-history]'); const direction = target.hasAttribute('data-interaction-undo') ? 'undo' : 'redo'; if (history) history.textContent = direction; change(root, {name: 'history', value: direction}); return; }
  };
  const input = event => {
    const field = closest(root, event.target, '[data-interaction-filter]');
    if (!field || field.disabled) return;
    root.querySelectorAll('.ely-interaction__items li').forEach(item => { item.hidden = !item.textContent.toLocaleLowerCase().includes(field.value.toLocaleLowerCase()); });
    change(root, {name: field.name || 'filter', value: field.value});
  };
  const keydown = event => {
    const tab = closest(root, event.target, '[data-interaction-tab]');
    if (tab && ['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) {
      event.preventDefault();
      const tabs = [...tab.closest('[data-interaction-tabs]').querySelectorAll('[data-interaction-tab]')];
      const index = tabs.indexOf(tab);
      const next = event.key === 'Home' ? 0 : event.key === 'End' ? tabs.length - 1 : (index + (event.key === 'ArrowRight' ? 1 : tabs.length - 1)) % tabs.length;
      tabs[next].click(); tabs[next].focus(); return;
    }
    const hotkeys = closest(root, event.target, '[data-interaction-hotkeys]');
    if (hotkeys && hotkeys.getAttribute('aria-disabled') !== 'true') { const key = [event.ctrlKey && 'Ctrl', event.metaKey && 'Meta', event.altKey && 'Alt', event.shiftKey && 'Shift', event.key].filter(Boolean).join('+'); hotkeys.nextElementSibling.textContent = key; emit(root, hotkeys.dataset.interactionAction || 'hotkey', {key}); return; }
    const menuItem = closest(root, event.target, '[role="menu"] button');
    if (menuItem && ['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(event.key)) {
      event.preventDefault();
      const items = [...menuItem.closest('[role="menu"]').querySelectorAll('button:not(:disabled)')];
      const index = items.indexOf(menuItem);
      const next = event.key === 'Home' ? 0 : event.key === 'End' ? items.length - 1 : (index + (event.key === 'ArrowDown' ? 1 : items.length - 1)) % items.length;
      items[next]?.focus(); return;
    }
    const reorder = closest(root, event.target, '[data-interaction-reorder] [data-interaction-draggable]');
    if (reorder && ['ArrowUp', 'ArrowDown'].includes(event.key)) {
      event.preventDefault();
      const sibling = event.key === 'ArrowUp' ? reorder.previousElementSibling : reorder.nextElementSibling;
      if (!sibling) return;
      if (event.key === 'ArrowUp') reorder.parentElement.insertBefore(reorder, sibling); else reorder.parentElement.insertBefore(sibling, reorder);
      change(root, {name: 'reorder', value: itemValue(reorder)});
      emit(root, 'reorder', {value: itemValue(reorder), direction: event.key === 'ArrowUp' ? 'up' : 'down'});
      reorder.focus(); return;
    }
    const roving = closest(root, event.target, '[data-interaction-roving]');
    if (roving && event.key.length === 1) { const item = [...roving.querySelectorAll('button')].find(button => button.textContent.trim().toLocaleLowerCase().startsWith(event.key.toLocaleLowerCase())); item?.focus(); }
  };
  let dragging;
  const dragstart = event => { const item = closest(root, event.target, '[data-interaction-draggable]'); if (!item || item.getAttribute('aria-disabled') === 'true') return; dragging = item; event.dataTransfer?.setData('text/plain', itemValue(item)); emit(root, 'drag-start', {value: itemValue(item)}); };
  const dragover = event => { if (dragging && closest(root, event.target, '[data-interaction-drop], [data-interaction-reorder]')) event.preventDefault(); };
  const drop = event => { const dropzone = closest(root, event.target, '[data-interaction-drop], [data-interaction-reorder]'); if (!dragging || !dropzone) return; event.preventDefault(); if (dropzone.matches('[data-interaction-reorder]')) dropzone.append(dragging); change(root, {name: 'drag', value: itemValue(dragging)}); emit(root, 'drop', {value: itemValue(dragging)}); dragging = undefined; };
  const holdStart = event => { const button = closest(root, event.target, '[data-interaction-hold], [data-interaction-longpress]'); if (!button) return; button.dataset.interactionTimer = String(later(() => { root.querySelector('[data-interaction-status]')?.replaceChildren(document.createTextNode('confirmed')); emit(root, button.dataset.interactionLongpress || button.dataset.interactionHold || 'confirm', {value: button.value || ''}); }, 700)); };
  const holdEnd = event => { const button = closest(root, event.target, '[data-interaction-hold], [data-interaction-longpress]'); const timer = Number(button?.dataset.interactionTimer); if (timer) { clearTimeout(timer); timers.delete(timer); } };
  root.querySelectorAll('[data-interaction-lottie]').forEach(host => {
    const status = host.querySelector('[data-interaction-lottie-status]');
    const progress = host.querySelector('[data-interaction-lottie-progress]');
    const seek = host.querySelector('[data-interaction-lottie-seek]');
    const canvas = host.querySelector('[data-lottie-canvas]');
    const controller = new AbortController();
    let animation, destroyed = false;
    const fail = error => { if (destroyed) return; status.textContent = 'Animation unavailable'; root.dispatchEvent(new CustomEvent('ely:error', {bubbles: true, detail: {id: root.id, error}})); };
    const data = async () => {
      const url = host.dataset.lottieUrl;
      if (url) {
        const target = new URL(url, window.location.href);
        if (target.origin !== window.location.origin) throw new Error('Lottie animation URL must be same-origin');
        const response = await fetch(target, {signal: controller.signal, credentials: 'same-origin'});
        if (!response.ok) throw new Error(`Lottie animation request failed: ${response.status}`);
        return response.json();
      }
      const source = host.querySelector('[data-lottie-animation]')?.textContent.trim();
      if (!source) throw new Error('Lottie animation data is required');
      const parsed = JSON.parse(source);
      return typeof parsed === 'string' ? JSON.parse(parsed) : parsed;
    };
    const start = async () => {
      try {
        const library = window.lottie;
        if (!library?.loadAnimation) throw new Error('Lottie library is unavailable');
        const animationData = await data();
        if (destroyed) return;
        animation = library.loadAnimation({container: canvas, renderer: 'svg', loop: host.dataset.lottieLoop === 'true', autoplay: host.dataset.lottieAutoplay === 'true', animationData});
        animation.addEventListener('DOMLoaded', () => { if (!destroyed) status.textContent = 'Animation ready'; });
        animation.addEventListener('enterFrame', event => {
          if (destroyed || !animation.totalFrames) return;
          const value = Math.round(event.currentTime / animation.totalFrames * 100);
          progress.value = String(value); progress.textContent = `${value}%`; seek.value = String(value);
          change(root, {name: 'frame', value: String(value)});
        });
      } catch (error) { if (error.name !== 'AbortError') fail(error); }
    };
    const clickLottie = event => {
      if (!animation || !(event.target instanceof Element)) return;
      if (event.target.closest('[data-interaction-lottie-play]')) { animation.play(); status.textContent = 'Animation playing'; }
      if (event.target.closest('[data-interaction-lottie-pause]')) { animation.pause(); status.textContent = 'Animation paused'; }
    };
    const seekLottie = event => {
      if (!animation || event.target !== seek) return;
      animation.goToAndStop(Number(seek.value) / 100 * animation.totalFrames, true);
      progress.value = seek.value; progress.textContent = `${seek.value}%`;
      change(root, {name: 'frame', value: seek.value});
    };
    host.addEventListener('click', clickLottie); host.addEventListener('input', seekLottie); start();
    lottieCleanups.push(() => { destroyed = true; controller.abort(); host.removeEventListener('click', clickLottie); host.removeEventListener('input', seekLottie); animation?.destroy(); });
  });
  root.addEventListener('click', click); root.addEventListener('input', input); root.addEventListener('keydown', keydown);
  root.addEventListener('dragstart', dragstart); root.addEventListener('dragover', dragover); root.addEventListener('drop', drop);
  root.addEventListener('pointerdown', holdStart); root.addEventListener('pointerup', holdEnd); root.addEventListener('pointercancel', holdEnd);
  root.querySelector('[autofocus]')?.focus();
  return () => { timers.forEach(clearTimeout); lottieCleanups.forEach(cleanup => cleanup()); root.removeEventListener('click', click); root.removeEventListener('input', input); root.removeEventListener('keydown', keydown); root.removeEventListener('dragstart', dragstart); root.removeEventListener('dragover', dragover); root.removeEventListener('drop', drop); root.removeEventListener('pointerdown', holdStart); root.removeEventListener('pointerup', holdEnd); root.removeEventListener('pointercancel', holdEnd); };
});
