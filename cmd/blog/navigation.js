document.addEventListener('DOMContentLoaded', () => {
  const nav = document.querySelector('.guigoes-nav');
  const links = document.querySelector('.guigoes-nav-links');
  const trigger = document.querySelector('.guigoes-menu-trigger');
  const dialog = document.querySelector('#guigoes-menu');
  const mobile = window.matchMedia('(max-width: 700px)');
  if (typeof dialog.showModal !== 'function') return;

  function placeLinks() {
    if (mobile.matches) {
      dialog.append(links);
    } else {
      const wasOpen = dialog.open;
      dialog.close();
      nav.insertBefore(links, nav.querySelector('[data-guigoes-color-mode]'));
      if (wasOpen) links.querySelector('a').focus();
    }
  }
  placeLinks();
  mobile.addEventListener('change', placeLinks);
  trigger.hidden = false;
  document.documentElement.classList.add('guigoes-menu-ready');

  trigger.addEventListener('click', () => {
    dialog.showModal();
    trigger.setAttribute('aria-expanded', 'true');
    document.documentElement.classList.add('guigoes-menu-open');
  });
  dialog.querySelector('button').addEventListener('click', () => dialog.close());
  dialog.addEventListener('click', event => {
    const bounds = dialog.getBoundingClientRect();
    if (event.target === dialog && (event.clientX < bounds.left ||
        event.clientX > bounds.right || event.clientY < bounds.top ||
        event.clientY > bounds.bottom)) dialog.close();
    if (event.target.closest('a')) dialog.close();
  });
  dialog.addEventListener('keydown', event => {
    if (event.key !== 'Tab') return;
    const items = dialog.querySelectorAll('button, a[href]');
    const first = items[0], last = items[items.length - 1];
    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  });
  dialog.addEventListener('close', () => {
    trigger.setAttribute('aria-expanded', 'false');
    document.documentElement.classList.remove('guigoes-menu-open');
  });
});
