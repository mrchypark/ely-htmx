import assert from 'node:assert/strict';
import {spawn} from 'node:child_process';
import {createServer} from 'node:net';
import {fileURLToPath} from 'node:url';
import {setTimeout as delay} from 'node:timers/promises';
import {chromium} from 'playwright';

// A real Go example, a real Chromium browser, and the actual bundled assets.
// Every request stays on this loopback server; no external service is involved.
const repo = fileURLToPath(new URL('../', import.meta.url));
const reservation = createServer();
await new Promise((resolve, reject) => {
  reservation.once('error', reject);
  reservation.listen(0, '127.0.0.1', resolve);
});
const port = reservation.address().port;
await new Promise(resolve => reservation.close(resolve));
const origin = `http://127.0.0.1:${port}`;
const child = spawn('go', ['run', './examples/sdui'], {
  cwd: repo, env: {...process.env, PORT: String(port), SDUI_SCHEMA: ''},
  detached: process.platform !== 'win32', stdio: ['ignore', 'pipe', 'pipe'],
});
let output = '', startupError, browser;
child.stdout.on('data', value => { output += value; });
child.stderr.on('data', value => { output += value; });
child.once('error', error => { startupError = error; });

async function eventually(check, message, timeout = 15_000) {
  const start = Date.now();
  while (Date.now() - start < timeout) {
    if (await check()) return;
    await delay(25);
  }
  throw new Error(message);
}

try {
  await eventually(async () => {
    if (startupError) throw startupError;
    if (child.exitCode !== null) throw new Error(`Go example exited (${child.exitCode}): ${output}`);
    try { return (await fetch(origin)).ok; } catch { return false; }
  }, `Go example did not start: ${output}`, 60_000);
  browser = await chromium.launch();
  const context = await browser.newContext({viewport: {width: 1100, height: 1000}, serviceWorkers: 'block'});
  const page = await context.newPage();
  const errors = [], requests = [], postRequests = [], assetResponses = [];
  page.on('pageerror', error => errors.push(String(error)));
  page.on('response', response => {
    if (new URL(response.url()).pathname.startsWith('/assets/')) assetResponses.push({url: response.url(), status: response.status()});
  });
  page.on('request', request => {
    assert.equal(new URL(request.url()).origin, origin, 'browser never calls an external origin');
    requests.push(new URL(request.url()).pathname);
    if (request.method() === 'POST') postRequests.push(request);
  });
  await page.addInitScript(() => {
    window.sduiResults = [];
    window.elyErrors = [];
    window.sduiSwaps = 0;
    for (const type of ['sdui:success', 'sdui:error']) {
      document.addEventListener(type, event => sduiResults.push({type, status: event.detail.status, validation: event.detail.validation}));
    }
    document.addEventListener('ely:error', event => elyErrors.push(String(event.detail.error)));
    document.addEventListener('htmx:after:swap', () => sduiSwaps++);
  });
  await page.goto(origin, {waitUntil: 'networkidle'});
  const name = page.locator('#profile-editor input[name="name"]');
  const toggle = page.locator('#updates-switch input[role="switch"]');
  const submit = page.locator('#save-profile button');
  const result = page.locator('#profile-result');
  const ready = () => page.waitForSelector('#updates-switch [data-ely-ready="true"]');
  await ready();
  assert.equal(await page.locator('#profile-screen').count(), 1);
  assert.equal(requests.filter(path => path === '/assets/sdui.js').length, 1);
  assert.equal(requests.filter(path => path === '/assets/ely-api-forms-entry.js').length, 1);
  assert.equal(requests.filter(path => path === '/assets/ely-api.js').length, 1);
  for (const family of ['foundation', 'interaction', 'data', 'workspace', 'communication', 'native']) {
    assert(!requests.includes(`/assets/ely-api-${family}-entry.js`), `unused ${family} entrypoint is absent`);
  }
  assert(!requests.includes('/assets/ely-api-all.js'), 'no all-family loader');

  // Server validation replaces the requested region even on HTTP 422, and
  // keyboard focus follows the input's stable generated ID across the swap.
  await name.focus();
  await name.press('Enter');
  await page.waitForFunction(() => sduiResults.some(result => result.status === 422));
  assert.match(await page.locator('#profile-error').innerText(), /Enter a display name\./);
  await ready();
  assert.equal(await name.evaluate(node => document.activeElement === node), true, 'focus restored after validation');
  assert.equal(await page.locator('#profile-form').count(), 1, 'no duplicate region IDs after 422');

  // Successful native FormData includes repeated names, the clicked submitter,
  // and a real File. File bytes are read by Go, not JSON-serialized or mocked.
  await name.fill('Ada Lovelace');
  const contents = Buffer.from('Real browser attachment.\n');
  await page.locator('input[name="attachment"]').setInputFiles({name: 'proof.txt', mimeType: 'text/plain', buffer: contents});
  await submit.click();
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Save #1.'));
  assert.match(await result.innerText(), /Saved Ada Lovelace\. Updates off\./);
  assert.match(await page.locator('#attachment-result').innerText(), new RegExp(`proof\\.txt \\(${contents.length} bytes\\)`));
  assert.match(await page.locator('#tags-result').innerText(), /Repeated fields: go, hypermedia/);
  assert.match(await page.locator('#intent-result').innerText(), /Submitter: save/);
  await ready();
  assert.equal(await submit.evaluate(node => document.activeElement === node), true, 'submit button focus restored');

  // Advanced component change events use the same transport and current form.
  let before = postRequests.length;
  await toggle.check();
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Save #2.'));
  await ready();
  assert.equal(postRequests.length, before + 1, 'one advanced change causes one request');
  assert.equal(await toggle.isChecked(), true);
  assert.match(await result.innerText(), /Updates on\./);
  assert.equal(await toggle.evaluate(node => document.activeElement === node), true, 'advanced control focus restored');

  // Replacing controls repeatedly must clean old family handlers and retain a
  // single delegated bridge. Retain an old root to explicitly test detachment.
  await page.evaluate(() => {
    window.detachedSwitch = document.querySelector('#updates-switch [data-ely-family]');
    window.detachedChanges = 0;
    detachedSwitch.addEventListener('ely:change', () => detachedChanges++);
  });
  before = postRequests.length;
  await toggle.uncheck();
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Save #3.'));
  await ready();
  assert.equal(postRequests.length, before + 1, 'replacement did not duplicate handlers');
  await page.evaluate(() => detachedSwitch.querySelector('input').dispatchEvent(new Event('change', {bubbles: true})));
  assert.equal(await page.evaluate(() => detachedChanges), 1, 'only the pre-removal change was observed on detached root');
  assert.equal(postRequests.length, before + 1, 'detached controls never send requests');

  // Repeated in-flight submit events on the same form are suppressed. The
  // server receives one request; later intentional submissions remain allowed.
  before = postRequests.length;
  await name.fill('Repeated submission');
  await page.locator('input[name="delay"]').evaluate(node => { node.value = '150'; });
  await page.evaluate(() => {
    const form = document.querySelector('#profile-editor');
    const button = document.querySelector('#save-profile button');
    form.requestSubmit(button);
    form.requestSubmit(button);
  });
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Save #4.'));
  await ready();
  assert.equal(postRequests.length, before + 1, 'busy duplicate submit does not reach server');

  // Force actual responses to arrive out of order. Disable AbortController's
  // cancellation just for this probe: a sequence guard must work even when a
  // transport cannot abort a request that already committed on the server.
  await page.evaluate(() => {
    window.realAbort = AbortController.prototype.abort;
    AbortController.prototype.abort = function () {};
  });
  before = postRequests.length;
  let olderResponseFinished = false;
  const observeOlder = response => {
    if (response.request().method() === 'POST' && response.request().postData()?.includes('Older response')) {
      void response.finished().then(() => { olderResponseFinished = true; });
    }
  };
  page.on('response', observeOlder);
  await name.fill('Older response');
  await page.locator('input[name="delay"]').evaluate(node => { node.value = '300'; });
  await page.evaluate(() => document.querySelector('#profile-editor').requestSubmit(document.querySelector('#save-profile button')));
  // Poll the authenticated loader, rather than relying on a sleep, to ensure
  // the older mutation committed before issuing the newer action.
  await eventually(async () => (await (await context.request.get(origin)).text()).includes('Saved Older response.'), 'older request did not commit');
  assert.equal(olderResponseFinished, false, 'race probe starts the newer action before the older response arrives');
  // One browser task avoids scheduling three automation round trips inside the
  // deliberately short server delay. A native click still runs Ely's handlers.
  await page.evaluate(() => {
    document.querySelector('input[name="name"]').value = 'Latest response';
    document.querySelector('input[name="delay"]').value = '0';
    document.querySelector('#updates-switch input[role="switch"]').click();
  });
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Saved Latest response.'));
  await ready();
  assert.match(await result.innerText(), /Save #6\./);
  await eventually(() => olderResponseFinished, 'older response did not finish');
  await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => resolve())));
  page.off('response', observeOlder);
  assert.match(await result.innerText(), /Saved Latest response\./, 'late stale response cannot replace the newest DOM');
  assert.equal(await name.inputValue(), 'Latest response');
  assert.equal(postRequests.length, before + 2);
  await page.evaluate(() => { AbortController.prototype.abort = realAbort; });
  const refreshed = await context.request.get(origin);
  assert.match(await refreshed.text(), /Saved Latest response\./, 'newest server state survives refresh');

  // HTTP negatives use the actual running server and leave the rendered region
  // untouched. The example never creates a session on unauthenticated POST.
  const csrf = await page.locator('[data-sdui-screen]').getAttribute('data-sdui-csrf');
  const post = (values, key) => context.request.post(origin, {form: values, headers: {'Idempotency-Key': key}});
  const payload = {_action: 'save-profile', _target: 'profile-form', _csrf: csrf, _request: 'browser-retry', name: 'Idempotent request'};
  assert.equal((await post({...payload, _csrf: 'wrong'}, 'bad-csrf')).status(), 403);
  assert.equal((await post({...payload, _action: 'not-registered'}, 'bad-action')).status(), 400);
  const anonymous = await fetch(origin, {method: 'POST', headers: {'Content-Type': 'application/x-www-form-urlencoded', 'Idempotency-Key': 'anonymous'}, body: new URLSearchParams(payload)});
  assert.equal(anonymous.status, 401);
  const first = await post(payload, 'browser-retry');
  const second = await post(payload, 'browser-retry');
  assert.equal(first.status(), 200);
  assert.equal(second.status(), 200);
  assert.equal(await first.text(), await second.text(), 'actual HTTP retry replays byte-identical HTML');
  assert.equal((await post({...payload, name: 'different'}, 'browser-retry')).status(), 409);

  const results = await page.evaluate(() => sduiResults);
  assert.equal(results.filter(item => item.status === 422 && item.validation).length, 1);
  assert.equal(results.filter(item => item.type === 'sdui:success').length, 5, 'superseded response emits no false success');
  assert.equal(await page.evaluate(() => sduiSwaps), 6, 'one swap per accepted response');
  assert.deepEqual(await page.evaluate(() => elyErrors), []);
  assert.deepEqual(errors, [], 'no uncaught browser errors');
  assert(assetResponses.length > 0);
  assert.deepEqual(assetResponses.filter(response => response.status !== 200), [], 'every initial asset request succeeds');

  // Normal document history remains usable. A rejected loader URL must not
  // break delegation after Back/Forward, including browsers using bfcache.
  const invalidPage = await page.goto(`${origin}/?unsupported=1`);
  assert.equal(invalidPage.status(), 400);
  await page.goBack();
  await ready();
  assert.equal(await page.locator('#profile-form').count(), 1);
  await page.goForward();
  assert.match(await page.locator('body').innerText(), /invalid page request/);
  await page.goBack();
  await ready();
  before = postRequests.length;
  await name.fill('After history navigation');
  await submit.click();
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Saved After history navigation.'));
  assert.match(await result.innerText(), /Save #8\./, 'server persisted the prior idempotent save exactly once');
  assert.equal(postRequests.length, before + 1, 'history restoration has one active transport handler');
  assert.deepEqual(errors, [], 'no browser errors after Back/Forward');
  // Native forms expose controls as named properties. Real Chromium, unlike
  // the isolated DOM fixture, can shadow these methods with innocent field
  // names. The transport must call native methods and recognize the form.
  before = postRequests.length;
  await page.evaluate(() => {
    const form = document.getElementById('profile-editor');
    for (const name of ['closest', 'getAttribute', 'setAttribute', 'removeAttribute', 'matches', 'querySelector', 'form']) {
      const field = document.createElement('input');
      field.type = 'hidden'; field.name = name; field.value = 'ordinary field value';
      form.appendChild(field);
    }
  });
  await name.fill('Named fields remain safe');
  await submit.click();
  await page.waitForFunction(() => document.querySelector('#profile-result')?.textContent.includes('Saved Named fields remain safe.'));
  assert.match(await result.innerText(), /Save #9\./, 'named properties do not interfere with native form submission');
  assert.equal(postRequests.length, before + 1, 'clobbering regression submits exactly once');
  assert.deepEqual(errors, [], 'named fields cause no browser errors');
  if (process.env.SDUI_SCREENSHOT) await page.screenshot({path: process.env.SDUI_SCREENSHOT, fullPage: true});
  console.log('SDUI browser checks passed: real Go server, selective bundled assets, validation, FormData/files/repeated fields/submitter, advanced changes, repeated swaps and cleanup, focus, duplicate suppression, late-response race, Back/Forward, named-field clobbering, auth/CSRF and idempotent replay.');
} finally {
  await browser?.close();
  if (child.pid && child.exitCode === null) {
    try { process.kill(process.platform === 'win32' ? child.pid : -child.pid, 'SIGTERM'); } catch (error) { if (error.code !== 'ESRCH') throw error; }
  }
}
