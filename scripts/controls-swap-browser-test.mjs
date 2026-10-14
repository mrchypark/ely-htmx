import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import {createServer} from 'node:http';
import {chromium} from 'playwright';

// Serve only the regression fixture and production assets it executes.
const files = new Map([
  ['/tests/controls_swap_browser.html', 'text/html; charset=utf-8'],
  ['/tests/controls_swap_response.html', 'text/html; charset=utf-8'],
  ['/assets/htmx.min.js', 'text/javascript; charset=utf-8'],
  ['/assets/ely-controls.js', 'text/javascript; charset=utf-8'],
]);
const server = createServer(async (request, response) => {
  const path = new URL(request.url, 'http://localhost').pathname;
  const type = files.get(path);
  if (!type) { response.writeHead(404).end(); return; }
  try {
    const contents = await readFile(new URL(`..${path}`, import.meta.url));
    response.writeHead(200, {'content-type': type, 'cache-control': 'no-store'}).end(contents);
  } catch (error) {
    response.writeHead(500).end(String(error));
  }
});
let browser;
try {
  await new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(0, '127.0.0.1', resolve);
  });
  browser = await chromium.launch();
  const page = await browser.newPage();
  const errors = [];
  page.on('pageerror', error => errors.push(String(error)));
  await page.goto(`http://127.0.0.1:${server.address().port}/tests/controls_swap_browser.html`);
  await page.waitForFunction(() => window.__controlsSwapResult, null, {timeout: 15000});
  const result = await page.evaluate(() => window.__controlsSwapResult);
  const report = await page.locator('#results').innerText();
  assert.equal(result.passed, true, report);
  assert.deepEqual(errors, [], 'Browser fixture must not throw uncaught errors');
  assert.equal(result.tests, 9, 'Real HTTP, all six swap styles, source replacement and OOB scenarios must run');
  console.log(report);
} finally {
  await browser?.close();
  await new Promise((resolve, reject) => server.close(error => {
    if (error && error.code !== 'ERR_SERVER_NOT_RUNNING') reject(error);
    else resolve();
  }));
}
