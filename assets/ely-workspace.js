(() => {
  document.addEventListener('click', event => {
    const file = event.target.closest('[data-workspace-file]');
    if (file) document.querySelectorAll('[id^="workspace-file-"]').forEach(node => { node.hidden = node.id !== `workspace-file-${file.dataset.workspaceFile}`; });
    const zoom = event.target.closest('button[data-workspace-zoom]');
    if (zoom) { const output = zoom.closest('.demo-workspace').querySelector('output[data-workspace-zoom]'); const next = Math.max(50, Math.min(200, Number.parseInt(output.textContent, 10) + (zoom.dataset.workspaceZoom === 'in' ? 10 : -10))); output.textContent = `${next}%`; }
    if (event.target.matches('[data-workspace-query-run]')) document.querySelectorAll('[data-workspace-query-result]').forEach(output => { output.textContent = 'Mock result · 2 rows · 4 ms · no network request'; });
  });
  document.addEventListener('submit', event => {
    const form = event.target.closest('[data-workspace-terminal]'); if (!form) return;
    event.preventDefault(); const command = form.querySelector('input').value.trim(); const output = document.querySelector('[data-workspace-terminal-output]'); if (command) output.textContent += `\n$ ${command}\nmock: command completed locally`;
  });
  document.addEventListener('input', event => {
    if (event.target.matches('[data-workspace-filter]')) { const term = event.target.value.toLowerCase(); document.querySelectorAll('[data-workspace-files] li,.demo-results li').forEach(row => { row.hidden = Boolean(term) && !row.textContent.toLowerCase().includes(term); }); }
    if (event.target.matches('[data-workspace-console]')) { const output = document.querySelector('[data-workspace-console-output]'); if (output) output.textContent = event.target.value ? `mock result: ${event.target.value}` : 'Ready.'; }
  });
})();
