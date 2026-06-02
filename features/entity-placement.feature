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
  Scenario Outline: Entity Placement 002 occupied rooms are valid cave rooms
    Given a new game created with seed <seed>
    When the occupied rooms are inspected
    Then every occupied room number is from 1 through 20

    Examples:
      | seed |
      | 1    |
      | 1973 |
      | 2026 |

  # Entity Placement 003
  Scenario Outline: Entity Placement 003 occupants are placed in distinct rooms
    Given a new game created with seed <seed>
    When the occupied rooms are inspected
    Then exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats

    Examples:
      | seed | occupied_count |
      | 1    | 6              |
      | 1973 | 6              |
      | 2026 | 6              |

  # Entity Placement 004
  Scenario Outline: Entity Placement 004 same seed creates the same setup
    Given a new game created with seed <seed>
    And another new game created with seed <seed>
    When both setups are inspected
    Then both setups have identical player, Wumpus, pit, and bat rooms

    Examples:
      | seed |
      | 1    |
      | 1973 |
      | 2026 |

  # Entity Placement 005
  Scenario Outline: Entity Placement 005 same setup replay preserves placement
    Given a completed game created with seed <seed>
    When a same setup replay is started
    Then the replay setup has identical player, Wumpus, pit, and bat rooms

    Examples:
      | seed |
      | 1973 |
