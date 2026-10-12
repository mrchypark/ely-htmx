(() => {
  const appendTag = input => {
    const value = input.value.trim().replace(/^#/, '');
    const tags = input.closest('.demo-stack')?.querySelector('#forms-tags');
    if (!value || !tags) return;
    const tag = document.createElement('button');
    tag.type = 'button'; tag.className = 'demo-chip'; tag.textContent = `${value} ×`;
    tag.setAttribute('aria-label', `Remove ${value}`); tags.append(tag); input.value = '';
  };
  const localAccountMessages = {
    'loginform-oauthbuttons': 'Local sign-in validation completed. No credentials were sent.',
    signupform: 'Local sign-up validation completed. No account was created.',
    'magiclinkform-forgotpassword': 'Local link preview completed. No email was sent.',
    'profilecard-profileeditor': 'Local profile preview saved. No profile was changed.',
    'onboardingwizard-welcometour': 'Local onboarding preview saved. The sample does not advance.',
    'feedbackwidget-contactsupport': 'Local feedback preview saved. No support request was sent.'
  };
  document.querySelectorAll('[data-forms-story]').forEach(story => {
    const message = localAccountMessages[story.dataset.formsStory];
    const form = story.querySelector('form');
    if (!message || !form) return;
    form.dataset.formsLocal = message;
    const output = document.createElement('output');
    output.dataset.formsLocalOutput = '';
    output.className = 'demo-muted';
    output.setAttribute('role', 'status');
    form.append(output);
  });
  document.addEventListener('keydown', event => {
    const active = document.activeElement;
    if (active?.matches('[data-forms-tag-input]') && event.key === 'Enter') { event.preventDefault(); appendTag(active); }
    const recorder = document.querySelector('[data-forms-hotkey][data-recording]');
    if (!recorder || event.key === 'Tab') return;
    if (event.key === 'Escape') {
      event.preventDefault();
      recorder.removeAttribute('data-recording'); recorder.textContent = 'Record shortcut';
      return;
    }
    if (['Control', 'Meta', 'Alt', 'Shift'].includes(event.key)) return;
    event.preventDefault();
    const keys = [event.metaKey && '⌘', event.ctrlKey && 'Ctrl', event.altKey && '⌥', event.shiftKey && '⇧', event.key.length === 1 ? event.key.toUpperCase() : event.key].filter(Boolean);
    const output = recorder.closest('.demo-stack')?.querySelector('#forms-hotkey-output');
    if (output) output.textContent = keys.join(' + ');
    recorder.removeAttribute('data-recording'); recorder.textContent = 'Record shortcut';
  });
  document.addEventListener('click', event => {
    if (event.target.closest('[data-clear="#forms-signature"]')) { document.querySelector('#forms-signature')?.replaceChildren(); return; }
    const target = event.target.closest('[data-forms-hotkey],[data-forms-add-row],#forms-tags .demo-chip');
    if (!target) return;
    if (target.matches('[data-forms-hotkey]')) { target.toggleAttribute('data-recording'); target.textContent = target.hasAttribute('data-recording') ? 'Press a key…' : 'Record shortcut'; return; }
    if (target.matches('#forms-tags .demo-chip')) { target.remove(); return; }
    const stack = target.closest('.demo-stack');
    const row = stack?.querySelector('template[data-forms-row]');
    if (row) target.before(row.content.cloneNode(true));
  });
  document.addEventListener('change', event => {
    const input = event.target;
    if (!(input instanceof HTMLInputElement) || input.type !== 'file') return;
    const host = input.closest('.demo-dropzone');
    const output = host?.querySelector('[data-forms-file-result]');
    if (output && input.files?.length) output.textContent = `${input.files.length} file${input.files.length === 1 ? '' : 's'} selected locally`;
  });
  document.addEventListener('dragover', event => { if (event.target.closest('.demo-dropzone')) event.preventDefault(); });
  document.addEventListener('drop', event => {
    const zone = event.target.closest('.demo-dropzone');
    if (!zone) return;
    event.preventDefault(); const count = event.dataTransfer?.files.length || 0;
    const output = zone.querySelector('[data-forms-file-result]');
    if (output) output.textContent = count ? `${count} file${count === 1 ? '' : 's'} selected locally` : 'Drop files here';
  });
  document.addEventListener('submit', event => {
    const form = event.target.closest('[data-forms-local]');
    if (!form) return;
    event.preventDefault();
    if (!form.checkValidity()) { form.reportValidity(); return; }
    const output = form.querySelector('[data-forms-local-output]');
    if (output) output.textContent = form.dataset.formsLocal;
  });
})();
