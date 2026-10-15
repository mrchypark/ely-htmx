/* All rendering and real gzip/Brotli/zstd work stays off the editor thread. */
'use strict';
(async () => {
  const start = performance.now();
  const version = new URL(self.location.href).search;
  const assetURL = path => path + version;
  importScripts(assetURL('./wasm_exec.js'));
  const go = new Go();
  const response = await fetch(assetURL('./renderer.wasm'));
  if (!response.ok) throw new Error(`Renderer download failed (${response.status})`);
  const bytes = await response.arrayBuffer();
  const result = await WebAssembly.instantiate(bytes,go.importObject);
  go.run(result.instance).catch(error => postMessage({fatal:error.message}));
  const manifest = await fetch(assetURL('./styles.json')).then(response => {if(!response.ok) throw new Error('Stylesheet manifest unavailable');return response.json();});
  const styles = {};
  await Promise.all(manifest.map(async path => {
    if (!/^assets\/[a-zA-Z0-9_./-]+\.css$/.test(path) || path.includes('..')) throw new Error('Invalid stylesheet path');
    const response = await fetch(assetURL('./'+path));
    if (!response.ok) throw new Error(`Stylesheet unavailable: ${path}`);
    styles[path] = await response.text();
  }));
  const initialized = globalThis.elyPlaygroundInit(JSON.stringify(styles));
  if (initialized) {const result = JSON.parse(initialized);if(result.error)throw new Error(result.error);}
  postMessage({metrics:{wasmBytes:bytes.byteLength,startupMs:performance.now()-start}});
  onmessage = ({data}) => {
    try {
      const result = JSON.parse(globalThis.elyPlaygroundRequest(data.method,data.path,data.body));
      const value = JSON.parse(result.body);
      if(result.status >= 400) throw new Error(value.error || 'Preview failed');
      postMessage({id:data.id,value});
    } catch(error) {postMessage({id:data.id,error:error.message});}
  };
  postMessage({ready:true});
})().catch(error => postMessage({fatal:error.message}));
