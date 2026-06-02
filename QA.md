# Console-Only QA Plan

This plan validates Hunt the Wumpus only through the playable terminal interface. Do not use acceptance tests, unit tests, property tests, mutation tests, generated runners, or source-level inspection as QA verification.

## 1. Console Launch Modes

Start normal play:

```sh
go run ./cmd/htwgo
```

Start QA play with inert hazards:

```sh
go run ./cmd/htwgo --qa-inert-hazards
```

Start QA play with all hidden state printed each turn:

```sh
go run ./cmd/htwgo --qa-reveal-state
```

Start deterministic QA play:

```sh
go run ./cmd/htwgo --qa-seed 1973
```

Flags may be combined:

```sh
go run ./cmd/htwgo --qa-inert-hazards --qa-reveal-state --qa-seed 1973
```

Expected QA banner when all flags are active:

```text
QA MODE ENABLED: HAZARDS INERT, STATE REVEALED, SEEDED SETUP
```

The reveal flag prints a line like:

```text
QA STATE: PLAYER=3 WUMPUS=13 WUMPUS_STATE=awake PITS=15,5 BATS=1,18 HHG=12 CARRYING_HHG=false PENDING_HHG=none FUSE=none ARROWS=5
```

The inert-hazards flag keeps hazards placed and warnings visible, but suppresses hazard consequences. Expected inert messages include:

```text
QA INERT: PIT IGNORED
QA INERT: BATS IGNORED
QA INERT: WUMPUS IGNORED
QA INERT: BLAST IGNORED
```

QA mode also enables setup commands typed at the game prompt:

```text
qa set player <room>
qa set wumpus <room>
qa set pits <room1> <room2>
qa set bats <room1> <room2>
qa set hhg <room>
qa set hhg none
qa set arrows <count>
```

Expected output examples:

```text
QA SET: WUMPUS=10
QA SET: PITS=10,12
QA SET: BATS=14,18
QA SET: HHG=12
```

These commands are available only when a QA launch flag is active. In normal play, `qa set wumpus 10` must be rejected as an invalid command.

## 2. Basic Console Contract

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. Answer the instruction prompt with `n`.
3. Verify the first turn contains:
   - `QA STATE: ...`
   - `YOU ARE IN ROOM <n>`
   - `TUNNELS LEAD TO <a> <b> <c>`
   - zero or more hazard warnings
   - `ARROWS LEFT: 5`
   - `SHOOT OR MOVE (S-M)?`
4. Type:

   ```text
   x
   ```

   Expected output: `X IS NOT A COMMAND`. The next `QA STATE` line must show unchanged player room and arrows.

5. Type a legal move using one displayed tunnel:

   ```text
   m <displayed tunnel>
   ```

   Expected output: the next room display and `QA STATE` show the destination room.

## 3. Cave Topology

Feature file: `features/cave-topology.feature`

Use inert hazards so the map can be traversed without random deaths.

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-inert-hazards --qa-reveal-state --qa-seed 1973
   ```

2. Visit all 20 rooms. Use the displayed tunnel list to choose moves.
3. For every visited room, record:
   - room number from `YOU ARE IN ROOM`
   - exact `TUNNELS LEAD TO` list
   - matching `PLAYER=<room>` from `QA STATE`
4. Verify every room lists exactly three tunnels.
5. Verify bidirectionality by moving from `A` to `B` and confirming `A` appears in `B`'s tunnel list.
6. Verify the canonical exits:
   - `1 -> 2 5 8`
   - `2 -> 1 3 10`
   - `3 -> 2 4 12`
   - `4 -> 3 5 14`
   - `5 -> 1 4 6`
   - `6 -> 5 7 15`
   - `7 -> 6 8 17`
   - `8 -> 1 7 9`
   - `9 -> 8 10 18`
   - `10 -> 2 9 11`
   - `11 -> 10 12 19`
   - `12 -> 3 11 13`
   - `13 -> 12 14 20`
   - `14 -> 4 13 15`
   - `15 -> 6 14 16`
   - `16 -> 15 17 20`
   - `17 -> 7 16 18`
   - `18 -> 9 17 19`
   - `19 -> 11 18 20`
   - `20 -> 13 16 19`

## 4. Entity Placement

Feature file: `features/entity-placement.feature`

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. Verify `QA STATE` shows one player, one Wumpus, two pits, two bats, one HHG, and 5 arrows.
3. Verify all rooms are 1 through 20.
4. Verify all occupied rooms are distinct.
5. Repeat with:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1975
   ```

6. Verify a different but valid placement is shown.

## 5. Setup Reproducibility And Replay

Feature files:

- `features/entity-placement.feature`
- `features/interactive-game-loop.feature`

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. Record the first `QA STATE` line.
3. Exit and restart the same command.
4. Verify the first `QA STATE` line is identical.
5. Start with seed `1975` and verify the first `QA STATE` line differs from seed `1973` but is valid.
6. In QA play, set a pit in an adjacent room and move into it:

   ```text
   qa set pits <displayed tunnel> <another room>
   m <displayed tunnel>
   ```

7. Verify the loss transcript includes:

   ```text
   HA HA HA - YOU LOSE!
   SAME SET UP (Y-N)?
   ```

8. Type `y` and verify the replay starts with the same setup.
9. Repeat, type `n`, and verify the new game does not preserve the prior setup.

## 6. Hazard Warnings

Feature file: `features/turn-warnings.feature`

Start with:

```sh
go run ./cmd/htwgo --qa-inert-hazards --qa-reveal-state --qa-seed 1973
```

Use `qa set` commands to place the player in room 1 and isolate each warning:

```text
qa set player 1
qa set wumpus 2
qa set pits 11 12
qa set bats 13 14
```

Expected warning: `I SMELL A WUMPUS`.

```text
qa set player 1
qa set wumpus 13
qa set pits 11 12
qa set bats 5 14
```

Expected warning: `BATS NEARBY`.

```text
qa set player 1
qa set wumpus 13
qa set pits 8 12
qa set bats 14 15
```

Expected warning: `I FEEL A DRAFT`.

Verify non-adjacent hazards do not warn. This setup reproduces the class of defect where a hazard two tunnels away prints a warning: room 3 has exits 2, 4, and 12, while the Wumpus, pits, and bats below are all non-adjacent.

```text
qa set player 3
qa set wumpus 13
qa set pits 15 5
qa set bats 1 18
```

Expected warnings: none. The transcript must not include:

```text
I SMELL A WUMPUS
BATS NEARBY
I FEEL A DRAFT
```

The turn display must go straight from `QA STATE: ...` to:

```text
YOU ARE IN ROOM 3
TUNNELS LEAD TO 2 4 12
ARROWS LEFT: 5
SHOOT OR MOVE (S-M)?
```

With all hazards adjacent, verify warning order:

```text
I SMELL A WUMPUS
BATS NEARBY
I FEEL A DRAFT
```

## 7. Movement And Hazards

Feature file: `features/movement-and-hazards.feature`

1. Start with reveal state:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. Legal move: `m <displayed tunnel>`. Expected result: `PLAYER=<destination>` and `YOU ARE IN ROOM <destination>`.
3. Illegal move: `m <non-neighbor>`. Expected output: `CAN'T MOVE THERE`; player room and arrows unchanged.
4. Pit entry:

   ```text
   qa set player 1
   qa set pits 2 12
   m 2
   ```

   Expected pit loss and same-setup prompt.

5. Bat entry:

   ```text
   qa set player 1
   qa set bats 2 17
   m 2
   ```

   Expected: `ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!`

6. Wumpus entry:

   ```text
   qa set player 1
   qa set wumpus 2
   m 2
   ```

   Expected Wumpus loss unless inert hazards are enabled.

7. Repeat hazard entries with `--qa-inert-hazards` and verify the matching `QA INERT` message appears and play continues.

## 8. Crooked Arrow Shooting

Feature file: `features/crooked-arrow.feature`

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. One-room shot: `s <displayed tunnel>`. Expected result: arrows decrease by one unless the shot wins or loses.
3. Multi-room shot: `s <room1> <room2> <room3>`. Expected result: path length 1 to 5 is accepted.
4. Empty shot: `s`. Expected output: `CAN'T SHOOT THERE`; arrows unchanged.
5. Six-room shot: `s 1 2 3 4 5 6`. Expected output: `CAN'T SHOOT THERE`; arrows unchanged.
6. Set the Wumpus on a known arrow path and shoot:

   ```text
   qa set player 1
   qa set wumpus 12
   s 2 3 12
   ```

   Expected win text:

   ```text
   AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!
   ```

7. Fire a path that returns through the player room:

   ```text
   qa set player 1
   qa set wumpus 12
   s 2 1
   ```

   Expected self-hit loss.

8. Miss with five arrows. Use `qa set wumpus <room not in the shot path>` and repeat legal misses until arrows reach zero. Expected result: fifth miss loses.

## 9. Command Validation

Feature file: `features/interactive-game-loop.feature`

Run these in a reveal-state session and verify invalid commands do not change `PLAYER`, `ARROWS`, carried grenade state, pending grenade, or outcome.

- Unknown command: `x`, expected `X IS NOT A COMMAND`
- Invalid move: `m <non-neighbor>`, expected `CAN'T MOVE THERE`
- Invalid shot: `s`, expected `CAN'T SHOOT THERE`
- Shot too long: `s 1 2 3 4 5 6`, expected `CAN'T SHOOT THERE`
- Throw without grenade: `t <displayed tunnel>`, expected `CAN'T THROW THERE`
- Case-insensitive commands: `M <room>`, `S <room>`, `R`, `ReSt`
- Normal play without QA flags: `qa set wumpus 10` is rejected as invalid

## 10. Holy Hand Grenade

Feature file: `features/holy-hand-grenade.feature`

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-inert-hazards --qa-reveal-state --qa-seed 1973
   ```

2. Put the HHG in an adjacent room and move there:

   ```text
   qa set player 1
   qa set hhg 2
   m 2
   ```

3. Expected acquisition text:

   ```text
   YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!
   ```

4. Expected next prompt: `SHOOT, MOVE OR THROW (S-M-T)?`
5. Invalid non-adjacent throw: `t <non-neighbor>`, expected `CAN'T THROW THERE`.
6. Legal throw: `t <displayed tunnel>`, expected `YOU HEAR TIC...TIC...`.
7. Verify `QA STATE` shows `CARRYING_HHG=false`, `PENDING_HHG=<target>`, and `FUSE=1`.
8. Take one legal next turn. Expected: `YOU HEAR A HORRENDOUS EXPLOSION!`
9. Arrange a Wumpus blast win and a player blast loss with `qa set` commands, then verify the exact win/loss transcripts.

## 11. Rest Command

Feature file: `features/rest-command.feature`

1. Start:

   ```sh
   go run ./cmd/htwgo --qa-reveal-state --qa-seed 1973
   ```

2. Record `PLAYER` and `ARROWS`.
3. Type `r`. Expected: `PLAYER` and `ARROWS` unchanged.
4. Type `rest`. Expected: same as short rest.
5. Put hazards adjacent to the player, then rest. Expected: warnings may display, but rest does not trigger pit, bat, or Wumpus room-entry effects by itself.
6. Put the HHG adjacent, acquire it, throw it, then rest. Expected: the grenade detonates after the rest turn.
7. Invalid rest syntax:

   ```text
   r 5
   rest 12
   ```

   Expected:

   ```text
   R IS NOT A COMMAND
   REST IS NOT A COMMAND
   ```

8. Case-insensitive rest: `R` and `ReSt` both rest successfully.

## 12. Final Console Regression

1. Start normal play:

   ```sh
   go run ./cmd/htwgo
   ```

2. Verify the console prints `INSTRUCTIONS (Y-N)?`.
3. Answer `y` and verify the instruction output includes:
   - `WELCOME TO 'HUNT THE WUMPUS'`
   - `THE WUMPUS LIVES IN A CAVE OF 20 ROOMS`
   - `HAZARDS:`
   - `WUMPUS:`
   - `YOU:`
   - `WARNINGS:`
4. Start normal play again, answer `n`, and verify the game proceeds directly to the turn display without instruction text.
5. Verify room display and prompt.
6. Verify invalid command rejection.
7. Verify legal movement.
8. Verify legal shooting.
9. Verify rest.
10. If the grenade is acquired, verify legal throwing.
11. Finish with either a win or loss.
12. Verify the final transcript contains the correct win or loss message.

## 13. Signoff Criteria

The build is ready for QA signoff only when:

1. Every feature listed here has been exercised through the terminal.
2. QA flags are used only for controllability and are clearly identified by the `QA MODE ENABLED` banner.
3. Every exact message observed in the terminal matches case, punctuation, and spacing.
4. Invalid commands do not advance room, arrows, grenade state, or game outcome.
5. Replay behavior is verified with both `y` and `n`.
6. Normal play without QA flags still has no QA banner, no `QA STATE` line, no inert-hazard messages, and no active `qa set ...` commands.
