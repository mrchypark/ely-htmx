import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import vm from 'node:vm';

// Exercise the shipped state/mutation/parser code in Node, without browser or
// network access. Full DOM, rendering and interaction checks run in hosted CI.
class Element {
  constructor(tag = 'div') {
    this.tagName = tag.toUpperCase(); this.children = []; this.dataset = {};
    this.attributes = {}; this.listeners = {}; this.value = ''; this.textContent = '';
    this.hidden = false; this.disabled = false; this.content = 'test-csrf-token';
    const classes = new Set();
    this.classList = {toggle: (name, set) => set ? classes.add(name) : classes.delete(name)};
  }
  append(...nodes) { for (const node of nodes) { node.parentElement = this; this.children.push(node); } }
  replaceChildren(...nodes) { this.children = []; this.append(...nodes); }
  addEventListener(name, callback) { this.listeners[name] = callback; }
  setAttribute(name, value) { this.attributes[name] = value; }
  removeAttribute(name) { delete this.attributes[name]; }
  setCustomValidity(value) { this.validationMessage = value; }
  reportValidity() { return !this.validationMessage; }
  closest() { return this.parentElement || this; }
  querySelectorAll() { return []; }
  get firstChild() { return this.children[0] || {textContent: 'Test setting'}; }
}
const nodes = new Map();
const get = id => { if (!nodes.has(id)) nodes.set(id, new Element()); return nodes.get(id); };
get('compression-profile').value = 'dynamic';
get('compression-strategy').value = 'smallest';
get('min-bytes').value = '256';
get('min-savings-bytes').value = '32';
get('min-savings-percent').value = '5';
let storageFails = false;
const storage = new Map();
const pending = [];
let confirmations = 0;
const context = vm.createContext({
  console, TextEncoder, AbortController, setTimeout, clearTimeout, URL, Blob,
  window: {confirm: () => { confirmations++; return true; }},
  document: {getElementById: get, createElement: tag => new Element(tag), querySelector: () => new Element(), querySelectorAll: () => []},
  localStorage: {
    setItem: (key, value) => { if (storageFails) throw new Error('Quota exceeded'); storage.set(key, value); },
    getItem: key => storage.get(key),
  },
  fetch: async (path, options) => new Promise(resolve => pending.push({path, options, resolve})),
});
let source = await readFile(new URL('../internal/playground/web/app.js', import.meta.url), 'utf8');
assert(source.includes('  let searchTimer;'), 'test seam stays before bootstrap');
source = source.replace('  let searchTimer;', `
  globalThis.test = {
    state, parseSource, atPath, defaultValue, templateFor, validateValue,
    mutate, duplicateNode, deleteNode, moveNode, historyStep, propertyField,
    compressionSettings, persist, restore, renderPreview, schedulePreview, importFile, renderInspector,
    useStateRefresh() { refresh = () => { parseDraft(); }; },
  };
  return;
  let searchTimer;`);
vm.runInContext(source, context, {filename: 'playground-app.js'});
const t = context.test;
const json = value => JSON.parse(JSON.stringify(value));
const fixture = {
  version: '1', id: 'composition', loader: 'projectLoader',
  templates: {row: {params: {label: {type: 'string', required: true}}, root: {id: 'row-text', type: 'text', bindings: {text: '/@params/label'}}}},
  root: {id: 'root', type: 'stack', props: {gap: 'md'}, children: [
    {id: 'bound', type: 'text', bindings: {text: '/headline'}, when: {op: 'exists', path: '/headline'}},
    {id: 'repeated', use: 'row', paramBindings: {label: '/@vars/row/label'}, repeat: {items: '/rows', as: 'row', key: '/id'}},
    {id: 'footer', type: 'text', props: {text: 'Footer'}},
  ]},
};
function load(value = fixture) {
  t.state.screenText = JSON.stringify(value, null, 2);
  t.state.screen = JSON.parse(t.state.screenText);
  t.state.dataText = '{}'; t.state.data = {};
  t.state.screenError = ''; t.state.dataError = ''; t.state.visualSafe = true;
  t.state.selection = ['root', 'children', 2]; t.state.undo = []; t.state.redo = [];
}
t.useStateRefresh();
assert.throws(() => t.parseSource('{"root":{},"ROOT":{}}', 'Screen'), /duplicate key/);
assert.throws(() => t.parseSource('{"a":{"x":1,"x":2}}', 'Screen'), /duplicate key/);
assert.throws(() => t.parseSource('{"a":', 'Screen'));
assert.throws(() => t.parseSource('[]', 'Data'), /must be a JSON object/);
assert.equal(t.parseSource('{"id":9007199254740993}', 'Screen').precise, false);
assert.equal(t.parseSource('{"value":0.123456789012345678901}', 'Screen').precise, false);
assert.equal(t.parseSource('{"value":1e3,"zero":-0,"other":1.25}', 'Screen').precise, true);
assert.equal(t.atPath(['constructor'], {}), undefined, 'paths never traverse inherited keys');
assert.equal(t.atPath(['__proto__'], {}), undefined);
assert.equal(t.atPath(['__proto__', 'kept'], JSON.parse('{"__proto__":{"kept":true}}')), true, 'own special keys remain data');
assert.throws(() => t.validateValue({constructor: 'unexpected'}, {type: 'object', properties: {}}), /not a declared property/);
assert.throws(() => t.validateValue([1, 'two'], {type: 'array', items: {type: 'integer'}}), /integer/);
assert.throws(() => t.validateValue(7, {type: 'integer', max: 6}), /at most/);
assert.throws(() => t.validateValue({name: 7}, {type: 'object', properties: {name: {type: 'string', required: true}}}), /string/);
assert.deepEqual(json(t.templateFor({name: 'text', properties: {text: {type: 'string', required: true}}})), {id: 'new-component', type: 'text', props: {text: 'Your text here'}});
load();
t.moveNode(-1); assert.equal(t.state.screen.root.children[1].id, 'footer');
t.moveNode(1); assert.deepEqual(json(t.state.screen), fixture);
t.duplicateNode();
assert.equal(t.state.screen.root.children.length, 4);
assert.equal(t.state.screen.root.children[3].id, 'footer-copy-1');
t.deleteNode(); assert.deepEqual(json(t.state.screen), fixture, 'tree edits preserve complete composition metadata');
t.historyStep('undo'); assert.equal(t.state.screen.root.children.length, 4);
t.historyStep('redo'); assert.deepEqual(json(t.state.screen), fixture);
const exact = t.state.screenText;
t.state.screenError = 'Unfinished JSON'; t.state.screenText = '{"root":';
assert.equal(t.mutate(screen => { screen.root = {}; }), false);
assert.equal(t.state.screenText, '{"root":', 'invalid drafts cannot be overwritten by visual mutation');
t.state.screenText = exact; t.state.screenError = '';
load({version: '1', id: 'slots', root: {id: 'root', type: 'stack', children: [{id: 'card', type: 'playground/card', props: {title: 'Card'}, slots: {content: [
  {id: 'region', type: 'region', children: [{id: 'action', type: 'button', props: {label: 'Run'}, event: {on: 'click', action: 'save', target: 'region'}}]},
]}}]}});
t.state.selection = ['root', 'children', 0];
t.duplicateNode();
const copy = t.state.screen.root.children[1];
assert.equal(copy.slots.content[0].children[0].event.target, copy.slots.content[0].id, 'duplicate remaps internal event targets');
t.state.selection = ['root', 'children', 1, 'slots', 'content', 0];
t.deleteNode(); assert.deepEqual(json(t.state.selection), ['root', 'children', 1], 'last slot child selects owning node');
t.state.definitions.set('playground/card', Promise.resolve({name: 'playground/card', slots: {content: {maxNodes: 1}}}));
t.state.selection = ['root', 'children', 0, 'slots', 'content', 0];
const fullSlotSource = t.state.screenText;
await t.duplicateNode();
assert.equal(t.state.screenText, fullSlotSource, 'a full registered slot cannot be duplicated into');
load();
const bound = t.propertyField('text', {type: 'string'}, fixture.root.children[0], ['root', 'children', 0]);
assert(bound.children.some(child => child.className === 'binding-label'), 'bound property is displayed without replacing binding');
assert(!bound.children.some(child => child.tagName === 'INPUT'), 'bound properties do not become literal editors');
const field = t.propertyField('text', {type: 'string'}, t.state.screen.root.children[2], ['root', 'children', 2]);
const input = field.children.find(child => child.tagName === 'INPUT');
input.value = 'Edited'; input.listeners.change();
assert.equal(t.state.screen.root.children[2].props.text, 'Edited');
assert.deepEqual(json(t.state.screen.templates), fixture.templates);
const importOriginal = t.state.screenText;
await t.importFile({size: 10, text: async () => '{"root":'});
assert.equal(t.state.screenText, importOriginal, 'invalid imports never replace an existing draft');
assert.equal(confirmations, 0, 'invalid imports are rejected before a replacement confirmation');
await t.importFile({size: 30, text: async () => '{"root":{},"ROOT":{}}'});
assert.equal(t.state.screenText, importOriginal, 'duplicate-key imports retain the existing draft');
let finishOlderImport;
const oldImport = t.importFile({size: 100, text: () => new Promise(resolve => { finishOlderImport = resolve; })});
await t.importFile({size: 10, text: async () => '{"root":'});
finishOlderImport(JSON.stringify(fixture)); await oldImport;
assert.equal(t.state.screenText, importOriginal, 'a late older import cannot replace a newer import attempt');
t.state.definitions.set('text', Promise.resolve({name: 'text', properties: {text: {type: 'string'}}}));
await t.renderInspector();
const idInput = get('property-fields').children[0].children.find(child => child.tagName === 'INPUT');
idInput.value = 'bound'; idInput.listeners.change();
assert.match(idInput.validationMessage, /already used/);
assert.equal(t.state.screen.root.children[2].id, 'footer', 'ID collisions are rejected');
storageFails = true; t.persist();
assert.equal(t.state.saved, false); assert.match(get('save-status').textContent, /storage unavailable/);
storageFails = false; t.persist(); assert.equal(t.state.saved, true);

// Ignore cancellation at the fake transport layer to prove the independent
// version check also protects against a server response delivered after abort.
load();
t.state.previewVersion = 1;
const first = t.renderPreview(1);
assert.equal(pending.length, 1);
t.state.previewVersion = 2;
const second = t.renderPreview(2);
assert.equal(pending.length, 2);
const response = html => ({ok: true, status: 200, json: async () => ({html, sizes: {json: {identity: 10, gzip: 20, br: 15, zstd: 24, selected: 'identity'}, html: {identity: 30, gzip: 25, br: 20, zstd: 23, selected: 'br'}}, events: [], assets: []})});
pending[1].resolve(response('newest-preview')); await second;
pending[0].resolve(response('stale-preview')); await first;
assert.equal(get('preview-frame').srcdoc, 'newest-preview', 'late replies never replace a newer preview');
assert.equal(get('preview-frame').attributes.sandbox, '', 'preview grants no sandbox capabilities');
t.state.screenText = '{"version":"1","id":"large","root":{"id":"root","type":"text","props":{"text":"safe"}},"number":9007199254740993}';
t.state.previewVersion = 3;
const large = t.renderPreview(3);
assert(pending[2].options.body.includes('9007199254740993'), 'preview sends original numeric tokens');
pending[2].resolve({ok: false, status: 422, json: async () => ({error: 'Invalid contract'})}); await large;
assert.equal(get('preview-frame').srcdoc, 'newest-preview', 'server errors retain last successful preview');
assert.match(get('preview-state').textContent, /Stale preview/);
console.log('SDUI playground Node checks passed: safe parsing, exact drafts, composition mutations, history, typed fields, storage failure and latest-wins previews.');
