import assert from 'node:assert/strict';
import {mkdtemp,mkdir,writeFile,readFile,rm} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {join} from 'node:path';
import {execFileSync} from 'node:child_process';
const root=await mkdtemp(join(tmpdir(),'ely-assets-'));
try {
 await mkdir(join(root,'assets'));
 await mkdir(join(root,'playground'));
 await writeFile(join(root,'index.html'),'<link href="/ely-htmx/assets/ui.css"><script src="/ely-htmx/assets/main.js"></script>');
 await writeFile(join(root,'assets/ui.css'),'body{color:black}');
 await writeFile(join(root,'assets/main.js'),"import './child.js';");
 await writeFile(join(root,'assets/child.js'),'console.log("English")');
 await writeFile(join(root,'playground/index.html'),'<script src="./app.js"></script>');
 await writeFile(join(root,'playground/app.js'),'console.log("editor")');
 execFileSync(process.execPath,['scripts/version-pages-assets.mjs',root,'/ely-htmx/']);
 assert.match(await readFile(join(root,'index.html'),'utf8'),/ui\.css\?v=[a-f0-9]{16}/);
 assert.match(await readFile(join(root,'assets/main.js'),'utf8'),/child\.js\?v=[a-f0-9]{16}/);
 assert.match(await readFile(join(root,'playground/index.html'),'utf8'),/app\.js\?v=[a-f0-9]{16}/);
 const provider=await readFile('internal/playground/web/static-provider.js','utf8');
 assert(provider.includes('fetch(assetURL(name))')&&provider.includes("new Worker(assetURL('renderer-worker.js'))"));
 const worker=await readFile('internal/playground/web/renderer-worker.js','utf8');
 for(const dependency of ['./wasm_exec.js','./renderer.wasm','./styles.json']) assert(worker.includes(`assetURL('${dependency}')`));
 assert(worker.includes("fetch(assetURL('./'+path))"));
 console.log('Static HTML, module, worker and stylesheet cache version checks passed.');
}finally{await rm(root,{recursive:true,force:true});}
