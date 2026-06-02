Feature: Holy Hand Grenade
  The Holy Hand Grenade is a collectible one-use delayed blast weapon with exact turn timing and messages.

  # Holy Hand Grenade 001
  Scenario Outline: Holy Hand Grenade 001 setup places one grenade in a valid empty room
    Given a new game created with seed <seed>
    When the setup is inspected
    Then there is 1 Holy Hand Grenade
    And the Holy Hand Grenade room is from 1 through 20
    And the Holy Hand Grenade room is not occupied by the player, Wumpus, pits, or bats

    Examples:
      | seed |
      | 1973 |
      | 2026 |

  # Holy Hand Grenade 002
  Scenario Outline: Holy Hand Grenade 002 same seed creates the same grenade placement
    Given a new game created with seed <seed>
    And another new game created with seed <seed>
    When both setups are inspected
    Then both setups have identical Holy Hand Grenade rooms

    Examples:
      | seed |
      | 1973 |
      | 2026 |

  # Holy Hand Grenade 003
  Scenario Outline: Holy Hand Grenade 003 entering grenade room acquires it once
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and the Holy Hand Grenade in room <grenade_room>
    When the player moves to room <grenade_room>
    Then the player carries the Holy Hand Grenade
    And the cave has no unclaimed Holy Hand Grenade
    And the turn messages are <messages>

    Examples:
      | from_room | grenade_room | wumpus_room | pit_rooms | bat_rooms | messages                                                |
      | 1         | 2            | 20          | 13, 14    | 16, 17    | YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!         |

  # Holy Hand Grenade 004
  Scenario Outline: Holy Hand Grenade 004 grenade produces no adjacency warning
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and the Holy Hand Grenade in room <grenade_room>
    When the turn warnings are requested
    Then the warning messages are <warnings>

    Examples:
      | player_room | grenade_room | wumpus_room | pit_rooms | bat_rooms | warnings |
      | 1           | 2            | 20          | 13, 14    | 16, 17    | none     |

  # Holy Hand Grenade 005
  Scenario Outline: Holy Hand Grenade 005 armed prompt includes throw
    Given an interactive game setup where the player carries the Holy Hand Grenade and has <arrows> arrows
    When the next turn is displayed
    Then the displayed lines include <prompt>

    Examples:
      | arrows | prompt                         |
      | 5      | SHOOT, MOVE OR THROW (S-M-T)?  |

  # Holy Hand Grenade 006
  Scenario Outline: Holy Hand Grenade 006 player without grenade cannot throw
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the displayed lines include <message>
    And the player does not carry the Holy Hand Grenade
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message           |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | t 2     | CAN'T THROW THERE |

  # Holy Hand Grenade 007
  Scenario Outline: Holy Hand Grenade 007 invalid throw target is rejected without consuming grenade
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the player carries the Holy Hand Grenade
    When the player enters command <command>
    Then the displayed lines include <message>
    And the player carries the Holy Hand Grenade
    And no Holy Hand Grenade detonation is pending

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message           |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | t 20    | CAN'T THROW THERE |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | T 20    | CAN'T THROW THERE |

  # Holy Hand Grenade 008
  Scenario Outline: Holy Hand Grenade 008 successful throw starts one-turn fuse
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the player carries the Holy Hand Grenade
    When the player enters command <command>
    Then the player does not carry the Holy Hand Grenade
    And the Holy Hand Grenade is pending detonation in room <target_room>
    And the displayed lines include <message>
    And the Wumpus is alive

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message               |
      | 1           | 2           | 20          | 13, 14    | 16, 17    | 5      | t 2     | YOU HEAR TIC...TIC... |
      | 1           | 2           | 20          | 13, 14    | 16, 17    | 5      | T 2     | YOU HEAR TIC...TIC... |

  # Holy Hand Grenade 009
  Scenario Outline: Holy Hand Grenade 009 pending grenade detonates after the next legal turn
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the displayed lines include <message>
    And no Holy Hand Grenade detonation is pending

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message                                  |
      | 1           | 10          | 20          | 13, 14    | 16, 17    | 5      | m 2     | YOU HEAR A HORRENDOUS EXPLOSION!         |
      | 1           | 10          | 20          | 13, 14    | 16, 17    | 5      | s 5     | YOU HEAR A HORRENDOUS EXPLOSION!         |

  # Holy Hand Grenade 010
  Scenario Outline: Holy Hand Grenade 010 blast kills Wumpus in target room or adjacent room
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the player wins
    And the displayed lines include <messages>

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                       |
      | 1           | 10          | 10          | 13, 14    | 16, 17    | 5      | m 5     | YOU HEAR A HORRENDOUS EXPLOSION!, AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |
      | 1           | 10          | 2           | 13, 14    | 16, 17    | 5      | m 5     | YOU HEAR A HORRENDOUS EXPLOSION!, AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Holy Hand Grenade 011
  Scenario Outline: Holy Hand Grenade 011 blast destroys bats but leaves pits
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the remaining bat rooms are <remaining_bat_rooms>
    And the pit rooms are <pit_rooms>
    And the displayed lines include <message>

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | remaining_bat_rooms | message                          |
      | 1           | 10          | 20          | 9, 14     | 2, 16     | 5      | m 5     | 16                  | YOU HEAR A HORRENDOUS EXPLOSION! |

  # Holy Hand Grenade 012
  Scenario Outline: Holy Hand Grenade 012 player in blast radius loses
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    When the player enters command <command>
    Then the player loses
    And the displayed lines include <messages>

    Examples:
      | player_room | target_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                                         |
      | 1           | 10          | 20          | 13, 14    | 16, 17    | 5      | m 2     | YOU HEAR A HORRENDOUS EXPLOSION!, YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!, HA HA HA - YOU LOSE! |
      | 1           | 10          | 20          | 13, 14    | 16, 17    | 5      | m 5     | YOU HEAR A HORRENDOUS EXPLOSION!, YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!, HA HA HA - YOU LOSE! |

  # Holy Hand Grenade 013
  Scenario Outline: Holy Hand Grenade 013 Wumpus moving into future blast zone is killed
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    And the next Wumpus wake choice is <wake_choice>
    When the player enters command <command>
    Then the Wumpus is in room <expected_wumpus_room>
    And the player wins
    And the displayed lines include <messages>

    Examples:
      | player_room | target_room | wumpus_room | wake_choice | expected_wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                       |
      | 1           | 10          | 11          | move to 10  | 10                   | 13, 14    | 16, 17    | 5      | s 5     | MISSED, YOU HEAR A HORRENDOUS EXPLOSION!, AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Holy Hand Grenade 014
  Scenario Outline: Holy Hand Grenade 014 bat transport before detonation can move player into blast
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    And the Holy Hand Grenade is pending detonation in room <target_room>
    And the next bat relocation room is <relocation_room>
    When the player enters command <command>
    Then the player loses
    And the displayed lines include <messages>

    Examples:
      | player_room | target_room | relocation_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                                         |
      | 1           | 10          | 2               | 20          | 13, 14    | 5, 17     | 5      | m 5     | ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, YOU HEAR A HORRENDOUS EXPLOSION!, YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!, HA HA HA - YOU LOSE! |

  # Holy Hand Grenade 015
  Scenario Outline: Holy Hand Grenade 015 same setup replay preserves grenade placement and pending state
    Given an interactive game setup with seed <seed>
    And the player carries the Holy Hand Grenade
    And the Holy Hand Grenade is pending detonation in room <target_room>
    And the player has lost
    When the player answers same setup prompt with <answer>
    Then the replay setup has identical player, Wumpus, pit, bat, and Holy Hand Grenade rooms
    And the replay Holy Hand Grenade pending detonation room is <target_room>

    Examples:
      | seed | target_room | answer |
      | 1973 | 10          | y      |
