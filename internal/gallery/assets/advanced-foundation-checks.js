import {registerChecks, assert} from './advanced-checks.js';
import {mount as mountAPI} from './ely-api.js';

const replayMotion = async (root, selector, verify) => {
  const original = root.querySelector(selector);
  assert(original, `missing ${selector}`);
  const originalValue = original.textContent;
  const clone = root.cloneNode(true);
  clone.id = `${root.id}-motion-audit`;
  clone.removeAttribute('data-foundation-mounted');
  clone.removeAttribute('data-ely-ready');
  const replay = clone.querySelector(selector);
  if (replay.matches('[data-foundation-number]')) replay.textContent = replay.dataset.foundationFrom || '0';
  if (replay.matches('[data-foundation-typewriter]')) replay.textContent = '';
  document.body.append(clone);
  let cleanup;
  try {
    cleanup = mountAPI(clone);
    await verify(replay);
    assert(original.textContent === originalValue, 'motion replay changed the original audit instance');
  } finally {
    cleanup?.();
    clone.remove();
  }
};

const actions = {
  show: root => { const input = root.querySelector('[data-foundation-show]'); assert(input, 'Show checkbox is missing'); if (input.disabled) return assert(input.disabled, 'disabled Show state was lost'); input.checked = !input.checked; input.dispatchEvent(new Event('change', {bubbles: true})); assert(root.querySelector('[data-foundation-conditional]').hidden === !input.checked, 'conditional content did not follow its checkbox'); },
  scope: root => { const trigger = root.querySelector('[data-foundation-open]'); assert(trigger, 'FocusScope trigger is missing'); if (trigger.disabled) return assert(trigger.disabled, 'disabled FocusScope state was lost'); trigger.click(); assert(!root.querySelector('[data-foundation-scope]').hidden, 'focus scope did not open'); },
  menu: root => { const trigger = root.querySelector('[data-foundation-menu]'); assert(trigger, 'ClickOutside trigger is missing'); if (trigger.disabled) return assert(trigger.disabled, 'disabled menu state was lost'); trigger.click(); assert(trigger.getAttribute('aria-expanded') === 'true', 'menu did not expose its expanded state'); },
  press: root => { const button = root.querySelector('[data-foundation-press]'); assert(button, 'Pressable button is missing'); if (button.disabled) return assert(button.disabled, 'disabled press state was lost'); button.click(); assert(root.querySelector('[data-foundation-output]').textContent === 'Pressed', 'press did not update its instance output'); },
  keys: root => { const target = root.querySelector('[data-foundation-keys]'); assert(target, 'KeyboardHandler target is missing'); target.dispatchEvent(new KeyboardEvent('keydown', {key: 'Enter', bubbles: true})); assert(target.querySelector('output').textContent === 'Enter', 'keyboard handler did not receive Enter'); },
  copy: root => { const button = root.querySelector('[data-foundation-copy-button]'); const value = root.querySelector('[data-foundation-copy]'); assert(button && value, 'copy control is missing'); if (button.disabled) return assert(button.disabled, 'disabled copy state was lost'); button.click(); assert(getSelection().toString() === value.textContent, 'copy action did not select caller text'); },
  split: root => { const range = root.querySelector('[data-foundation-split] input[type=range]'); assert(range, 'split range is missing'); if (range.disabled) return assert(range.disabled, 'disabled split state was lost'); range.value = '60'; range.dispatchEvent(new Event('change', {bubbles: true})); assert(range.closest('[data-foundation-split]').style.getPropertyValue('--foundation-split') === '60%', 'split position did not remain local'); },
  details: root => { const details = root.querySelector('details'); assert(details, 'disclosure is missing'); details.open = true; details.dispatchEvent(new Event('change', {bubbles: true})); assert(details.open, 'disclosure did not retain its native open state'); },
  dialog: root => { const trigger = root.querySelector('[data-foundation-dialog-open]'); const dialog = root.querySelector('dialog'); assert(trigger && dialog, 'dialog control is missing'); if (trigger.disabled) return assert(trigger.disabled, 'disabled dialog state was lost'); trigger.click(); assert(dialog.open, 'dialog did not open'); root.querySelector('[data-foundation-dialog-close]').click(); },
  master: root => { const item = root.querySelector('[data-foundation-master] nav button'); assert(item, 'master-detail item is missing'); if (item.disabled) return assert(item.disabled, 'disabled detail state was lost'); item.click(); assert(item.getAttribute('aria-current') === 'true', 'detail selection was not local'); },
  scroll: root => { const button = root.querySelector('[data-foundation-scroll-top]'); assert(button, 'scroll-to-top button is missing'); if (button.disabled) return assert(button.disabled, 'disabled scroll state was lost'); button.click(); assert(root.querySelector('[data-foundation-scroll]'), 'scroll container is missing'); },
  drop: root => { const zone = root.querySelector('[data-foundation-dropzone]'); assert(zone, 'dock zone is missing'); zone.dispatchEvent(new KeyboardEvent('keydown', {key: 'Enter', bubbles: true})); assert(zone.dataset.active === 'true', 'dock zone did not retain keyboard focus state'); },
  number: root => replayMotion(root, '[data-foundation-number]', async node => { const final = node.dataset.foundationNumber; if (node.getAttribute('aria-disabled') === 'true' || globalThis.matchMedia?.('(prefers-reduced-motion: reduce)').matches) return assert(node.textContent === final, 'reduced or disabled number did not finish immediately'); const initial = node.textContent; await new Promise(resolve => setTimeout(resolve, 60)); assert(node.textContent !== initial, 'number did not progress over time'); await new Promise(resolve => setTimeout(resolve, Number(node.dataset.foundationDuration) + 80)); assert(node.textContent === final, 'number did not reach the caller value'); }),
  typewriter: root => replayMotion(root, '[data-foundation-typewriter]', async node => { const final = node.dataset.foundationTypewriter; if (node.getAttribute('aria-disabled') === 'true' || globalThis.matchMedia?.('(prefers-reduced-motion: reduce)').matches) return assert(node.textContent === final, 'reduced or disabled typewriter did not finish immediately'); await new Promise(resolve => setTimeout(resolve, 60)); assert(final.startsWith(node.textContent) && node.textContent !== final, 'typewriter did not progress through caller text'); await new Promise(resolve => setTimeout(resolve, Number(node.dataset.foundationDuration) + 80)); assert(node.textContent === final, 'typewriter did not reach caller text'); })
};

const check = (selector, action) => async ({root}) => {
  const subject = root.querySelector(selector);
  assert(subject, `missing ${selector}`);
  assert(subject.outerHTML.toLowerCase().includes('caller'), `${selector} did not preserve caller data`);
  await actions[action]?.(root);
};

const primitives = {
  Box:['.ely-foundation__box'], Text:['span'], Icon:['.ely-foundation__icon'], IconSet:['.ely-foundation__icons button'], Svg:['svg'], Image:['figure img'], Canvas:['canvas[data-foundation-canvas]'], Divider:['.ely-foundation__divider'], Spacer:['.ely-foundation__spacer'], Portal:['.ely-foundation__portal aside'], Overlay:['.ely-foundation__overlay'], Slot:['.ely-foundation__surface footer'], Show:['[data-foundation-conditional]','show'], For:['ul li'], Fragment:['.ely-foundation__fragment'], VisuallyHidden:['.ely-foundation__sr'], FocusRing:['.ely-foundation__focus button'], FocusScope:['[data-foundation-open]','scope'], ClickOutside:['[data-foundation-menu]','menu'], HoverArea:['.ely-foundation__hover'], Pressable:['[data-foundation-press]','press'], KeyboardHandler:['[data-foundation-keys]','keys'], Measure:['[data-foundation-measure]'], IntersectionObserver:['[data-foundation-observe]'], Clipboard:['[data-foundation-copy-button]','copy'], 'Tooltip Trigger':['[role=tooltip]']
};
const typography = {
  Heading:['h1'], Title:['h2'], Paragraph:['.ely-foundation__prose'], Label:['label[for]'], Caption:['figcaption'], Overline:['.ely-foundation__overline'], Code:['code'], Kbd:['kbd'], KbdCombo:['kbd'], Blockquote:['blockquote'], Highlight:['mark'], Truncate:['.ely-foundation__truncate'], LineClamp:['.ely-foundation__clamp'], MiddleEllipsis:['.ely-foundation__middle'], EllipsisTooltip:['[title]'], Link:['a[href]'], ExternalLink:['a[target=_blank]'], RelativeTime:['time[data-foundation-format=relative]'], DateTimeText:['time[data-foundation-format=datetime]'], NumberText:['[data-foundation-format=number]'], CurrencyText:['[data-foundation-format=currency]'], PercentText:['[data-foundation-format=percent]'], FileSizeText:['[data-foundation-format=filesize]'], DurationText:['[data-foundation-format=duration]'], PluralText:['[data-foundation-format=plural]'], CopyableText:['[data-foundation-copy-button]','copy'], SelectableText:['.ely-foundation__selectable'], Emoji:['[aria-label]'], Latex:['.ely-foundation__math'], AnimatedNumber:['[data-foundation-number]','number'], Typewriter:['[data-foundation-typewriter]','typewriter'], GradientText:['.ely-foundation__gradient'], ShimmerText:['.ely-foundation__shimmer']
};
const layout = {
  Stack:['.ely-foundation__stack'], Flex:['.ely-foundation__flex'], Grid:['.ely-foundation__grid'], SimpleGrid:['.ely-foundation__grid'], Masonry:['.ely-foundation__masonry'], Center:['.ely-foundation__center'], Container:['.ely-foundation__container'], AspectRatio:['.ely-foundation__aspect'], Inset:['.ely-foundation__inset'], Wrap:['.ely-foundation__wrap'], Absolute:['.ely-foundation__absolute'], Layer:['.ely-foundation__layers'], ScrollArea:['.ely-foundation__scroll'], Scrollbar:['.ely-foundation__scrollbar'], ScrollShadow:['.ely-foundation__scroll-shadow'], ScrollToTop:['[data-foundation-scroll-top]','scroll'], StickyHeader:['.ely-foundation__sticky header'], SplitPane:['[data-foundation-split]','split'], ResizeHandle:['[data-foundation-split]','split'], Collapsible:['details','details'], Accordion:['[data-foundation-accordion] details','details'], Card:['.ely-foundation__card'], Panel:['.ely-foundation__panel'], Section:['.ely-foundation__section'], Group:['fieldset'], Frame:['.ely-foundation__frame'], Well:['.ely-foundation__well'], Sidebar:['.ely-foundation__sidebar'], Drawer:['[data-foundation-dialog-open]','dialog'], Sheet:['[data-foundation-dialog-open]','dialog'], Page:['.ely-foundation__page'], AppShell:['.ely-foundation__app-shell'], MasterDetail:['[data-foundation-master]','master'], Dock:['.ely-foundation__dock'], DockZone:['[data-foundation-dropzone]','drop'], FloatingPanel:['.ely-foundation__floating'], Pane:['.ely-foundation__panes'], Workspace:['.ely-foundation__workspace'], Viewport:['.ely-foundation__viewport']
};

registerChecks('foundation', Object.fromEntries([
  ...Object.entries(primitives).map(([name, [selector, action]]) => [`primitives/${name}`, check(selector, action)]),
  ...Object.entries(typography).map(([name, [selector, action]]) => [`typography/${name}`, check(selector, action)]),
  ...Object.entries(layout).map(([name, [selector, action]]) => [`layout/${name}`, check(selector, action)])
]));
