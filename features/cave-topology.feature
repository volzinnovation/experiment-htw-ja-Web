Feature: Cave topology
  The cave is the fixed twenty-room dodecahedron used by Hunt the Wumpus.

  # Cave Topology 001
  Scenario Outline: Cave Topology 001 canonical room exits
    Given a new cave
    When the exits for room <room> are queried
    Then the exits are <exits>

    Examples:
      | room | exits      |
      | 1    | 2, 5, 8    |
      | 2    | 1, 3, 10   |
      | 3    | 2, 4, 12   |
      | 4    | 3, 5, 14   |
      | 5    | 1, 4, 6    |
      | 6    | 5, 7, 15   |
      | 7    | 6, 8, 17   |
      | 8    | 1, 7, 9    |
      | 9    | 8, 10, 18  |
      | 10   | 2, 9, 11   |
      | 11   | 10, 12, 19 |
      | 12   | 3, 11, 13  |
      | 13   | 12, 14, 20 |
      | 14   | 4, 13, 15  |
      | 15   | 6, 14, 16  |
      | 16   | 15, 17, 20 |
      | 17   | 7, 16, 18  |
      | 18   | 9, 17, 19  |
      | 19   | 11, 18, 20 |
      | 20   | 13, 16, 19 |

  # Cave Topology 002
  Scenario Outline: Cave Topology 002 valid room has exactly three exits
    Given a new cave
    When the exits for room <room> are queried
    Then the exit count is <exit_count>

    Examples:
      | room | exit_count |
      | 1    | 3          |
      | 10   | 3          |
      | 20   | 3          |

  # Cave Topology 003
  Scenario: Cave Topology 003 cave is connected
    Given a new cave
    When the cave is traversed from room 1
    Then every room from 1 through 20 is reachable

  # Cave Topology 004
  Scenario Outline: Cave Topology 004 tunnel links are bidirectional
    Given a new cave
    When the tunnel from room <from_room> to room <to_room> is queried
    Then the reverse tunnel also exists

    Examples:
      | from_room | to_room |
      | 1         | 2       |
      | 1         | 5       |
      | 1         | 8       |
      | 10        | 11      |
      | 16        | 20      |

  # Cave Topology 005
  Scenario Outline: Cave Topology 005 rooms do not connect to themselves
    Given a new cave
    When the exits for room <room> are queried
    Then room <room> is not one of the exits

    Examples:
      | room |
      | 1    |
      | 7    |
      | 14   |
      | 20   |

  # Cave Topology 006
  Scenario Outline: Cave Topology 006 adjacent hazard query reports only neighboring hazards
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When adjacent hazards are queried from the player room
    Then the adjacent hazard types are <hazards>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | hazards          |
      | 1           | 2           | 3, 4      | 5, 6      | Wumpus, Bats     |
      | 10          | 2           | 9, 18     | 6, 7      | Wumpus, Pit      |
      | 13          | 7           | 12, 14    | 20, 1     | Pit, Bats        |
      | 6           | 20          | 1, 2      | 3, 4      | none             |
