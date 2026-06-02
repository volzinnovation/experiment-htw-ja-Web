# mutation-stamp: sha256=f854056a1a1eafc04c987fccfaf7a2a89abe91a408dc0d3085f6c98b199d2b3e
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T14:14:50Z","feature_name":"Movement and hazard resolution","feature_path":"features/movement-and-hazards.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:f7b4a1e34ca682886d47a42f7705487d63c9191a01fef6729a523c1beb19850c","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Movement and hazard resolution
  Player movement validates cave tunnels, relocates the player, and resolves arrival hazards immediately.

  # Movement And Hazards 001
  Scenario: Movement And Hazards 001 legal move enters an empty adjacent room
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, and bats in rooms 16, 17
    When the player moves to room 2
    Then the player is in room 2
    And the game is still in progress
    And the turn messages are none

  # Movement And Hazards 002
  Scenario: Movement And Hazards 002 illegal move is rejected without relocating
    Given a game setup with the player in room 1, the Wumpus in room 19, pits in rooms 13, 14, and bats in rooms 16, 17
    When the player moves to room 20
    Then the move is rejected with message CAN'T MOVE THERE
    And the player is in room 1
    And the game is still in progress

  # Movement And Hazards 003
  Scenario: Movement And Hazards 003 moving into a pit loses immediately
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 2, 14, and bats in rooms 16, 17
    When the player moves to room 2
    Then the player loses
    And the turn messages are YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!

  # Movement And Hazards 004
  Scenario: Movement And Hazards 004 moving into bats relocates the player
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, and bats in rooms 2, 17
    And the next bat relocation room is 3
    When the player moves to room 2
    Then the player is in room 3
    And the game is in progress
    And the turn messages are ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!

  Scenario: Movement And Hazards 004 bat relocation can resolve destination hazards
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, and bats in rooms 2, 17
    And the next bat relocation room is 13
    When the player moves to room 2
    Then the player is in room 13
    And the game is lost
    And the turn messages are ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!

  # Movement And Hazards 005
  Scenario: Movement And Hazards 005 bat relocation into the Wumpus room can lose
    Given a game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, and bats in rooms 2, 17
    And the next bat relocation room is 10
    And the next Wumpus wake choice is stay
    When the player moves to room 2
    Then the Wumpus is in room 10
    And the game is lost
    And the turn messages are ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE!

  Scenario: Movement And Hazards 005 bat relocation into the Wumpus room can continue
    Given a game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, and bats in rooms 2, 17
    And the next bat relocation room is 10
    And the next Wumpus wake choice is move to 11
    When the player moves to room 2
    Then the Wumpus is in room 11
    And the game is in progress
    And the turn messages are ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!

  # Movement And Hazards 006
  Scenario: Movement And Hazards 006 moving into the Wumpus room can lose
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the next Wumpus wake choice is stay
    When the player moves to room 2
    Then the Wumpus is in room 2
    And the game is lost
    And the turn messages are TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE!

  Scenario: Movement And Hazards 006 moving into the Wumpus room can continue
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the next Wumpus wake choice is move to 3
    When the player moves to room 2
    Then the Wumpus is in room 3
    And the game is in progress
    And the turn messages are none

  # Movement And Hazards 007
  Scenario: Movement And Hazards 007 Wumpus wake movement ignores pits and bats
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 10, 14, and bats in rooms 3, 17
    And the next Wumpus wake choice is move to 10
    When the player moves to room 2
    Then the Wumpus is in room 10
    And the game is in progress
