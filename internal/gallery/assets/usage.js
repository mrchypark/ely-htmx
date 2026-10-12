document.querySelector('[data-usage-form]')?.addEventListener('submit', event => {
  event.preventDefault();
  const entries = [...new FormData(event.currentTarget)].map(([name,value]) => [name, value instanceof File ? value.name : value]);
  event.currentTarget.querySelector('[data-usage-values]').textContent = JSON.stringify(entries, null, 2);
});
