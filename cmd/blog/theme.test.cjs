const { test } = require('node:test');
const assert = require('node:assert/strict');
const { readFileSync } = require('node:fs');
const { runInNewContext } = require('node:vm');
const script = readFileSync(`${__dirname}/theme.js`, 'utf8');

function page({ saved, dark = false, blocked = false } = {}) {
  const events = {};
  const media = { matches: dark, addEventListener: (_, callback) => { events.system = callback; } };
  const control = { hidden: true, value: 'system', addEventListener: (_, callback) => { events.change = callback; } };
  const root = { dataset: { colorMode: dark ? 'dark' : 'light' } };
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
      removeItem: key => { access(); storage.delete(key); },
    },
  });
  events.ready();
  return { root, control, storage, choose(value) { control.value = value; events.change(); }, system(dark) {
    media.matches = dark;
    // Margo's earlier listener resolves the operating-system preference first.
    root.dataset.colorMode = dark ? 'dark' : 'light';
    events.system();
  } };
}

test('system default follows OS changes and shows the selector', () => {
  const state = page({ dark: true });
  assert.equal(state.control.hidden, false);
  assert.equal(state.control.value, 'system');
  assert.equal(state.root.dataset.colorMode, 'dark');
  state.system(false);
  assert.equal(state.root.dataset.colorMode, 'light');
});

test('explicit choices survive OS changes and reload; system clears persistence', () => {
  const state = page();
  state.choose('dark');
  state.system(false);
  assert.equal(state.root.dataset.colorMode, 'dark');
  const reload = page({ saved: state.storage.get('guigoes.color-mode') });
  assert.equal(reload.control.value, 'dark');
  assert.equal(reload.root.dataset.colorMode, 'dark');
  state.choose('light');
  state.system(true);
  assert.equal(state.root.dataset.colorMode, 'light');
  state.choose('system');
  assert.equal(state.root.dataset.colorMode, 'dark');
  assert.equal(state.storage.has('guigoes.color-mode'), false);
});

test('invalid or unavailable storage falls back to system without disabling controls', () => {
  for (const options of [{ saved: 'invalid' }, { blocked: true }]) {
    const state = page(options);
    assert.equal(state.control.value, 'system');
    state.choose('dark');
    assert.equal(state.root.dataset.colorMode, 'dark');
    state.choose('system');
    assert.equal(state.root.dataset.colorMode, 'light');
  }
});
