import { WumpusGame } from './game.js';
import { ROOM_ARTWORK } from './artwork.js';

function escapeHtml(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;');
}

export function createAppMarkup() {
  return `
    <section class="panel status-panel">
      <h1>Hunt the Wumpus (Browser Edition)</h1>
      <p class="subtitle">Static HTML/CSS/JS implementation with no server-side runtime.</p>
      <div id="status"></div>
      <div id="warnings" class="warnings"></div>
      <div id="event-banner" class="event-banner" aria-live="polite"></div>
    </section>

    <section class="panel controls-panel">
      <form id="command-form" class="command-form">
        <label for="command-input">Command</label>
        <input id="command-input" name="command" type="text" autocomplete="off" placeholder="m 5, s 2 3, t 8, r" />
        <button type="submit">Execute</button>
      </form>
      <div class="quick-buttons">
        <button type="button" data-cmd="rest">Rest</button>
        <button type="button" data-cmd="new">New Game</button>
      </div>
      <ul class="help-list">
        <li><code>m &lt;room&gt;</code> move to adjacent room</li>
        <li><code>s &lt;r1&gt; ... &lt;rN&gt;</code> shoot through 1-5 rooms</li>
        <li><code>t &lt;room&gt;</code> throw grenade when carrying it</li>
        <li><code>r</code> rest one turn</li>
      </ul>

      <div class="challenge-box">
        <div class="challenge-heading">
          <div>
            <h2>Seeded Challenge</h2>
            <p id="challenge-summary" class="muted"></p>
          </div>
          <span id="challenge-difficulty" class="difficulty-chip">Scoutable</span>
        </div>
        <form id="challenge-form" class="challenge-form">
          <label for="challenge-seed">Seed</label>
          <input id="challenge-seed" name="seed" type="text" inputmode="numeric" autocomplete="off" />
          <button type="submit">Start Seed</button>
        </form>
        <div class="challenge-actions">
          <button type="button" id="random-challenge">Lucky Cave</button>
          <button type="button" id="copy-challenge">Copy Link</button>
        </div>
        <p id="challenge-link" class="challenge-link"></p>
      </div>
    </section>

    <section class="panel log-panel">
      <h2>Turn Log</h2>
      <ol id="log" aria-live="polite"></ol>
    </section>

    <section class="panel replay-panel">
      <h2>Replay Trail</h2>
      <p id="route-summary" class="muted"></p>
      <ol id="replay-log"></ol>
    </section>

    <section class="panel room-panel">
      <h2>Room Map</h2>
      <p class="muted">Room artwork by room type (safe, player, Wumpus, pits, bats, grenade).</p>
      <div id="rooms" class="rooms"></div>
    </section>

    <section class="panel artwork-panel">
      <h2>Artwork Gallery</h2>
      <div id="artwork-gallery" class="artwork-gallery"></div>
    </section>
  `;
}

export function normalizeSeed(value) {
  const seed = Number.parseInt(String(value ?? '').trim(), 10);
  if (!Number.isFinite(seed) || seed < 1) {
    return null;
  }
  return Math.min(seed, 999999999);
}

export function randomChallengeSeed() {
  return Math.floor(Math.random() * 999999999) + 1;
}

export function seedFromLocation(locationObject = globalThis.location ?? { search: '' }) {
  const params = new URLSearchParams(locationObject.search || '');
  return normalizeSeed(params.get('seed'));
}

export function buildChallengeUrl(seed, href = globalThis.location?.href ?? 'http://localhost/') {
  const url = new URL(href);
  url.searchParams.set('seed', String(seed));
  url.hash = '';
  return url.toString();
}

export class WumpusApp {
  constructor(root, options = {}) {
    this.root = root;
    this.options = options;
    this.turn = 1;
    this.game = null;
    this.challengeSeed = normalizeSeed(options.challengeSeed) || seedFromLocation() || randomChallengeSeed();
    this.commandHistory = [];
    this.roomTrail = [];

    if (!this.root.innerHTML.trim()) {
      this.root.innerHTML = createAppMarkup();
    }

    this.statusNode = this.root.querySelector('#status');
    this.warningsNode = this.root.querySelector('#warnings');
    this.logNode = this.root.querySelector('#log');
    this.roomsNode = this.root.querySelector('#rooms');
    this.eventNode = this.root.querySelector('#event-banner');
    this.commandForm = this.root.querySelector('#command-form');
    this.commandInput = this.root.querySelector('#command-input');
    this.galleryNode = this.root.querySelector('#artwork-gallery');
    this.challengeForm = this.root.querySelector('#challenge-form');
    this.challengeSeedInput = this.root.querySelector('#challenge-seed');
    this.challengeSummaryNode = this.root.querySelector('#challenge-summary');
    this.challengeDifficultyNode = this.root.querySelector('#challenge-difficulty');
    this.challengeLinkNode = this.root.querySelector('#challenge-link');
    this.randomChallengeButton = this.root.querySelector('#random-challenge');
    this.copyChallengeButton = this.root.querySelector('#copy-challenge');
    this.routeSummaryNode = this.root.querySelector('#route-summary');
    this.replayLogNode = this.root.querySelector('#replay-log');
    this.challengeSeedInput.value = String(this.challengeSeed);

    this.bindEvents();
    this.renderGallery();
    this.startNewGame();
  }

  bindEvents() {
    this.commandForm.addEventListener('submit', (event) => {
      event.preventDefault();
      const command = this.commandInput.value;
      if (!command.trim()) {
        return;
      }
      this.execute(command);
      this.commandInput.value = '';
      this.commandInput.focus();
    });

    this.root.querySelectorAll('.quick-buttons button').forEach((button) => {
      button.addEventListener('click', () => {
        if (button.dataset.cmd === 'new') {
          this.startNewGame();
          return;
        }
        this.execute('r');
      });
    });

    this.challengeForm.addEventListener('submit', (event) => {
      event.preventDefault();
      const seed = normalizeSeed(this.challengeSeedInput.value);
      if (!seed) {
        this.showBanner('event-lose', 'Seed must be a positive number');
        return;
      }
      this.setChallengeSeed(seed, true);
      this.startNewGame();
    });

    this.randomChallengeButton.addEventListener('click', () => {
      this.setChallengeSeed(randomChallengeSeed(), true);
      this.startNewGame();
    });

    this.copyChallengeButton.addEventListener('click', () => {
      this.copyChallengeLink();
    });
  }

  startNewGame() {
    const gameOptions = { ...(this.options.gameOptions ?? {}) };
    if (!gameOptions.setup && !gameOptions.rng && gameOptions.seed === undefined) {
      gameOptions.seed = this.challengeSeed;
    }
    this.game = new WumpusGame(gameOptions);
    this.turn = 1;
    this.commandHistory = [];
    this.roomTrail = [this.game.snapshot().player];
    this.logNode.innerHTML = '';
    this.replayLogNode.innerHTML = '';
    this.updateView([]);
  }

  execute(command) {
    if (this.game.snapshot().status !== WumpusGame.status.IN_PROGRESS) {
      this.pushLog('Game over. Start a new game to continue.');
      return;
    }

    const before = this.game.snapshot();
    const lines = this.game.executeCommand(command);
    const after = this.game.snapshot();
    this.pushLog(`> ${command}`);
    if (lines.length === 0) {
      this.pushLog('(no visible effect)');
    }
    lines.forEach((line) => this.pushLog(line));
    this.recordCommand(command, before, after, lines);

    this.turn += 1;
    this.updateView(lines, command);
  }

  pushLog(text) {
    const li = document.createElement('li');
    li.textContent = text;
    this.logNode.prepend(li);
  }

  updateView(lastLines = [], command = '') {
    const state = this.game.snapshot();
    const exits = WumpusGame.caveExits[state.player];
    const prompt = state.carriesGrenade ? 'SHOOT, MOVE OR THROW (S-M-T)?' : 'SHOOT OR MOVE (S-M)?';

    this.statusNode.innerHTML = `
      <p><strong>Turn:</strong> ${this.turn}</p>
      <p><strong>Status:</strong> ${escapeHtml(state.status.toUpperCase())}</p>
      <p><strong>You are in room:</strong> ${state.player}</p>
      <p><strong>Tunnels lead to:</strong> ${exits.join(', ')}</p>
      <p><strong>Arrows left:</strong> ${state.arrows}</p>
      <p><strong>Prompt:</strong> ${prompt}</p>
    `;

    const warnings = this.game.turnWarnings();
    this.warningsNode.innerHTML = warnings.length
      ? warnings.map((warning) => `<span class="warning-chip">${escapeHtml(warning)}</span>`).join('')
      : '<span class="warning-chip muted">No immediate warnings</span>';

    this.renderRooms(state);
    this.playEventAnimation(lastLines, command);
    this.renderChallenge(state);
    this.renderReplay(state);
  }

  renderRooms(state) {
    const cards = [];
    for (let room = 1; room <= 20; room += 1) {
      const type = this.roomType(state, room);
      const art = ROOM_ARTWORK[type];
      cards.push(`
        <article class="room-card type-${type}" data-room="${room}">
          <header>Room ${room}</header>
          <div class="icon">${art.art}</div>
          <p>${art.label}</p>
        </article>
      `);
    }
    this.roomsNode.innerHTML = cards.join('');
  }

  renderGallery() {
    const cards = Object.entries(ROOM_ARTWORK).map(([key, value]) => `
      <article class="gallery-card type-${key}">
        <div class="icon">${value.art}</div>
        <p>${value.label}</p>
      </article>
    `);
    this.galleryNode.innerHTML = cards.join('');
  }

  roomType(state, room) {
    if (state.player === room) {
      return 'player';
    }
    if (state.wumpus === room) {
      return 'wumpus';
    }
    if (state.pits.includes(room)) {
      return 'pit';
    }
    if (state.bats.includes(room)) {
      return 'bats';
    }
    if (state.grenadeRoom === room) {
      return 'grenade';
    }
    return 'safe';
  }

  playEventAnimation(lines, command) {
    const text = lines.join(' ');
    let kind = 'event-neutral';
    let label = 'Turn updated';

    if (text.includes('AHA! YOU GOT THE WUMPUS')) {
      kind = 'event-win';
      label = 'Wumpus defeated';
    } else if (text.includes('HA HA HA - YOU LOSE!')) {
      kind = 'event-lose';
      label = 'You lost the round';
    } else if (text.includes('HORRENDOUS EXPLOSION')) {
      kind = 'event-explosion';
      label = 'Grenade explosion';
    } else if (text.includes('SUPER BAT SNATCH')) {
      kind = 'event-bats';
      label = 'Bat relocation';
    } else if (text.includes('MISSED')) {
      kind = 'event-miss';
      label = 'Arrow missed';
    } else if (String(command).toLowerCase().startsWith('m ')) {
      kind = 'event-move';
      label = 'Player moved';
    }

    this.eventNode.className = `event-banner ${kind}`;
    this.eventNode.textContent = label;
  }

  setChallengeSeed(seed, syncUrl = false) {
    this.challengeSeed = seed;
    this.challengeSeedInput.value = String(seed);
    if (syncUrl && window.history?.replaceState) {
      window.history.replaceState({}, '', buildChallengeUrl(seed));
    }
  }

  challengeUrl() {
    return buildChallengeUrl(this.challengeSeed);
  }

  async copyChallengeLink() {
    const url = this.challengeUrl();
    this.challengeLinkNode.textContent = url;
    try {
      if (!navigator.clipboard?.writeText) {
        throw new Error('Clipboard API unavailable');
      }
      await navigator.clipboard.writeText(url);
      this.showBanner('event-win', 'Challenge link copied');
    } catch (error) {
      this.showBanner('event-neutral', 'Challenge link ready');
    }
  }

  recordCommand(command, before, after, lines) {
    if (after.player !== this.roomTrail.at(-1)) {
      this.roomTrail.push(after.player);
    }
    this.commandHistory.push({
      turn: this.turn,
      command: String(command).trim(),
      before: before.player,
      after: after.player,
      status: after.status,
      message: lines[0] || '(no visible effect)',
    });
  }

  renderChallenge(state) {
    const warnings = this.game.turnWarnings();
    const difficulty = this.challengeDifficulty(state, warnings);
    this.challengeDifficultyNode.textContent = difficulty;
    this.challengeDifficultyNode.dataset.level = difficulty.toLowerCase();
    this.challengeSummaryNode.textContent = `Seed ${this.challengeSeed} starts in room ${this.roomTrail[0]} and can be shared as a permalink.`;
    this.challengeLinkNode.textContent = this.challengeUrl();
  }

  renderReplay(state) {
    const route = this.roomTrail.join(' -> ');
    const commandCount = this.commandHistory.length;
    this.routeSummaryNode.textContent = `Rooms visited: ${route}. Commands: ${commandCount}. Result: ${state.status}.`;
    this.replayLogNode.innerHTML = this.commandHistory
      .slice(-8)
      .reverse()
      .map((item) => `
        <li>
          <span>Turn ${item.turn}: <code>${escapeHtml(item.command)}</code></span>
          <small>${item.before} -> ${item.after} · ${escapeHtml(item.message)}</small>
        </li>
      `)
      .join('');
  }

  challengeDifficulty(state, warnings) {
    let score = warnings.length;
    if (warnings.includes('I SMELL A WUMPUS')) {
      score += 2;
    }
    if (warnings.includes('I FEEL A DRAFT')) {
      score += 1;
    }
    if (state.carriesGrenade) {
      score -= 1;
    }
    if (state.grenadeRoom && WumpusGame.caveExits[state.player].includes(state.grenadeRoom)) {
      score -= 1;
    }
    if (score >= 4) return 'Brutal';
    if (score >= 2) return 'Tense';
    return 'Scoutable';
  }

  showBanner(kind, label) {
    this.eventNode.className = `event-banner ${kind}`;
    this.eventNode.textContent = label;
  }
}
