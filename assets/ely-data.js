document.addEventListener('click', (event) => {
  const tag = event.target.closest('[data-data-tags] [data-data-tag]');
  if (tag) { tag.hidden = true; return; }
  const selectable = event.target.closest('[data-data-selectable]');
  if (selectable) { selectListOption(selectable); return; }
  const range = event.target.closest('[data-data-range]');
  if (range) {
    const root = range.closest('.demo-stack') || document;
    root.dataset.dataChartRange = range.dataset.dataRange;
    root.querySelectorAll('[data-data-range]').forEach((item) => {
      item.setAttribute('aria-pressed', String(item === range));
    });
    updateChartStatus(root);
    return;
  }
  const draw = event.target.closest('[data-data-draw]');
  if (draw) {
    const root = draw.closest('.demo-drawing');
    root.querySelectorAll('[data-data-draw]').forEach((item) => {
      item.setAttribute('aria-pressed', String(item === draw));
    });
    root.querySelector('[data-data-draw-line]').toggleAttribute('hidden', draw.dataset.dataDraw !== 'line');
    root.querySelector('[data-data-draw-note]').toggleAttribute('hidden', draw.dataset.dataDraw !== 'note');
    return;
  }
  const action = event.target.closest('[data-data-action]');
  if (!action) return;
  const root = action.closest('.demo-stack') || document;
  const name = action.dataset.dataAction;
  if (name === 'increment') {
    const count = root.querySelector('[data-data-count]');
    count.textContent = String(Number(count.textContent) + 1);
  }
  if (name === 'reset-tags') root.querySelectorAll('[data-data-tags] .demo-chip').forEach((tag) => { tag.hidden = false; });
  if (name === 'remove-row') action.closest('div').remove();
  if (name === 'load-more') { action.insertAdjacentHTML('beforebegin', '<div>이전 기록 · 09:48</div><div>이전 기록 · 09:40</div>'); action.remove(); }
  if (name === 'append-feed') root.querySelector('[data-data-feed]').textContent = '새 활동 · 연구 질문에 메모를 남겼습니다.';
  if (name === 'add-label') {
    const input = root.querySelector('[data-data-label-input]');
    const label = input.value.trim();
    if (label) {
      const chip = document.createElement('span');
      chip.className = 'demo-chip';
      chip.textContent = `${label} ×`;
      input.before(chip);
      input.value = '';
    }
  }
  if (name === 'move-row') { const row = action.closest('li'); const next = row && row.nextElementSibling; if (row && next) next.after(row); }
  if (name === 'table-columns') { const table = root.querySelector('[data-data-table] .demo-table'); if (table) { table.classList.toggle('is-updated-hidden'); action.textContent = table.classList.contains('is-updated-hidden') ? '갱신 열 보이기' : '갱신 열 숨기기'; } }
  if (name === 'table-density') { const table = root.querySelector('[data-data-table] .demo-table'); if (table) { table.classList.toggle('is-dense'); action.textContent = table.classList.contains('is-dense') ? '여유 있게 보기' : '촘촘히 보기'; } }
  if (name === 'timer') { const timer = root.querySelector('[data-data-timer]'); timer.textContent = timer.textContent === '01:35:00' ? '01:34:59' : '01:35:00'; action.textContent = action.textContent === '시작' ? '중지' : '시작'; }
  if (name === 'pomodoro') { const timer = root.querySelector('[data-data-pomodoro]'); timer.textContent = timer.textContent === '25:00' ? '24:59' : '25:00'; action.textContent = action.textContent === '시작' ? '일시 정지' : '시작'; }
  if (name === 'ack-alert') { action.closest('article').querySelector('span').textContent = '확인됨 · 로컬 예제'; action.remove(); }
  if (name === 'pause-stream') { const output = root.querySelector('[data-data-stream-state]'); const paused = action.textContent === '일시 정지'; action.textContent = paused ? '계속' : '일시 정지'; output.textContent = paused ? '스트림을 일시 정지했습니다.' : '스트림을 다시 표시합니다.'; }
  if (name === 'realtime') root.querySelector('[data-data-realtime]').textContent = '고정 표본 9개 · 새 표본 1개 추가됨';
  if (name === 'carousel-next' || name === 'carousel-prev') {
    const images = ['dunes-sun.jpg', 'atrium-stair.jpg', 'dunes-grass.jpg']; const image = root.querySelector('[data-data-slide]'); const index = images.indexOf(image.src.split('/').pop()); const next = (index + (name === 'carousel-next' ? 1 : 2)) % images.length; image.src = `/assets/ely/assets/${images[next]}`; root.querySelector('[data-data-slide-label]').textContent = `${next + 1} / 3`;
  }
});
document.addEventListener('change', (event) => {
  if (event.target.matches('[data-data-file]')) event.target.closest('.demo-row').querySelector('[data-data-file-name]').textContent = event.target.files[0] ? `${event.target.files[0].name} · 로컬 선택` : '사진을 선택하지 않았습니다.';
  if (event.target.matches('[data-data-status]')) event.target.closest('.demo-stack').querySelector('[data-data-status-output]').textContent = `${event.target.value} · 민서 · 2026-10-14`;
  if (event.target.matches('[data-data-dashboard-filter]')) event.target.closest('.demo-stack').querySelector('[data-data-dashboard-state]').textContent = `${event.target.value} · web · 고정 표본`;
  if (event.target.matches('[data-data-chart-type]')) updateChartStatus(event.target.closest('.demo-stack') || document);
  if (event.target.matches('[data-data-table-status]')) filterDataTable(event.target.closest('.demo-stack') || document);
});

function updateChartStatus(root) {
  const status = root.querySelector('[data-data-chart-status]');
  if (!status) return;
  const range = root.dataset.dataChartRange || root.querySelector('[data-data-range]')?.dataset.dataRange || '1일';
  status.textContent = `${range} · ${root.querySelector('[data-data-chart-type]').value}`;
}

function selectListOption(option) {
  const listbox = option.closest('[role="listbox"]');
  listbox.querySelectorAll('[data-data-selectable]').forEach((item) => {
    const selected = item === option;
    item.setAttribute('aria-selected', String(selected));
    item.tabIndex = selected ? 0 : -1;
  });
  option.focus();
}

function filterDataTable(root) {
  const queryInput = root.querySelector('[data-data-table-filter]');
  const statusInput = root.querySelector('[data-data-table-status]');
  const table = root.querySelector('[data-data-table] .demo-table');
  const summary = root.querySelector('[data-data-table-summary]');
  if (!table) return;
  const query = queryInput ? queryInput.value.trim().toLocaleLowerCase('ko-KR') : '';
  const status = statusInput ? statusInput.value : '';
  let visible = 0;
  table.querySelectorAll('tbody tr').forEach((row) => {
    const text = row.textContent.toLocaleLowerCase('ko-KR');
    const matches = (!query || text.includes(query)) && (!status || row.textContent.includes(status));
    row.hidden = !matches;
    if (matches) visible += 1;
  });
  if (summary) summary.textContent = `표시 ${visible}개 · 평균 갱신 1.5일`;
}

document.addEventListener('input', (event) => {
  if (event.target.matches('[data-data-table-filter]')) filterDataTable(event.target.closest('.demo-stack') || document);
});

document.addEventListener('keydown', (event) => {
  const option = event.target.closest('[data-data-selectable]');
  if (!option || !['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(event.key)) return;
  const options = [...option.closest('[role="listbox"]').querySelectorAll('[data-data-selectable]')];
  const index = options.indexOf(option);
  const next = event.key === 'Home' ? 0 : event.key === 'End' ? options.length - 1 : (index + (event.key === 'ArrowDown' ? 1 : options.length - 1)) % options.length;
  event.preventDefault();
  selectListOption(options[next]);
});
