const icon = (svg) => `<svg viewBox="0 0 64 64" aria-hidden="true">${svg}</svg>`;

export const ROOM_ARTWORK = Object.freeze({
  safe: {
    label: 'Safe Room',
    art: icon('<circle cx="32" cy="32" r="24" fill="#2f855a"/><path d="M20 34h24" stroke="#c6f6d5" stroke-width="4"/><circle cx="32" cy="24" r="4" fill="#c6f6d5"/>'),
  },
  player: {
    label: 'Player Room',
    art: icon('<circle cx="32" cy="20" r="10" fill="#f6e05e"/><path d="M16 50c0-8 7-14 16-14s16 6 16 14" fill="#ecc94b"/><circle cx="28" cy="20" r="2"/><circle cx="36" cy="20" r="2"/>'),
  },
  wumpus: {
    label: 'Wumpus Lair',
    art: icon('<circle cx="22" cy="30" r="12" fill="#9f1239"/><circle cx="42" cy="30" r="12" fill="#9f1239"/><circle cx="32" cy="38" r="14" fill="#be123c"/><circle cx="27" cy="35" r="2" fill="#fff"/><circle cx="37" cy="35" r="2" fill="#fff"/><path d="M26 43h12" stroke="#fff" stroke-width="3"/>'),
  },
  pit: {
    label: 'Pit Room',
    art: icon('<ellipse cx="32" cy="36" rx="22" ry="14" fill="#1a202c"/><ellipse cx="32" cy="36" rx="12" ry="7" fill="#000"/><path d="M16 22l6 8m26-8-6 8" stroke="#718096" stroke-width="3"/>'),
  },
  bats: {
    label: 'Bat Room',
    art: icon('<path d="M8 34c8-12 16-12 24 0 8-12 16-12 24 0-8 2-14 6-18 12-4-6-10-10-18-12z" fill="#6b46c1"/><circle cx="32" cy="30" r="4" fill="#f7fafc"/>'),
  },
  grenade: {
    label: 'Grenade Room',
    art: icon('<circle cx="30" cy="36" r="14" fill="#4a5568"/><rect x="30" y="18" width="14" height="8" rx="2" fill="#718096"/><path d="M38 14c4-4 8-4 12 0" stroke="#a0aec0" stroke-width="3" fill="none"/>'),
  },
});
