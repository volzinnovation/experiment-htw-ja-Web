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
    </section>

    <section class="panel log-panel">
      <h2>Turn Log</h2>
      <ol id="log" aria-live="polite"></ol>
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

export class WumpusApp {
  constructor(root, options = {}) {
    this.root = root;
    this.options = options;
    this.turn = 1;
    this.game = null;

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
  }

  startNewGame() {
    const gameOptions = this.options.gameOptions ?? {};
    this.game = new WumpusGame(gameOptions);
    this.turn = 1;
    this.logNode.innerHTML = '';
    this.updateView([]);
  }

  execute(command) {
    if (this.game.snapshot().status !== WumpusGame.status.IN_PROGRESS) {
      this.pushLog('Game over. Start a new game to continue.');
      return;
    }

    const lines = this.game.executeCommand(command);
    this.pushLog(`> ${command}`);
    if (lines.length === 0) {
      this.pushLog('(no visible effect)');
    }
    lines.forEach((line) => this.pushLog(line));

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
}
