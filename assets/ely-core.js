(() => {
  "use strict";
  const $ = (selector, root = document) => root.querySelector(selector);
  const $$ = (selector, root = document) => [...root.querySelectorAll(selector)];
  const closeMenus = except => $$('[data-core-menu]').forEach(trigger => {
    const panel = $("#" + trigger.dataset.coreMenu);
    if (panel && panel !== except) {
      panel.hidden = true;
      trigger.setAttribute('aria-expanded', 'false');
    }
  });
  const toggleMenu = trigger => {
    const panel = $("#" + trigger.dataset.coreMenu);
    if (!panel) return;
    const willOpen = panel.hidden;
    closeMenus(panel);
    panel.hidden = !willOpen;
    trigger.setAttribute('aria-expanded', String(willOpen));
  };
  const closeFocusTrap = trap => {
    trap.hidden = true;
    const entry = $(`[data-core-focus-trap-open="${trap.id}"]`);
    if (entry) {
      entry.setAttribute('aria-expanded', 'false');
      entry.focus();
    }
  };

  document.addEventListener('click', event => {
    const menu = event.target.closest('[data-core-menu]');
    if (menu) {
      event.preventDefault();
      toggleMenu(menu);
      return;
    }
    if (!event.target.closest('.core-menu')) closeMenus();

    const trapEntry = event.target.closest('[data-core-focus-trap-open]');
    if (trapEntry) {
      const trap = $("#" + trapEntry.dataset.coreFocusTrapOpen);
      if (trap) {
        trap.hidden = false;
        trapEntry.setAttribute('aria-expanded', 'true');
        $('input, button, [tabindex="0"]', trap)?.focus();
      }
      return;
    }

    const open = event.target.closest('[data-core-dialog-open]');
    if (open) {
      const dialog = $("#" + open.dataset.coreDialogOpen);
      if (dialog && typeof dialog.showModal === 'function') dialog.showModal();
      return;
    }
    if (event.target.closest('[data-core-dialog-close]')) {
      event.target.closest('dialog')?.close();
      return;
    }
    const action = event.target.closest('[data-core-action="press-status"]');
    if (action) {
      const output = $("#" + action.dataset.coreAction);
      if (output) output.textContent = 'Pressed locally.';
      return;
    }
    const toggle = event.target.closest('[data-core-toggle-button]');
    if (toggle) {
      toggle.setAttribute('aria-pressed', String(toggle.getAttribute('aria-pressed') !== 'true'));
      return;
    }
    const choice = event.target.closest('[data-core-choice] button');
    if (choice) {
      $$('button', choice.parentElement).forEach(button => button.setAttribute('aria-pressed', String(button === choice)));
      return;
    }
    const scroll = event.target.closest('[data-core-scrolltop]');
    if (scroll) {
      $("#" + scroll.dataset.coreScrolltop)?.scrollTo({top: 0, behavior: 'smooth'});
      return;
    }
    const copy = event.target.closest('[data-core-copy-label]');
    if (copy) {
      const text = $("#" + copy.dataset.coreCopyLabel);
      if (text) {
        const range = document.createRange();
        range.selectNodeContents(text);
        const selection = getSelection();
        selection?.removeAllRanges();
        selection?.addRange(range);
        copy.textContent = 'Text selected';
      }
      return;
    }
    const toast = event.target.closest('[data-core-toast]');
    if (toast) {
      const output = $('[data-core-toast-output]');
      if (output) {
        output.textContent = toast.dataset.coreToast;
        setTimeout(() => { output.textContent = ''; }, 2600);
      }
      return;
    }
    const coach = event.target.closest('[data-core-coach-close]');
    if (coach) {
      coach.closest('.core-coach')?.querySelector('aside')?.setAttribute('hidden', '');
      return;
    }
    const loading = event.target.closest('[data-core-loading]');
    if (loading) {
      const old = loading.textContent;
      loading.disabled = true;
      loading.textContent = 'Preparing…';
      setTimeout(() => { loading.disabled = false; loading.textContent = old; }, 700);
      return;
    }
    const zoom = event.target.closest('[data-core-zoom]');
    if (zoom) {
      const output = $('[data-core-zoom-output]');
      let value = Number.parseInt(output?.textContent) || 100;
      value = Math.max(70, Math.min(140, value + (zoom.dataset.coreZoom === 'in' ? 10 : -10)));
      if (output) output.textContent = value + '%';
    }
  });

  document.addEventListener('contextmenu', event => {
    const context = event.target.closest('[data-core-context]');
    if (!context) return;
    event.preventDefault();
    const trigger = $('[data-core-menu]', context);
    if (trigger) toggleMenu(trigger);
  });
  document.addEventListener('keydown', event => {
    const trap = event.target.closest?.('[data-core-focus-trap]');
    if (!trap) {
      if (event.key === 'Escape') closeMenus();
      return;
    }
    if (event.key === 'Escape') {
      event.preventDefault();
      closeFocusTrap(trap);
      return;
    }
    if (event.key !== 'Tab') return;
    const focusable = $$('a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex="0"]', trap);
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (!first || !last) return;
    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  });

  $$('[data-core-show]').forEach(input => input.addEventListener('change', () => {
    const target = $('#' + input.dataset.coreShow);
    if (target) target.hidden = !input.checked;
  }));
  $$('[data-core-keyboard]').forEach(box => box.addEventListener('keydown', event => {
    if (event.key === 'Enter' || event.key === 'Escape') {
      const output = $('output', box);
      if (output) output.textContent = event.key + ' received locally.';
    }
  }));
  $$('[data-core-tabs]').forEach(tabs => $$('[data-core-tab]', tabs).forEach(button => button.addEventListener('click', () => {
    $$('[data-core-tab]', tabs).forEach(item => item.setAttribute('aria-selected', String(item === button)));
    $$('[data-core-panel]', tabs).forEach(panel => { panel.hidden = panel.dataset.corePanel !== button.dataset.coreTab; });
  })));
  $$('[data-core-filter]').forEach(input => input.addEventListener('input', () => {
    $$('#' + input.dataset.coreFilter + ' button').forEach(row => { row.hidden = !row.textContent.toLowerCase().includes(input.value.toLowerCase()); });
  }));
  $$('[data-core-hold]').forEach(button => {
    let timer;
    const status = $('[data-core-hold-status]');
    const clear = () => { clearTimeout(timer); timer = undefined; button.classList.remove('is-holding'); };
    const start = () => {
      if (button.disabled || timer) return;
      button.classList.add('is-holding');
      timer = setTimeout(() => {
        if (status) status.textContent = 'Local review marker set for this specimen.';
        button.textContent = 'Reviewed';
        button.disabled = true;
        timer = undefined;
      }, 1000);
    };
    button.addEventListener('pointerdown', start);
    button.addEventListener('pointerup', clear);
    button.addEventListener('pointerleave', clear);
    button.addEventListener('pointercancel', clear);
    button.addEventListener('blur', clear);
    button.addEventListener('keydown', event => {
      if ((event.key === ' ' || event.key === 'Enter') && !event.repeat) {
        event.preventDefault();
        start();
      }
    });
    button.addEventListener('keyup', event => {
      if (event.key === ' ' || event.key === 'Enter') {
        event.preventDefault();
        clear();
      }
    });
  });
  $$('[data-core-typewriter]').forEach(node => {
    const text = node.dataset.coreTypewriter || '';
    let index = 0;
    const tick = () => { node.textContent = text.slice(0, ++index); if (index < text.length) setTimeout(tick, 35); };
    tick();
  });
  $$('[data-core-counter]').forEach(node => {
    const target = Number(node.dataset.coreCounter) || 0;
    let value = 0;
    const span = $('span', node);
    const tick = () => {
      value = Math.min(target, value + Math.ceil(target / 12));
      node.firstChild.nodeValue = value + ' ';
      if (value < target) setTimeout(tick, 35);
      else if (span) span.textContent = 'coded excerpts';
    };
    tick();
  });
  $$('[data-core-measure]').forEach(node => {
    const draw = () => { node.textContent = 'Current specimen width: ' + Math.round(node.getBoundingClientRect().width) + ' px'; };
    draw();
    addEventListener('resize', draw);
  });
})();
