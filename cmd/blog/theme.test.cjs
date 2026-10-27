const { test } = require('node:test');
const assert = require('node:assert/strict');
const { readFileSync } = require('node:fs');
const { runInNewContext } = require('node:vm');
const script = readFileSync(`${__dirname}/theme.js`, 'utf8');

function page({ saved, dark = false, blocked = false } = {}) {
  const events = {};
  const media = { matches: dark, addEventListener: (_, callback) => { events.system = callback; } };
  const control = {
    hidden: true, attributes: {},
    setAttribute(name, value) { this.attributes[name] = value; },
    addEventListener: (name, callback) => { events[name] = callback; },
  };
  const classes = new Set();
  const root = { dataset: { colorMode: dark ? 'dark' : 'light' }, classList: {
    toggle(name, enabled) { if (enabled) classes.add(name); else classes.delete(name); },
    contains(name) { return classes.has(name); },
  } };
  const storage = new Map(saved ? [['guigoes.color-mode', saved]] : []);
  const access = () => { if (blocked) throw new Error('storage unavailable'); };
  runInNewContext(script, {
    window: { matchMedia: () => media },
    document: {
      documentElement: root,
      addEventListener: (_, callback) => { events.ready = callback; },
      querySelector: () => control,
    },
    localStorage: {
      getItem: key => { access(); return storage.get(key); },
      setItem: (key, value) => { access(); storage.set(key, value); },
    },
  });
  events.ready();
  return { root, control, storage, click() { events.click(); }, system(dark) {
    media.matches = dark;
    // Margo's earlier listener resolves the operating-system preference first.
    root.dataset.colorMode = dark ? 'dark' : 'light';
    events.system();
  } };
}

test('system default follows OS changes and synchronizes the accessible switch', () => {
  const state = page({ dark: true });
  assert.equal(state.control.hidden, false);
  assert.equal(state.control.attributes['aria-checked'], 'true');
  assert.equal(state.root.dataset.colorMode, 'dark');
  assert.equal(state.root.classList.contains('dark'), true);
  state.system(false);
  assert.equal(state.root.dataset.colorMode, 'light');
  assert.equal(state.root.classList.contains('dark'), false);
  assert.equal(state.control.attributes['aria-checked'], 'false');
  assert.equal(state.storage.size, 0);
});

test('toggle choices survive OS changes and reload', () => {
  const state = page();
  state.click();
  state.system(false);
  assert.equal(state.root.dataset.colorMode, 'dark');
  const reload = page({ saved: state.storage.get('guigoes.color-mode') });
  assert.equal(reload.control.attributes['aria-checked'], 'true');
  assert.equal(reload.root.classList.contains('dark'), true);
  assert.equal(reload.control.title, 'Usar tema claro');
  reload.click();
  reload.system(true);
  assert.equal(reload.root.dataset.colorMode, 'light');
  assert.equal(reload.root.classList.contains('dark'), false);
  assert.equal(reload.control.attributes['aria-checked'], 'false');
  assert.equal(reload.storage.get('guigoes.color-mode'), 'light');
});

test('invalid or unavailable storage falls back to system without disabling the toggle', () => {
  for (const options of [{ saved: 'invalid' }, { blocked: true }]) {
    const state = page(options);
    assert.equal(state.root.dataset.colorMode, 'light');
    state.click();
    assert.equal(state.root.dataset.colorMode, 'dark');
    state.click();
    assert.equal(state.root.dataset.colorMode, 'light');
  }
});
