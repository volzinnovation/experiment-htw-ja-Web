Feature: Jumping Wumpus
  On a turn, the Wumpus can make two legal jumps with exact messages and landing outcomes.

  # Jumping Wumpus 001
  Scenario Outline: Jumping Wumpus 001 seeded turn can trigger double jump event
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is <jump_event>
    And the next Wumpus jump path is <jump_path>
    When the next turn begins
    Then the Wumpus is in room <expected_wumpus_room>
    And the displayed lines include <message>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | jump_event | jump_path | expected_wumpus_room | message                 |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | jumps      | 11, 12    | 12                   | YOU HEAR WHUMP, WHUMP. |

  # Jumping Wumpus 002
  Scenario Outline: Jumping Wumpus 002 no jump event leaves Wumpus in place
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is <jump_event>
    When the next turn begins
    Then the Wumpus is in room <wumpus_room>
    And the displayed lines do not include <message>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | jump_event | message                 |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | no jump    | YOU HEAR WHUMP, WHUMP. |

  # Jumping Wumpus 003
  Scenario Outline: Jumping Wumpus 003 first jump landing on player can trample
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    And the next first jump player landing outcome is <landing_outcome>
    When the next turn begins
    Then the game is <game_status>
    And the displayed lines include <messages>

    Examples:
      | player_room | wumpus_room | jump_path | landing_outcome | pit_rooms | bat_rooms | arrows | game_status | messages                                                     |
      | 2           | 10          | 2, 1      | trample         | 13, 14    | 16, 17    | 5      | lost        | YOU HEAR WHUMP, WHUMP., THE WUMPUS TRAMPLES YOU TO DEATH!, HA HA HA - YOU LOSE! |
      | 2           | 10          | 2, 1      | slam            | 13, 14    | 16, 17    | 5      | in progress | YOU HEAR WHUMP, WHUMP., YOU ARE SLAMMED AGAINST THE CAVE WALL BY THE SNARLING WUMPUS! |

  # Jumping Wumpus 004
  Scenario Outline: Jumping Wumpus 004 second jump landing on player grants escape turn
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    When the next turn begins
    Then the game is still in progress
    And the player may take the next command
    And the displayed lines include <messages>

    Examples:
      | player_room | wumpus_room | jump_path | pit_rooms | bat_rooms | arrows | messages                                                                                |
      | 1           | 10          | 2, 1      | 13, 14    | 16, 17    | 5      | YOU HEAR WHUMP, WHUMP., YOU SEE THE BLOODSTAINED EYES OF THE WUMPUS APPRAISING YOU!     |

  # Jumping Wumpus 005
  Scenario Outline: Jumping Wumpus 005 jumps only follow legal tunnels
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    When the next turn begins
    Then every Wumpus jump segment is a legal tunnel
    And the Wumpus is in room <expected_wumpus_room>

    Examples:
      | player_room | wumpus_room | jump_path | expected_wumpus_room | pit_rooms | bat_rooms | arrows |
      | 1           | 10          | 11, 12    | 12                   | 13, 14    | 16, 17    | 5      |
      | 1           | 20          | 19, 18    | 18                   | 13, 14    | 16, 17    | 5      |

  # Jumping Wumpus 006
  Scenario Outline: Jumping Wumpus 006 jump can occur before a rest, move, or shoot command
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    And the next Wumpus wake choice is <wake_choice>
    When the player enters command <command>
    Then the displayed lines include <message>
    And the game is <game_status>

    Examples:
      | player_room | wumpus_room | jump_path | wake_choice | pit_rooms | bat_rooms | arrows | command | message                 | game_status |
      | 1           | 10          | 11, 12    | stay        | 13, 14    | 16, 17    | 5      | m 5     | YOU HEAR WHUMP, WHUMP. | in progress |
      | 1           | 10          | 11, 12    | stay        | 13, 14    | 16, 17    | 5      | s 5     | YOU HEAR WHUMP, WHUMP. | in progress |

  # Jumping Wumpus 007
  Scenario Outline: Jumping Wumpus 007 pending grenade still detonates after a jump turn
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    When the player enters command <command>
    Then the displayed lines include <messages>
    And no Holy Hand Grenade detonation is pending

    Examples:
      | player_room | wumpus_room | target_room | jump_path | pit_rooms | bat_rooms | arrows | command | messages                                                   |
      | 1           | 10          | 13          | 11, 12    | 14, 15    | 16, 17    | 5      | m 5     | YOU HEAR WHUMP, WHUMP., YOU HEAR A HORRENDOUS EXPLOSION!  |

  # Jumping Wumpus 008
  Scenario Outline: Jumping Wumpus 008 seeded jump events are reproducible
    Given a new game created with seed <seed>
    And another new game created with seed <seed>
    When both games evaluate jumping Wumpus behavior for <turn_count> turns
    Then both games produce identical jumping Wumpus events

    Examples:
      | seed | turn_count |
      | 1973 | 20         |
      | 2026 | 20         |
