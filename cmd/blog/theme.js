(() => {
  const root = document.documentElement;
  const system = window.matchMedia('(prefers-color-scheme: dark)');
  let preference = 'system';
  let control;
  try {
    const saved = localStorage.getItem('guigoes.color-mode');
    if (saved === 'light' || saved === 'dark') preference = saved;
  } catch (_) { /* Storage may be unavailable in private browsing. */ }

  function apply() {
    root.dataset.colorMode = preference === 'system'
      ? (system.matches ? 'dark' : 'light') : preference;
    if (control) {
      const dark = root.dataset.colorMode === 'dark';
      control.setAttribute('aria-checked', String(dark));
      control.title = dark ? 'Usar tema claro' : 'Usar tema escuro';
    }
  }
  apply();
  // Registered after Margo's system listener so an explicit choice stays active.
  system.addEventListener('change', apply);
  document.addEventListener('DOMContentLoaded', () => {
    control = document.querySelector('[data-guigoes-color-mode]');
    apply();
    control.hidden = false;
    control.addEventListener('click', () => {
      preference = root.dataset.colorMode === 'dark' ? 'light' : 'dark';
      apply();
      try {
        localStorage.setItem('guigoes.color-mode', preference);
      } catch (_) { /* The selected mode still works for this page. */ }
    });
  });
})();
