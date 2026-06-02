# Architecture Improvement Plan

## Current Shape

The project has a useful three-layer structure:

- `cmd/htwgo`: terminal adapter, flags, QA commands, and transcript rendering.
- `internal/interactive`: command/session orchestration.
- `internal/wumpus`: game rules and state transitions.

The acceptance runner is intentionally outside the game layers and drives the domain through feature steps.

## Weaknesses To Address

1. Random setup and random in-game events share one RNG stream.
   - A new setup draw or event draw can shift every later seeded behavior.
   - Plan: split setup placement RNG from event RNG and expose only focused random-event helpers inside `Game`.

2. Domain mutation has unsafe setup setters.
   - QA/test controls can bypass placement invariants.
   - Plan: name these APIs as QA controls, keep explicit room validation at the CLI adapter, and make normal setup construction remain validated through `NewGameWithSetup`.

3. Turn sequencing is spread across command handlers.
   - Jumping Wumpus, command execution, grenade detonation, turn count, and loss prompt rules are coordinated manually.
   - Plan: introduce a small turn executor in `interactive.Session` so command handlers provide only the command action.

4. Replay naming is too narrow.
   - `ReplaySameSetup` restores a pre-loss snapshot, not only setup placement.
   - Plan: introduce `ReplaySnapshot` and keep `ReplaySameSetup` as a compatibility wrapper while session code uses the clearer name.

5. Acceptance handler registration is too centralized.
   - `NewHandlers` is a large mixed registry.
   - Plan: split registration into feature-oriented helper functions while preserving the same handler map.

## Execution Checklist

1. [x] Add an event RNG seam in `internal/wumpus`.
2. [x] Update random event consumers to use the event RNG seam.
3. [x] Rename CLI QA mutation calls to explicit QA setup controls.
4. [x] Refactor session command execution through one turn pipeline.
5. [x] Add `ReplaySnapshot` and move session replay to it.
6. [x] Split acceptance handler registration helpers.
7. [x] Run `go test ./...` and `./scripts/acceptance.sh`.
