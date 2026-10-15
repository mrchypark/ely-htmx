/* Static transport: names/contracts are lazy files; only preview loads Go/WASM. */
(() => {
  'use strict';
  const base = new URL('./', document.currentScript.src);
  const version = new URL(document.currentScript.src).search;
  const assetURL = name => { const url = new URL(name, base); url.search = version; return url; };
  const pending = new Map();
  let catalogue, worker, sequence = 0, ready;
  async function json(name) {
    const response = await fetch(assetURL(name));
    if (!response.ok) throw new Error(`Static asset unavailable (${response.status}): ${name}`);
    return response.json();
  }
  function start() {
    if (ready) return ready;
    ready = new Promise((resolve, reject) => {
      if (!globalThis.Worker || !globalThis.WebAssembly) { reject(new Error('This browser needs Web Workers and WebAssembly. JSON editing/export still works.')); return; }
      worker = new Worker(assetURL('renderer-worker.js'));
      const timeout = setTimeout(() => fail(new Error('Go renderer loading timed out. Reload to retry; export JSON to keep your draft.')), 120000);
      function fail(error) { clearTimeout(timeout); reject(error); for (const item of pending.values()) item.reject(error); pending.clear(); worker.terminate(); }
      worker.onerror = event => fail(new Error(event.message || 'Go renderer could not load. Reload to retry.'));
      worker.onmessage = ({data}) => {
        if (data.ready) { clearTimeout(timeout); resolve(); return; }
        if (data.fatal) { fail(new Error(data.fatal)); return; }
        if (data.metrics) { globalThis.elyRendererMetrics = data.metrics; return; }
        const item = pending.get(data.id);
        if (!item) return;
        pending.delete(data.id);
        if (data.error) item.reject(new Error(data.error)); else item.resolve(data.value);
      };
    });
    return ready;
  }
  globalThis.elyStaticAPI = async (path, options = {}) => {
    const url = new URL(path, location.origin);
    options.signal?.throwIfAborted();
    if (url.pathname === '/api/examples') return json('examples.json');
    if (url.pathname === '/api/catalogue' || url.pathname === '/api/component') {
      const items = await (catalogue ||= json('catalogue.json'));
      options.signal?.throwIfAborted();
      if (url.pathname === '/api/component') {
        const item = items.find(item => item.name === url.searchParams.get('name'));
        if (!item) throw new Error('Unknown component');
        return json(item.file);
      }
      const q = (url.searchParams.get('q') || '').toLowerCase();
      if(q.length > 192) throw new Error('Search is too long');
      const matches = items.filter(item => `${item.name} ${item.family}`.toLowerCase().includes(q));
      matches.sort((a,b) => Number(b.name.toLowerCase() === q) - Number(a.name.toLowerCase() === q));
      return {items: matches.slice(0,40), total: items.length, matches: matches.length};
    }
    if (url.pathname !== '/api/preview') throw new Error('Unsupported static operation');
    const label = document.getElementById('preview-state');
    if (!ready && label) label.textContent = 'Loading Go renderer…';
    await start();
    options.signal?.throwIfAborted();
    return new Promise((resolve,reject) => {
      const id = ++sequence;
      pending.set(id,{resolve,reject});
      options.signal?.addEventListener('abort',() => {pending.delete(id);reject(new DOMException('Aborted','AbortError'));},{once:true});
      worker.postMessage({id, method:'POST',path:'/api/preview',body:options.body});
    });
  };
})();
