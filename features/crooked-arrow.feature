Feature: Crooked arrow shooting
  The player shoots one to five-room crooked arrows that follow legal tunnels or deviate randomly on invalid segments.

  # Crooked Arrow 001
  Scenario Outline: Crooked Arrow 001 arrow path that reaches the Wumpus wins immediately
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    When the player shoots the path <path>
    Then the player wins
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path       | remaining_arrows | messages                                                                 |
      | 1           | 2           | 13, 14    | 16, 17    | 5      | 2          | 4                | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | 2, 10      | 4                | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |
      | 6           | 18          | 13, 14    | 16, 17    | 5      | 7, 8, 9, 18 | 4               | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Crooked Arrow 002
  Scenario Outline: Crooked Arrow 002 invalid arrow segment deviates to a random adjacent room
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    And the next arrow deviation room is <deviation_room>
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the arrow traversed rooms are <traversed_rooms>
    And the game is <game_status>
    And the player has <remaining_arrows> arrows

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path  | deviation_room | traversed_rooms | wake_choice | game_status | remaining_arrows |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 3     | 5              | 5               | stay        | in progress | 4                |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 3, 4  | 5              | 5, 4            | stay        | in progress | 4                |

  # Crooked Arrow 003
  Scenario Outline: Crooked Arrow 003 arrow can hit the player after a legal or deviated path
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    And the next arrow deviation room is <deviation_room>
    When the player shoots the path <path>
    Then the player loses
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path       | deviation_room | remaining_arrows | messages                                      |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 2, 1       | none           | 4                | OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE!   |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 3          | 1              | 4                | OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE!   |
      | 6           | 20          | 13, 14    | 16, 17    | 5      | 7, 8, 1, 5, 6 | none        | 4                | OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE!   |

  # Crooked Arrow 004
  Scenario Outline: Crooked Arrow 004 missed arrow wakes the Wumpus
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path | wake_choice | expected_wumpus_room | game_status | remaining_arrows | messages                                            |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | 5    | stay        | 10                   | in progress | 4                | MISSED                                              |
      | 1           | 2           | 13, 14    | 16, 17    | 5      | 5    | move to 1   | 1                    | lost        | 4                | MISSED, TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE! |

  # Crooked Arrow 005
  Scenario Outline: Crooked Arrow 005 shooting spends one arrow
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the player has <remaining_arrows> arrows
    And the game is <game_status>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path | wake_choice | remaining_arrows | game_status |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | 5    | stay        | 4                | in progress |
      | 1           | 10          | 13, 14    | 16, 17    | 2      | 5    | stay        | 1                | in progress |

  # Crooked Arrow 006
  Scenario Outline: Crooked Arrow 006 missing with the last arrow loses
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the player loses
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path | wake_choice | remaining_arrows | messages                                       |
      | 1           | 10          | 13, 14    | 16, 17    | 1      | 5    | stay        | 0                | MISSED, YOU RAN OUT OF ARROWS, HA HA HA - YOU LOSE! |

  # Crooked Arrow 007
  Scenario Outline: Crooked Arrow 007 shooting path must contain one to five rooms
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the player has <arrows> arrows
    When the player shoots the path <path>
    Then the shot is rejected with message <message>
    And the player has <arrows> arrows
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path             | message           |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | none             | CAN'T SHOOT THERE |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | 2, 3, 4, 5, 6, 7 | CAN'T SHOOT THERE |
