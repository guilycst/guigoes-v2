(() => {
  const root = document.documentElement;
  const system = window.matchMedia('(prefers-color-scheme: dark)');
  let preference = 'system';
  try {
    const saved = localStorage.getItem('guigoes.color-mode');
    if (saved === 'light' || saved === 'dark') preference = saved;
  } catch (_) { /* Storage may be unavailable in private browsing. */ }

  function apply() {
    root.dataset.colorMode = preference === 'system'
      ? (system.matches ? 'dark' : 'light') : preference;
  }
  apply();
  // Registered after Margo's system listener so an explicit choice stays active.
  system.addEventListener('change', apply);
  document.addEventListener('DOMContentLoaded', () => {
    const control = document.querySelector('[data-guigoes-color-mode]');
    control.value = preference;
    control.hidden = false;
    control.addEventListener('change', () => {
      preference = control.value;
      apply();
      try {
        if (preference === 'system') localStorage.removeItem('guigoes.color-mode');
        else localStorage.setItem('guigoes.color-mode', preference);
      } catch (_) { /* The selected mode still works for this page. */ }
    });
  });
})();
