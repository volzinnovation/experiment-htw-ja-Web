Feature: Rest command
  Rest consumes a full turn without moving, shooting, or throwing, while normal turn effects still occur.

  # Rest Command 001
  Scenario Outline: Rest Command 001 rest leaves player in place and preserves arrows
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player is in room <player_room>
    And the player has <arrows> arrows
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | r       |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | R       |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | rest    |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | ReSt    |

  # Rest Command 002
  Scenario Outline: Rest Command 002 rest consumes a full turn
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the turn count is <turn_count>
    When the player enters command <command>
    Then the turn count is <expected_turn_count>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | turn_count | command | expected_turn_count |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 1          | r       | 2                   |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 7          | rest    | 8                   |

  # Rest Command 003
  Scenario Outline: Rest Command 003 rest turn still displays ordinary warnings before command
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the next turn is displayed
    Then the displayed lines include <warnings>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | warnings                                      |
      | 1           | 2           | 5, 14     | 8, 17     | 5      | I SMELL A WUMPUS, BATS NEARBY, I FEEL A DRAFT |

  # Rest Command 004
  Scenario Outline: Rest Command 004 rest can trigger Jumping Wumpus behavior
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is <jump_path>
    When the player enters command <command>
    Then the displayed lines include <message>
    And the Wumpus is in room <expected_wumpus_room>

    Examples:
      | player_room | wumpus_room | jump_path | expected_wumpus_room | pit_rooms | bat_rooms | arrows | command | message                 |
      | 1           | 10          | 11, 12    | 12                   | 13, 14    | 16, 17    | 5      | r       | YOU HEAR WHUMP, WHUMP. |

  # Rest Command 005
  Scenario Outline: Rest Command 005 rest after throwing grenade detonates at end of turn
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the displayed lines include <message>
    And no Holy Hand Grenade detonation is pending

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message                          |
      | 1           | 13          | 20          | 14, 15    | 16, 17    | 5      | r       | YOU HEAR A HORRENDOUS EXPLOSION! |

  # Rest Command 006
  Scenario Outline: Rest Command 006 rest outside grenade blast survives
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the game is still in progress
    And the player is in room <player_room>
    And the displayed lines include <message>

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message                          |
      | 1           | 13          | 20          | 14, 15    | 16, 17    | 5      | rest    | YOU HEAR A HORRENDOUS EXPLOSION! |

  # Rest Command 007
  Scenario Outline: Rest Command 007 rest inside grenade blast loses
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the player loses
    And the displayed lines include <messages>

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                                         |
      | 1           | 2           | 20          | 13, 14    | 16, 17    | 5      | r       | YOU HEAR A HORRENDOUS EXPLOSION!, YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!, HA HA HA - YOU LOSE! |
      | 1           | 1           | 20          | 13, 14    | 16, 17    | 5      | rest    | YOU HEAR A HORRENDOUS EXPLOSION!, YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!, HA HA HA - YOU LOSE! |

  # Rest Command 008
  Scenario Outline: Rest Command 008 rest alone does not resolve room hazards
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player is in room <player_room>
    And the game is still in progress
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | r       | none     |

  # Rest Command 009
  Scenario Outline: Rest Command 009 invalid rest syntax does not advance the game
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the turn count is <turn_count>
    When the player enters command <command>
    Then the displayed lines include <message>
    And the turn count is <turn_count>
    And the player is in room <player_room>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | turn_count | command | message            |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 1          | r 5     | R IS NOT A COMMAND |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | 1          | rest 12 | REST IS NOT A COMMAND |
