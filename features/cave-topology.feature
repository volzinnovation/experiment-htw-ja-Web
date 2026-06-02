# mutation-stamp: sha256=281c6087a54cdc8a387bb2b49e39f64582d4faf960f1b58f8ed6c1c57e9792f6
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T14:14:49Z","feature_name":"Cave topology","feature_path":"features/cave-topology.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:1a9ae523adb25befd2d869df4a2abf1ff2f613dfdebf887e8fa9a1f59c54490a","scenarios":[{"index":0,"name":"Cave Topology 001 canonical room exits","scenario_hash":"69f077d60b370d45ea016344c3f1f507025b91d8518e18ae47ed91db7d3fe2c1","mutation_count":40,"result":{"Total":40,"Killed":40,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:03:01Z"},{"index":3,"name":"Cave Topology 004 tunnel links are bidirectional","scenario_hash":"b016c11eff193b2b036daf12ced4386810fb116605c7346ef6fe59150f893818","mutation_count":10,"result":{"Total":10,"Killed":10,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:03:01Z"}]}
# acceptance-mutation-manifest-end

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
  Scenario: Cave Topology 002 every valid room has exactly three exits
    Given a new cave
    When the cave invariants are inspected
    Then every room has exactly three exits

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
  Scenario: Cave Topology 005 rooms do not connect to themselves
    Given a new cave
    When the cave invariants are inspected
    Then no room is one of its own exits

  # Cave Topology 006
  Scenario: Cave Topology 006 adjacent hazard query reports Wumpus and bats
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 3 and 4, and bats in rooms 5 and 6
    When adjacent hazards are queried from the player room
    Then the adjacent hazard types are Wumpus and Bats

  # Cave Topology 007
  Scenario: Cave Topology 007 adjacent hazard query reports Wumpus and pit
    Given a game setup with the player in room 10, the Wumpus in room 2, pits in rooms 9 and 18, and bats in rooms 6 and 7
    When adjacent hazards are queried from the player room
    Then the adjacent hazard types are Wumpus and Pit

  # Cave Topology 008
  Scenario: Cave Topology 008 adjacent hazard query reports no neighboring hazards
    Given a game setup with the player in room 6, the Wumpus in room 20, pits in rooms 1 and 2, and bats in rooms 3 and 4
    When adjacent hazards are queried from the player room
    Then there are no adjacent hazard types
