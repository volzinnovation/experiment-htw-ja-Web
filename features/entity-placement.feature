# mutation-stamp: sha256=24149fde36ab104faa185032d96dfacf124c9657c26f5aa861727e491ce28feb
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T14:14:49Z","feature_name":"Entity placement","feature_path":"features/entity-placement.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:f8324f871f5a0fedf4a32f905b77a7767e6e4c3622a7ca4c4cd95eab15ddbc84","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Entity placement
  New Hunt the Wumpus games place the player, Wumpus, pits, and bats in valid distinct rooms.

  # Entity Placement 001
  Scenario: Entity Placement 001 random setup creates the required occupants
    Given a new game created with seed 1973
    When the setup is inspected
    Then there is 1 player
    And there is 1 Wumpus
    And there are 2 pits
    And there are 2 bats

  # Entity Placement 002
  Scenario: Entity Placement 002 occupied rooms are valid cave rooms
    Given a new game created with seed 1973
    When the occupied rooms are inspected
    Then every occupied room number is from 1 through 20

  # Entity Placement 003
  Scenario: Entity Placement 003 occupants are placed in distinct rooms
    Given a new game created with seed 1973
    When the occupied rooms are inspected
    Then exactly 6 distinct rooms are occupied by the player, Wumpus, pits, and bats

  # Entity Placement 004
  Scenario: Entity Placement 004 same seed creates the same setup
    Given a new game created with seed 1973
    And another new game created with seed 1973
    When both setups are inspected
    Then both setups have identical player, Wumpus, pit, and bat rooms

  # Entity Placement 005
  Scenario: Entity Placement 005 same setup replay preserves placement
    Given a completed game created with seed 1973
    When a same setup replay is started
    Then the replay setup has identical player, Wumpus, pit, and bat rooms
