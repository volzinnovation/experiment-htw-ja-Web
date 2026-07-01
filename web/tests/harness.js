import { WumpusGame } from '../game.js';
import {
  WumpusApp,
  buildChallengeUrl,
  createAppMarkup,
  normalizeSeed,
  seedFromLocation,
} from '../app.js';

const tests = [];

function test(name, fn) {
  tests.push({ name, fn });
}

function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

function assertEquals(actual, expected, message) {
  if (actual !== expected) {
    throw new Error(`${message} (expected ${expected}, got ${actual})`);
  }
}

function makeFixedRandom(values) {
  const queue = [...values];
  return () => (queue.length ? queue.shift() : 0);
}

test('move rejects non-adjacent rooms', () => {
  const game = new WumpusGame({
    setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [6, 7], grenadeRoom: 8 },
    rng: makeFixedRandom([0]),
  });
  const lines = game.executeCommand('m 20');
  assertEquals(lines[0], "CAN'T MOVE THERE", 'invalid move should be rejected');
});

test('moving into a pit loses the game', () => {
  const game = new WumpusGame({
    setup: { player: 1, wumpus: 2, pits: [5, 4], bats: [6, 7], grenadeRoom: 8 },
    rng: makeFixedRandom([0]),
  });
  const lines = game.executeCommand('m 5');
  assert(lines.includes('YYYIIIIEEEE . . . FELL IN PIT'), 'pit message should be present');
  assertEquals(game.snapshot().status, WumpusGame.status.LOST, 'status should be lost');
});

test('shooting Wumpus directly wins', () => {
  const game = new WumpusGame({
    setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [6, 7], grenadeRoom: 8 },
    rng: makeFixedRandom([0]),
  });
  const lines = game.executeCommand('s 2');
  assert(lines.some((line) => line.includes('AHA! YOU GOT THE WUMPUS')), 'win message should be present');
  assertEquals(game.snapshot().status, WumpusGame.status.WON, 'status should be won');
});

test('grenade detonates on next rest turn', () => {
  const game = new WumpusGame({
    setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [6, 7], grenadeRoom: 8 },
    carriesGrenade: true,
    rng: makeFixedRandom([0]),
  });
  const throwLines = game.executeCommand('t 2');
  assert(throwLines.includes('YOU HEAR TIC...TIC...'), 'throw should arm grenade');
  const restLines = game.executeCommand('r');
  assert(restLines.includes('YOU HEAR A HORRENDOUS EXPLOSION!'), 'detonation should occur');
  assertEquals(game.snapshot().status, WumpusGame.status.WON, 'grenade blast should win if wumpus in blast');
});

test('bat room relocates player', () => {
  const game = new WumpusGame({
    setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [5, 7], grenadeRoom: 8 },
    rng: makeFixedRandom([0.45]),
  });
  const lines = game.executeCommand('m 5');
  assert(lines.some((line) => line.includes('SUPER BAT SNATCH')), 'bat message should be present');
  assert(game.snapshot().player !== 5, 'player should no longer be in bat room');
});

test('frontend updates status and log after command', () => {
  const host = document.createElement('div');
  host.innerHTML = createAppMarkup();
  document.body.append(host);

  const app = new WumpusApp(host, {
    gameOptions: {
      setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [6, 7], grenadeRoom: 8 },
      rng: makeFixedRandom([0]),
    },
  });

  app.commandInput.value = 's 2';
  app.commandForm.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

  assert(app.statusNode.textContent.includes('WON'), 'status panel should show win');
  assert(app.logNode.textContent.includes('AHA! YOU GOT THE WUMPUS'), 'log should include win message');
});

test('seeded games produce deterministic setups', () => {
  const first = new WumpusGame({ seed: 1973 }).snapshot();
  const second = new WumpusGame({ seed: 1973 }).snapshot();
  const different = new WumpusGame({ seed: 1974 }).snapshot();

  assertEquals(first.player, second.player, 'same seed should preserve player room');
  assertEquals(first.wumpus, second.wumpus, 'same seed should preserve wumpus room');
  assert(first.player !== different.player || first.wumpus !== different.wumpus, 'different seed should usually change setup');
});

test('challenge URL helpers normalize and preserve seed params', () => {
  assertEquals(normalizeSeed('001973'), 1973, 'numeric seeds should normalize');
  assertEquals(normalizeSeed('nope'), null, 'non-numeric seeds should be rejected');
  assertEquals(seedFromLocation({ search: '?seed=4242' }), 4242, 'seed should be read from URL search params');
  assertEquals(
    buildChallengeUrl(1973, 'https://example.test/play?mode=quiet#rooms'),
    'https://example.test/play?mode=quiet&seed=1973',
    'challenge URL should preserve existing params and clear hash',
  );
});

test('frontend renders seeded challenge and replay trail', () => {
  const host = document.createElement('div');
  host.innerHTML = createAppMarkup();
  document.body.append(host);

  const app = new WumpusApp(host, {
    challengeSeed: 1973,
    gameOptions: {
      setup: { player: 1, wumpus: 2, pits: [3, 4], bats: [6, 7], grenadeRoom: 8 },
      rng: makeFixedRandom([0]),
    },
  });

  assertEquals(app.challengeSeedInput.value, '1973', 'challenge seed should be visible');
  assert(app.challengeSummaryNode.textContent.includes('Seed 1973'), 'challenge summary should mention seed');
  app.commandInput.value = 's 2';
  app.commandForm.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

  assert(app.routeSummaryNode.textContent.includes('Commands: 1'), 'route summary should count commands');
  assert(app.replayLogNode.textContent.includes('s 2'), 'replay log should include command');
  assert(app.challengeLinkNode.textContent.includes('seed=1973'), 'challenge link should include seed');
});

async function run() {
  const results = document.getElementById('results');
  const summary = document.getElementById('summary');

  let passed = 0;
  const failures = [];

  for (const item of tests) {
    const line = document.createElement('li');
    try {
      await item.fn();
      passed += 1;
      line.textContent = `PASS: ${item.name}`;
      line.className = 'pass';
    } catch (error) {
      failures.push({ name: item.name, error });
      line.textContent = `FAIL: ${item.name} -> ${error.message}`;
      line.className = 'fail';
    }
    results.append(line);
  }

  if (failures.length === 0) {
    summary.textContent = `All ${passed} tests passed.`;
    summary.dataset.status = 'pass';
  } else {
    summary.textContent = `${failures.length} / ${tests.length} tests failed.`;
    summary.dataset.status = 'fail';
  }
}

run();
