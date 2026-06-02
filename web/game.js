const CAVE_EXITS = Object.freeze({
  1: [2, 5, 8],
  2: [1, 3, 10],
  3: [2, 4, 12],
  4: [3, 5, 14],
  5: [1, 4, 6],
  6: [5, 7, 15],
  7: [6, 8, 17],
  8: [1, 7, 9],
  9: [8, 10, 18],
  10: [2, 9, 11],
  11: [10, 12, 19],
  12: [3, 11, 13],
  13: [12, 14, 20],
  14: [4, 13, 15],
  15: [6, 14, 16],
  16: [15, 17, 20],
  17: [7, 16, 18],
  18: [9, 17, 19],
  19: [11, 18, 20],
  20: [13, 16, 19],
});

const STATUS = Object.freeze({
  IN_PROGRESS: 'in progress',
  LOST: 'lost',
  WON: 'won',
});

function createSeededRandom(seed) {
  let value = Math.abs(Number(seed)) || 1;
  return () => {
    value = (value * 1664525 + 1013904223) >>> 0;
    return value / 0x100000000;
  };
}

function randomInt(random, max) {
  return Math.floor(random() * max);
}

function hasTunnel(from, to) {
  const exits = CAVE_EXITS[from];
  return Array.isArray(exits) && exits.includes(to);
}

function uniqueRooms(setup) {
  const values = [setup.player, setup.wumpus, ...setup.pits, ...setup.bats];
  if (Number.isInteger(setup.grenadeRoom)) {
    values.push(setup.grenadeRoom);
  }
  return new Set(values).size === values.length;
}

function randomSetup(random) {
  const rooms = Array.from({ length: 20 }, (_, i) => i + 1);
  for (let i = rooms.length - 1; i > 0; i -= 1) {
    const j = randomInt(random, i + 1);
    [rooms[i], rooms[j]] = [rooms[j], rooms[i]];
  }
  return {
    player: rooms[0],
    wumpus: rooms[1],
    pits: [rooms[2], rooms[3]],
    bats: [rooms[4], rooms[5]],
    grenadeRoom: rooms[6],
  };
}

export class WumpusGame {
  constructor(options = {}) {
    this.random = options.rng ?? createSeededRandom(options.seed ?? Date.now());
    const setup = options.setup ?? randomSetup(this.random);
    if (!Array.isArray(setup.pits) || setup.pits.length !== 2) {
      throw new Error('setup must contain exactly two pits');
    }
    if (!Array.isArray(setup.bats) || setup.bats.length !== 2) {
      throw new Error('setup must contain exactly two bat rooms');
    }
    if (!uniqueRooms(setup)) {
      throw new Error('setup rooms must be unique');
    }
    this.state = {
      player: setup.player,
      wumpus: setup.wumpus,
      pits: [...setup.pits],
      bats: [...setup.bats],
      grenadeRoom: Number.isInteger(setup.grenadeRoom) ? setup.grenadeRoom : null,
      carriesGrenade: Boolean(options.carriesGrenade),
      pendingGrenadeRoom: null,
      arrows: Number.isInteger(options.arrows) ? options.arrows : 5,
      status: STATUS.IN_PROGRESS,
    };
  }

  static get caveExits() {
    return CAVE_EXITS;
  }

  static get status() {
    return STATUS;
  }

  snapshot() {
    return {
      ...this.state,
      pits: [...this.state.pits],
      bats: [...this.state.bats],
    };
  }

  turnWarnings() {
    if (this.state.status !== STATUS.IN_PROGRESS) {
      return [];
    }
    const exits = CAVE_EXITS[this.state.player] || [];
    const warnings = [];
    if (exits.includes(this.state.wumpus)) {
      warnings.push('I SMELL A WUMPUS');
    }
    if (exits.some((room) => this.state.bats.includes(room))) {
      warnings.push('BATS NEARBY');
    }
    if (exits.some((room) => this.state.pits.includes(room))) {
      warnings.push('I FEEL A DRAFT');
    }
    return warnings;
  }

  executeCommand(raw) {
    const fields = String(raw).trim().split(/\s+/).filter(Boolean);
    if (fields.length === 0) {
      return [' IS NOT A COMMAND'];
    }

    const action = fields[0].toLowerCase();
    if (this.state.status !== STATUS.IN_PROGRESS) {
      return [];
    }

    let result;
    const shouldDetonate = this.state.pendingGrenadeRoom !== null && (action === 'm' || action === 's' || action === 'r' || action === 'rest');

    if (action === 'm') {
      const room = Number.parseInt(fields[1], 10);
      result = this.move(room);
    } else if (action === 's') {
      const path = fields.slice(1).map((value) => Number.parseInt(value, 10));
      result = this.shoot(path);
    } else if (action === 't') {
      const room = Number.parseInt(fields[1], 10);
      result = this.throwGrenade(room);
    } else if (action === 'r' || action === 'rest') {
      if (fields.length !== 1) {
        return [`${fields[0].toUpperCase()} IS NOT A COMMAND`];
      }
      result = { rejectedMessage: '', messages: [] };
    } else {
      return [`${fields[0].toUpperCase()} IS NOT A COMMAND`];
    }

    if (result.rejectedMessage) {
      return [result.rejectedMessage];
    }

    const lines = [...result.messages];
    if (shouldDetonate) {
      lines.push(...this.detonateGrenade());
    }
    return lines;
  }

  move(to) {
    if (!Number.isInteger(to) || !hasTunnel(this.state.player, to)) {
      return { rejectedMessage: "CAN'T MOVE THERE", messages: [] };
    }
    this.state.player = to;
    return { rejectedMessage: '', messages: this.resolveArrival() };
  }

  shoot(path) {
    if (!Array.isArray(path) || path.length < 1 || path.length > 5 || path.some((room) => !Number.isInteger(room) || room < 1 || room > 20)) {
      return { rejectedMessage: "CAN'T SHOOT THERE", messages: [], traversedRooms: [] };
    }

    this.state.arrows -= 1;
    let arrowRoom = this.state.player;
    const traversed = [];

    for (const requested of path) {
      arrowRoom = hasTunnel(arrowRoom, requested) ? requested : this.randomExit(arrowRoom);
      traversed.push(arrowRoom);
      if (arrowRoom === this.state.wumpus) {
        this.state.status = STATUS.WON;
        return {
          rejectedMessage: '',
          traversedRooms: traversed,
          messages: ["AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"],
        };
      }
      if (arrowRoom === this.state.player) {
        return { rejectedMessage: '', traversedRooms: traversed, messages: this.lose('OUCH! ARROW GOT YOU!') };
      }
    }

    const messages = ['MISSED', ...this.wakeWumpus()];
    if (this.state.status === STATUS.IN_PROGRESS && this.state.arrows === 0) {
      messages.push(...this.lose('YOU RAN OUT OF ARROWS'));
    }

    return { rejectedMessage: '', traversedRooms: traversed, messages };
  }

  throwGrenade(target) {
    if (!this.state.carriesGrenade || !Number.isInteger(target) || !hasTunnel(this.state.player, target)) {
      return { rejectedMessage: "CAN'T THROW THERE", messages: [] };
    }
    this.state.carriesGrenade = false;
    this.state.pendingGrenadeRoom = target;
    return { rejectedMessage: '', messages: ['YOU HEAR TIC...TIC...'] };
  }

  detonateGrenade() {
    const target = this.state.pendingGrenadeRoom;
    if (target === null) {
      return [];
    }
    this.state.pendingGrenadeRoom = null;
    const blastRooms = new Set([target, ...CAVE_EXITS[target]]);
    this.state.bats = this.state.bats.filter((room) => !blastRooms.has(room));

    const messages = ['YOU HEAR A HORRENDOUS EXPLOSION!'];
    if (blastRooms.has(this.state.wumpus)) {
      this.state.status = STATUS.WON;
      messages.push("AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!");
      return messages;
    }
    if (blastRooms.has(this.state.player)) {
      messages.push(...this.lose('YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!'));
      return messages;
    }
    return messages;
  }

  resolveArrival() {
    if (this.state.pits.includes(this.state.player)) {
      return this.lose('YYYIIIIEEEE . . . FELL IN PIT');
    }
    if (this.state.bats.includes(this.state.player)) {
      const messages = ['ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!'];
      this.state.player = this.randomRoom();
      return [...messages, ...this.resolveArrival()];
    }
    if (this.state.player === this.state.wumpus) {
      return this.wakeWumpus();
    }
    if (this.state.grenadeRoom === this.state.player) {
      this.state.grenadeRoom = null;
      this.state.carriesGrenade = true;
      return ['YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!'];
    }
    return [];
  }

  wakeWumpus() {
    if (this.state.status !== STATUS.IN_PROGRESS) {
      return [];
    }
    if (randomInt(this.random, 4) !== 0) {
      this.state.wumpus = this.randomExit(this.state.wumpus);
    }
    if (this.state.wumpus === this.state.player) {
      return this.lose('TSK TSK TSK - WUMPUS GOT YOU!');
    }
    return [];
  }

  randomRoom() {
    return randomInt(this.random, 20) + 1;
  }

  randomExit(room) {
    const exits = CAVE_EXITS[room];
    return exits[randomInt(this.random, exits.length)];
  }

  lose(reason) {
    this.state.status = STATUS.LOST;
    return [reason, 'HA HA HA - YOU LOSE!'];
  }
}

export { createSeededRandom };
