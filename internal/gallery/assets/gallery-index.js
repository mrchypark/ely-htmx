(() => {
  const root = document.querySelector('[data-gallery-index]');
  if (!root) return;
  const form = root.querySelector('.gallery-filters');
  const query = root.querySelector('#component-query');
  const category = root.querySelector('#gallery-category');
  const cards = [...root.querySelectorAll('[data-gallery-card]')];
  const normal = value => value.normalize('NFKC').trim().toLocaleLowerCase();
  function filter(updateURL = true) {
    const q = normal(query.value);
    let chapters = 0, stories = 0;
    for (const card of cards) {
      const categoryMatch = !category.value || card.dataset.category === category.value;
      const titleMatch = normal(card.querySelector('.chapter-title').textContent).includes(q);
      let matches = 0;
      for (const link of card.querySelectorAll('.chapter-stories a')) {
        link.hidden = !(categoryMatch && (titleMatch || normal(link.dataset.search).includes(q)));
        if (!link.hidden) matches++;
      }
      card.hidden = matches === 0;
      card.querySelector('[data-gallery-examples]').open = Boolean(q && matches);
      if (matches) { chapters++; stories += matches; }
    }
    root.querySelector('#gallery-results').textContent = `${chapters}개 분류 · ${stories}개 예제${q ? ' 검색 결과' : ''}`;
    root.querySelector('#gallery-empty').hidden = chapters !== 0;
    if (updateURL) {
      const url = new URL(location.href);
      query.value.trim() ? url.searchParams.set('q', query.value.trim()) : url.searchParams.delete('q');
      category.value ? url.searchParams.set('category', category.value) : url.searchParams.delete('category');
      history.replaceState(null, '', url);
    }
  }
  function restore() {
    const params = new URLSearchParams(location.search);
    query.value = (params.get('q') || '').slice(0, 200);
    category.value = params.get('category') || '';
    filter(false);
  }
  function reset() { query.value = ''; category.value = ''; filter(); query.focus(); }
  form.addEventListener('input', () => filter());
  form.addEventListener('change', () => filter());
  form.addEventListener('submit', event => { event.preventDefault(); filter(); });
  form.addEventListener('reset', event => { event.preventDefault(); reset(); });
  root.querySelector('[data-gallery-clear]').addEventListener('click', reset);
  addEventListener('popstate', restore);
  addEventListener('pageshow', restore);
  restore();
})();
