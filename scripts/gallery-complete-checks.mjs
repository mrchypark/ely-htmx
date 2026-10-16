import assert from 'node:assert/strict';
import {writeFile} from 'node:fs/promises';
import {resolve} from 'node:path';

// This is browser execution, not an HTTP-only crawl. Behavior assertions use the
// shipped synthetic host fixtures; they do not verify real external services.
export async function checkEveryGalleryPage({browser, origin, base, staticMode, storyURLs, output}) {
  const mode = staticMode ? 'static' : 'live';
  const families = ['foundation', 'interaction', 'data', 'workspace', 'communication', 'native', 'forms'];
  const auditURL = (family, aliases) => staticMode
    ? `${origin}${base}components/audit/${family ? family + '/' : ''}${aliases ? 'aliases/' : ''}`
    : `${origin}/components/audit?${new URLSearchParams({...family && {family}, ...aliases && {aliases: '1'}})}`;
  const reports = {mode, scope: 'All gallery pages and synthetic component/alias behavior; no external services or native platform execution.', limitations: ['Host services use synthetic fixtures; callbacks do not prove authentication, payment, device, AI, or persistence integrations.', 'Aliases reuse their resolved canonical behavior checks. Passing every name does not prove every possible prop combination, control path, or alias-specific semantic variation.', ...(staticMode ? ['Markdown preview checks the explicit source-only static fallback; full Go Markdown rendering is tested only in live mode.'] : [])], pages: [], behaviors: []};
  const normalizeURL = url => new URL(url).pathname.replace(/\/index\.html$/, '/').replace(/\/$/, '');
  const stories = new Set(storyURLs.map(normalizeURL));
  assert.equal(stories.size, 624, 'unique story inventory');
  let urls;
  if (staticMode) {
    const response = await fetch(`${origin}${base}gallery-manifest.json`);
    assert(response.ok, 'static manifest loads');
    const manifest = await response.json();
    assert.equal(manifest.stories, 624);
    assert.equal(manifest.pages.length, 646, 'complete export page inventory');
    urls = manifest.pages.map(path => `${origin}${base}${path}`);
  } else {
    urls = [...storyURLs, ...['', 'components', 'components/usage', 'components/themes', 'components/api', 'components/demo'].map(path => `${origin}/${path}`),
      ...['', ...families].flatMap(family => [false, true].map(aliases => auditURL(family, aliases)))];
  }
  assert.equal(urls.length, 646);
  const failures = [];
  const persist = () => writeFile(resolve(output, `${mode}-complete-report.json`), JSON.stringify(reports, null, 2));
  async function session(url, exercise) {
    const context = await browser.newContext({viewport: {width: 1440, height: 1000}, serviceWorkers: 'block'});
    const page = await context.newPage();
    const errors = [];
    page.on('pageerror', error => errors.push(String(error)));
    page.on('response', response => { if (response.status() >= 400) errors.push(`HTTP ${response.status()} ${response.url()}`); });
    page.on('requestfailed', request => {
      if (request.failure()?.errorText !== 'net::ERR_ABORTED') errors.push(`${request.failure()?.errorText} ${request.url()}`);
    });
    await context.route('**/*', route => {
      if (new URL(route.request().url()).origin === origin) return route.continue();
      errors.push(`External dependency: ${route.request().url()}`);
      return route.abort();
    });
    try {
      const response = await page.goto(url, {waitUntil: 'load', timeout: 90_000});
      assert.equal(response.status(), 200, url);
      await page.waitForFunction(() => [...document.querySelectorAll('[data-ely-family]')].every(node => node.closest('[inert]') || node.dataset.elyReady === 'true'), null, {timeout: 30_000});
      await exercise(page);
      await page.waitForTimeout(100); // Surface asynchronous mount/module failures.
      assert.deepEqual(errors, [], `${url}: browser or resource errors`);
      return {url, status: 'passed'};
    } catch (error) {
      const failure = {url, status: 'failed', error: String(error), errors};
      failures.push(failure);
      return failure;
    } finally {
      await context.close();
    }
  }
  try {
    // Four isolated pages bound memory while ensuring every story loads its own
    // dependencies without benefitting from a previous page's module cache.
    const largePages = urls.filter(url => url.includes('/components/audit'));
    const ordinaryPages = urls.filter(url => !url.includes('/components/audit'));
    let next = 0;
    await Promise.all(Array.from({length: 4}, async () => {
      while (next < ordinaryPages.length) {
        const url = ordinaryPages[next++];
        reports.pages.push(await session(url, async page => {
          if (stories.has(normalizeURL(url))) {
            assert(await page.locator('[data-demo-story]').count() > 0, 'story wrapper exists');
            assert(await page.locator('[data-demo-story] [data-ely-family]').count() > 0, 'story renders public components');
          }
        }));
        if (reports.pages.length % 50 === 0) console.log(`[gallery ${mode}] loaded ${reports.pages.length}/646 pages`);
      }
    }));
    for (const url of largePages) reports.pages.push(await session(url, async () => {}));
    for (const family of families) {
      const result = await session(auditURL(family, true), async page => {
        const expected = await page.locator('[data-audit-case]').count();
        assert(expected > 0, `${family} has registered cases`);
        await page.locator('[data-audit-run]').click();
        await page.waitForFunction(() => document.querySelector('[data-audit-results]')?.textContent.trim().startsWith('{'), null, {timeout: 300_000});
        const report = JSON.parse(await page.locator('[data-audit-results]').textContent());
        reports.behaviors.push({family, expected, ...report});
        assert.equal(report.results.length, expected, `${family}: every case ran`);
        assert.equal(report.failed, 0, JSON.stringify(report.results.filter(result => result.status !== 'passed')));
        assert.equal(report.passed, expected);
        assert.equal(report.runtime.status, 'passed', JSON.stringify(report.runtime));
      });
      console.log(`[gallery ${mode}] ${family} behavior: ${result.status}`);
      await persist();
    }
    const names = reports.behaviors.flatMap(report => report.results.map(result => result.name));
    assert.equal(new Set(names).size, 1313, 'every resolved component name has a behavior result');
    assert.equal(names.length, 1313, 'no duplicated behavior coverage');
    assert.deepEqual(failures, [], 'complete gallery browser failures; see complete-report.json');
  } finally {
    await persist();
  }
}
