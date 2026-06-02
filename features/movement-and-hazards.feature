Feature: Movement and hazard resolution
  Player movement validates cave tunnels, relocates the player, and resolves arrival hazards immediately.

  # Movement And Hazards 001
  Scenario Outline: Movement And Hazards 001 legal move enters an empty adjacent room
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the player moves to room <to_room>
    Then the player is in room <to_room>
    And the game is still in progress
    And the turn messages are <messages>

    Examples:
      | from_room | to_room | wumpus_room | pit_rooms | bat_rooms | messages |
      | 1         | 2       | 20          | 13, 14    | 16, 17    | none     |
      | 10        | 11      | 20          | 13, 14    | 16, 17    | none     |

  # Movement And Hazards 002
  Scenario Outline: Movement And Hazards 002 illegal move is rejected without relocating
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the player moves to room <to_room>
    Then the move is rejected with message <message>
    And the player is in room <from_room>
    And the game is still in progress

    Examples:
      | from_room | to_room | wumpus_room | pit_rooms | bat_rooms | message          |
      | 1         | 20      | 19          | 13, 14    | 16, 17    | CAN'T MOVE THERE |
      | 10        | 14      | 19          | 13, 15    | 16, 17    | CAN'T MOVE THERE |

  # Movement And Hazards 003
  Scenario Outline: Movement And Hazards 003 moving into a pit loses immediately
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the player moves to room <pit_room>
    Then the player loses
    And the turn messages are <messages>

    Examples:
      | from_room | pit_room | wumpus_room | pit_rooms | bat_rooms | messages                                                |
      | 1         | 2        | 20          | 2, 14     | 16, 17    | YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!     |

  # Movement And Hazards 004
  Scenario Outline: Movement And Hazards 004 moving into bats relocates the player
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next bat relocation room is <relocation_room>
    When the player moves to room <bat_room>
    Then the player is in room <relocation_room>
    And the game is <game_status>
    And the turn messages are <messages>

    Examples:
      | from_room | bat_room | relocation_room | wumpus_room | pit_rooms | bat_rooms | game_status | messages                                                                 |
      | 1         | 2        | 3               | 20          | 13, 14    | 2, 17     | in progress | ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!                         |
      | 1         | 2        | 13              | 20          | 13, 14    | 2, 17     | lost        | ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE! |

  # Movement And Hazards 005
  Scenario Outline: Movement And Hazards 005 bat relocation into the Wumpus room resolves Wumpus wake
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next bat relocation room is <wumpus_room>
    And the next Wumpus wake choice is <wake_choice>
    When the player moves to room <bat_room>
    Then the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>
    And the turn messages are <messages>

    Examples:
      | from_room | bat_room | wumpus_room | wake_choice | expected_wumpus_room | pit_rooms | bat_rooms | game_status | messages                                                                                     |
      | 1         | 2        | 10          | stay        | 10                   | 13, 14    | 2, 17     | lost        | ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE! |
      | 1         | 2        | 10          | move to 11  | 11                   | 13, 14    | 2, 17     | in progress | ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!                                             |

  # Movement And Hazards 006
  Scenario Outline: Movement And Hazards 006 moving into the Wumpus room wakes the Wumpus
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next Wumpus wake choice is <wake_choice>
    When the player moves to room <wumpus_room>
    Then the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>
    And the turn messages are <messages>

    Examples:
      | from_room | wumpus_room | wake_choice | expected_wumpus_room | pit_rooms | bat_rooms | game_status | messages                                                       |
      | 1         | 2           | stay        | 2                    | 13, 14    | 16, 17    | lost        | TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE!            |
      | 1         | 2           | move to 3   | 3                    | 13, 14    | 16, 17    | in progress | none                                                           |

  # Movement And Hazards 007
  Scenario Outline: Movement And Hazards 007 Wumpus wake movement ignores pits and bats
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next Wumpus wake choice is <wake_choice>
    When the player moves to room <wumpus_room>
    Then the Wumpus is in room <expected_wumpus_room>
    And the game is in progress

    Examples:
      | from_room | wumpus_room | wake_choice | expected_wumpus_room | pit_rooms | bat_rooms |
      | 1         | 2           | move to 10  | 10                   | 10, 14    | 3, 17     |
