import {registerChecks, assert} from './advanced-checks.js';

const stories = {
  buttons: ['Button','Sizes','With icons','Icon button','ShortcutHint','ButtonGroup','SplitButton','ToggleButton','ToggleGroup','SegmentedControl','FloatingActionButton','ConfirmButton','CopyButton','ActionBar','BulkActionBar','ActionSheet','QuickActions'],
  navigation: ['Tabs','EditorTabs','Breadcrumb','Pagination','Steps','Sidebar Navigation','NavigationMenu','BackForwardNavigation','Anchor','JumpTo','CommandPalette','QuickOpen','QuickSwitcher','SearchPalette'],
  menus: ['Menu','CheckboxMenuItem','DropdownMenu','ContextMenu','MenuBar Menu','MegaMenu','SearchableMenu','PieMenu'],
  overlays: ['Popover','HoverCard','Dialog','Lightbox','Spotlight','FloatingToolbar','Callout','Tooltip'],
  feedback: ['Toast','Notification','Alert','Banner','Callout','InlineMessage','EmptyState','ErrorView','ErrorBoundary','Result','ConfirmationCard','SavingIndicator','Countdown'],
  motion: ['Spinner','Skeleton','Shimmer','LoadingButton','LoadMore','ProgressBar','LoadingOverlay','LazyLoad','Suspense','RefreshIndicator','UploadList','Transition','Stagger','AnimatePresence','Flip','Reorder','Marquee','Ripple','Glow','Blink','ParticleBackground'],
  interaction: ['Draggable','Resizable','LongPress','Hotkeys','Undo','VirtualScroller','Autofocus'],
};

const contracts = {
  buttons: {
    'Button':'[data-ely-action="save"]','Sizes':'[data-ely-action="press"]','With icons':'[data-ely-action="press"]','Icon button':'button[aria-label]','ShortcutHint':'.ely-interaction__shortcut','ButtonGroup':'[data-interaction-select]','SplitButton':'[data-interaction-menu]','ToggleButton':'[data-interaction-toggle]','ToggleGroup':'[data-interaction-toggle]','SegmentedControl':'[data-interaction-choice]','FloatingActionButton':'.ely-interaction__fab','ConfirmButton':'[data-interaction-hold]','CopyButton':'[data-ely-action="save"]','ActionBar':'[data-ely-action="action"]','BulkActionBar':'[data-interaction-bulk]','ActionSheet':'[data-interaction-dialog]','QuickActions':'.ely-interaction__quick [data-ely-action]'
  },
  navigation: {
    'Tabs':'[role="tab"]','EditorTabs':'[data-interaction-tab]','Breadcrumb':'.ui-breadcrumb','Pagination':'[data-interaction-page]','Steps':'.ely-interaction__steps','Sidebar Navigation':'.ely-interaction__nav a','NavigationMenu':'nav a','BackForwardNavigation':'[data-ely-action="back"]','Anchor':'.ely-interaction__nav a','JumpTo':'form[data-ely-submit="jump"]','CommandPalette':'[data-interaction-dialog]','QuickOpen':'form[data-ely-submit="open"]','QuickSwitcher':'[data-interaction-roving]','SearchPalette':'form[data-ely-submit="search"]'
  },
  menus: {
    'Menu':'[role="menu"]','CheckboxMenuItem':'input[type="checkbox"]','DropdownMenu':'[data-interaction-menu]','ContextMenu':'[data-interaction-context]','MenuBar Menu':'nav [data-interaction-menu]','MegaMenu':'.ely-interaction__mega section','SearchableMenu':'[data-interaction-filter]','PieMenu':'.ely-interaction__pie [role="menuitem"]'
  },
  overlays: {
    'Popover':'[role="dialog"]','HoverCard':'[role="tooltip"]','Dialog':'dialog','Lightbox':'dialog figure','Spotlight':'.ely-interaction__spotlight','FloatingToolbar':'[role="toolbar"]','Callout':'.ely-interaction__callout','Tooltip':'[role="tooltip"]'
  },
  feedback: {
    'Toast':'.ely-interaction__toast','Notification':'.ely-interaction__notification','Alert':'aside.ui-alert','Banner':'.ely-interaction__banner','Callout':'.ely-interaction__callout','InlineMessage':'.ely-interaction__message','EmptyState':'.ui-empty','ErrorView':'.ely-interaction__error','ErrorBoundary':'.ely-interaction__error','Result':'.ely-interaction__result','ConfirmationCard':'article.ui-card','SavingIndicator':'.ely-interaction__saving','Countdown':'[data-interaction-countdown]'
  },
  motion: {
    'Spinner':'.ely-interaction__spinner','Skeleton':'.ely-interaction__skeleton','Shimmer':'.ely-interaction__shimmer','LoadingButton':'[data-interaction-loading]','LoadMore':'[data-interaction-more]','ProgressBar':'progress','LoadingOverlay':'.ely-interaction__loading','LazyLoad':'[data-interaction-lazy]','Suspense':'[aria-busy]','RefreshIndicator':'[data-interaction-refresh]','UploadList':'.ely-interaction__items','Transition':'[data-interaction-transition]','Stagger':'[data-interaction-stagger]','AnimatePresence':'[data-interaction-presence]','Flip':'[data-interaction-flip]','Reorder':'[data-interaction-reorder]','Marquee':'.ely-interaction__marquee','Ripple':'[data-interaction-ripple]','Glow':'.ely-interaction__glow','Blink':'.ely-interaction__blink','ParticleBackground':'.ely-interaction__particles, [data-interaction-lottie]'
  },
  interaction: {
    'Draggable':'[data-interaction-draggable]','Resizable':'[data-interaction-resizable]','LongPress':'[data-interaction-longpress]','Hotkeys':'[data-interaction-hotkeys]','Undo':'[data-interaction-undo]','VirtualScroller':'[data-interaction-scroller]','Autofocus':'[data-interaction-roving]'
  },
};

const requires = {
  filter: new Set(['navigation/CommandPalette','menus/SearchableMenu']),
  toggle: new Set(['buttons/ToggleButton','buttons/ToggleGroup','overlays/FloatingToolbar']),
  choice: new Set(['buttons/SegmentedControl']),
  tab: new Set(['navigation/Tabs','navigation/EditorTabs']),
  menu: new Set(['buttons/SplitButton','menus/Menu','menus/DropdownMenu','menus/ContextMenu','menus/MenuBar Menu','overlays/Popover']),
  dialog: new Set(['buttons/ActionSheet','navigation/CommandPalette','overlays/Dialog','overlays/Lightbox']),
  resize: new Set(['interaction/Resizable']), undo: new Set(['interaction/Undo']), hotkeys: new Set(['interaction/Hotkeys']),
  ripple: new Set(['motion/Ripple']), reorder: new Set(['motion/Reorder']), hold: new Set(['buttons/ConfirmButton','interaction/LongPress']),
  form: new Set(['navigation/JumpTo','navigation/QuickOpen','navigation/SearchPalette','overlays/Dialog']),
};

const check = async ({root}, selector) => {
  const key = root.dataset.elyComponent;
  const need = (name, node) => { if (requires[name]?.has(key)) assert(node, `${key} is missing its ${name} control`); return node; };
  const emitted = [];
  root.addEventListener('ely:action', event => emitted.push(event.detail));
  assert(root.dataset.elyReady === 'true', 'interaction root did not mount');
  if (root.querySelector('[data-interaction-lottie]')) return lottieCheck({root});
  const contract = root.querySelector(selector);
  assert(contract, `${key} is missing interaction contract ${selector}`);
  assert(contract.getAttribute('data-interaction-unsupported') !== 'lottie', `${key} used an unsupported Lottie placeholder`);
  const filter = need('filter', root.querySelector('[data-interaction-filter]'));
  if (filter && !filter.disabled) { filter.value = 'second'; filter.dispatchEvent(new Event('input', {bubbles: true})); const first = root.querySelector('.ely-interaction__items li'); assert(first?.hidden === true, 'filter did not hide a nonmatching item'); }
  const toggle = need('toggle', root.querySelector('[data-interaction-toggle]'));
  if (toggle && !toggle.disabled) { const was = toggle.getAttribute('aria-pressed'); toggle.click(); assert(toggle.getAttribute('aria-pressed') !== was, 'toggle state did not change'); }
  const choice = need('choice', root.querySelector('[data-interaction-choice]'));
  if (choice && !choice.disabled) { choice.click(); assert(choice.getAttribute('aria-checked') === 'true', 'segmented choice did not select'); }
  const tab = need('tab', root.querySelector('[data-interaction-tab]'));
  if (tab && !tab.disabled) { tab.click(); assert(tab.getAttribute('aria-selected') === 'true', 'tab did not select'); tab.dispatchEvent(new KeyboardEvent('keydown', {key: 'ArrowRight', bubbles: true})); }
  const menu = need('menu', root.querySelector('[data-interaction-menu]'));
  if (menu && !menu.disabled) { menu.click(); assert(menu.getAttribute('aria-expanded') === 'true', 'menu did not open'); const item = root.querySelector('[role=menu] button:not(:disabled)'); if (item) { item.focus(); item.dispatchEvent(new KeyboardEvent('keydown', {key: 'ArrowDown', bubbles: true})); assert(root.querySelector('[role=menu]').contains(document.activeElement), 'menu arrow navigation lost menu focus'); } }
  const dialog = need('dialog', root.querySelector('[data-interaction-dialog]'));
  if (dialog && !dialog.disabled) { dialog.click(); const target = root.querySelector(`#${CSS.escape(dialog.dataset.interactionDialog)}`); assert(target?.open || target?.hidden === false, 'dialog did not open'); target?.querySelector('[data-interaction-close]')?.click(); }
  const resize = need('resize', root.querySelector('[data-interaction-resize]'));
  if (resize && !resize.disabled) { resize.click(); const panel=resize.closest('[data-interaction-resizable]'); assert(panel.classList.contains('is-expanded'), 'resize did not change local size'); assert(panel.getBoundingClientRect().width<=panel.parentElement.getBoundingClientRect().width,'expanded panel exceeds its container'); }
  const undo = need('undo', root.querySelector('[data-interaction-undo]'));
  if (undo && !undo.disabled) { undo.click(); assert(root.querySelector('[data-interaction-history]').textContent === 'undo', 'undo did not update history'); }
  const hotkeys = need('hotkeys', root.querySelector('[data-interaction-hotkeys]'));
  if (hotkeys && hotkeys.getAttribute('aria-disabled') !== 'true') { hotkeys.dispatchEvent(new KeyboardEvent('keydown', {key: 'k', ctrlKey: true, bubbles: true})); const event = emitted.at(-1); assert(event?.action === hotkeys.dataset.interactionAction && event.key === 'Ctrl+k', 'hotkey payload was not exact'); }
  const ripple = need('ripple', root.querySelector('[data-interaction-ripple]'));
  if (ripple && !ripple.disabled) { ripple.click(); assert(ripple.classList.contains('is-rippling'), 'ripple did not begin'); }
  const reorder = need('reorder', root.querySelector('[data-interaction-reorder] [data-interaction-draggable]'));
  if (reorder) { reorder.dispatchEvent(new KeyboardEvent('keydown', {key: 'ArrowDown', bubbles: true})); const event = emitted.at(-1); assert(event?.action === 'reorder' && event.value === reorder.dataset.interactionValue && event.direction === 'down', 'keyboard reorder payload was not exact'); }
  const action = root.querySelector('[data-ely-action]:not(:disabled)');
  if (action) { action.click(); const event = emitted.find(item => item.action === action.dataset.elyAction); assert(event?.id === root.id && event.component === root.dataset.elyComponent && event.value === (action.value || action.dataset.value || ''), 'button action payload was not exact'); }
  const form = need('form', root.querySelector('form[data-ely-submit]'));
  if (root.dataset.elyDisabled !== 'true' && form?.checkValidity()) { form.dispatchEvent(new Event('submit', {bubbles: true, cancelable: true})); const event = emitted.find(item => item.action === form.dataset.elySubmit); assert(event?.id === root.id && event.component === root.dataset.elyComponent && Array.isArray(event.values), 'form action payload was not exact'); }
  const hold = need('hold', root.querySelector('[data-interaction-hold], [data-interaction-longpress]'));
  if (hold && !hold.disabled) { hold.dispatchEvent(new PointerEvent('pointerdown', {bubbles: true})); await new Promise(resolve => setTimeout(resolve, 720)); hold.dispatchEvent(new PointerEvent('pointerup', {bubbles: true})); const event = emitted.at(-1); const expected = hold.dataset.interactionLongpress || hold.dataset.interactionHold; assert(event?.action === expected && event.value === hold.value, 'hold action payload was not exact'); }
  const disabled = root.querySelector('button:disabled');
  if (disabled) { const before = emitted.length; disabled.click(); assert(emitted.length === before, 'disabled button emitted an action'); }
};

const lottieCheck = async ({root}) => {
  const host = root.querySelector('[data-interaction-lottie]');
  const play = root.querySelector('[data-interaction-lottie-play]');
  const pause = root.querySelector('[data-interaction-lottie-pause]');
  const seek = root.querySelector('[data-interaction-lottie-seek]');
  const progress = root.querySelector('[data-interaction-lottie-progress]');
  const status = root.querySelector('[data-interaction-lottie-status]');
  assert(host && play && pause && seek && progress && status, 'Lottie Player is missing playback controls');
  for (let attempt = 0; attempt < 10 && !host.querySelector('svg'); attempt++) await new Promise(resolve => setTimeout(resolve, 50));
  assert(host.querySelector('svg'), 'Lottie Player did not render SVG frames');
  seek.value = '50'; seek.dispatchEvent(new Event('input', {bubbles: true}));
  assert(progress.textContent === '50%', 'Lottie seek did not update frame progress');
  pause.click(); assert(status.textContent === 'Animation paused', 'Lottie pause did not change playback state');
  play.click(); assert(status.textContent === 'Animation playing', 'Lottie play did not change playback state');
};

const checks = Object.fromEntries(Object.entries(stories).flatMap(([chapter, names]) => names.map(name => [`${chapter}/${name}`, context => check(context, contracts[chapter][name])])));
checks['motion/Lottie Player'] = lottieCheck;
registerChecks('interaction', checks);
