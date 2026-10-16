import assert from 'node:assert/strict';
import {spawnSync} from 'node:child_process';
const valid=['/','/elysdui/','/releases/v1.2/','/org.github.io/project/','/a_b/c-d/'];
const invalid=['','relative/','/missing-trailing','//','/a//b/','/./','/../','/a/./b/','/a/../b/','/a/?q=1','/a/#hash','/a b/','/a%2fb/'];
for(const [values,expected] of [[valid,0],[invalid,1]]) for(const base of values) {
 const result=spawnSync('bash',['-c','source scripts/pages-base.sh; valid_pages_base "$1"','validator',base]);
 assert.equal(result.status,expected,JSON.stringify(base));
}
console.log('Pages base validator accepts dotted segments and rejects traversal/malformed paths.');
