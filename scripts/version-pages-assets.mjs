import {createHash} from 'node:crypto';
import {readdir, readFile, writeFile} from 'node:fs/promises';
import {resolve, dirname, relative, sep} from 'node:path';

// Each exported HTML page names the exact stylesheet/script content it needs.
// A cached stylesheet from an earlier deployment cannot style newer markup.
const root = resolve(process.argv[2] || 'dist');
const base = process.argv[3] || '/elysdui/';
async function files(dir) {
 const result=[];
 for (const entry of await readdir(dir,{withFileTypes:true})) {
  const path=resolve(dir,entry.name);
  if(entry.isDirectory()) result.push(...await files(path)); else result.push(path);
 }
 return result.sort();
}
const paths=await files(root);
// A release key also invalidates transitive module imports and lazy worker assets.
const release=createHash('sha256');
for(const path of paths) release.update(relative(root,path)).update(await readFile(path));
const version=release.digest('hex').slice(0,16);
const apiIdentities=new Set();
for(const path of paths.filter(path=>path.endsWith('.js'))) {
 let source=await readFile(path,'utf8');
 source=source.replace(/(\b(?:import|export)\s+(?:[^'"\n]*?\sfrom\s*)?['"])([^'"?#]+\.js)(['"])/g,(match,prefix,url,suffix)=>{
  // Absolute export-base imports and relative imports must resolve to the same
  // module identity, especially the shared component registry in ely-api.js.
  if(!url.startsWith('.')&&!url.startsWith(base)) return match;
  return `${prefix}${url}?v=${version}${suffix}`;
 });
 for(const [,specifier] of source.matchAll(/\b(?:import|export)\s+(?:[^'"\n]*?\sfrom\s*)?['"]([^'"]+)['"]/g)) {
  const url=new URL(specifier,new URL(base+relative(root,path).split(sep).join('/'),'https://static.invalid'));
  if(url.origin==='https://static.invalid'&&url.pathname.startsWith(base)&&url.pathname.endsWith('/ely-api.js')) apiIdentities.add(url.href);
 }
 await writeFile(path,`/* Static release ${version} */\n`+source);
}
if(apiIdentities.size>1) throw new Error(`Split component registry module identities: ${[...apiIdentities].join(', ')}`);
async function walk(dir) {
  for (const entry of await readdir(dir, {withFileTypes:true})) {
    const file = resolve(dir, entry.name);
    if (entry.isDirectory()) { await walk(file); continue; }
    if (!entry.name.endsWith('.html')) continue;
    const source = await readFile(file, 'utf8');
    const urls = [...source.matchAll(/(?:href|src)="([^"?#]+\.(?:css|js))"/g)];
    let output = source;
    for (const [, url] of urls) {
      if (/^(?:https?:)?\/\//.test(url)) continue;
      const asset = url.startsWith(base) ? resolve(root, url.slice(base.length)) : resolve(dirname(file), url);
      const path = relative(root, asset);
      if (path.startsWith('..'+sep) || path === '..') throw new Error(`Asset outside export: ${url}`);
      const hash = createHash('sha256').update(await readFile(asset)).digest('hex').slice(0,16);
      output = output.replaceAll(`"${url}"`, `"${url}?v=${hash}"`);
    }
    if (output !== source) await writeFile(file, output);
  }
}
await walk(root);
console.log('Versioned exported stylesheets, scripts and transitive assets.');
