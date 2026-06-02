Feature: Interactive game loop
  A human can play Hunt the Wumpus through compact single-line commands while the domain remains free of input and output.

  # Interactive Game Loop 001
  Scenario Outline: Interactive Game Loop 001 each turn displays room, tunnels, warnings, arrows, and prompt
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the next turn is displayed
    Then the displayed lines are <lines>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | lines                                                                 |
      | 1           | 2           | 13, 14    | 5, 17     | 5      | I SMELL A WUMPUS, BATS NEARBY, YOU ARE IN ROOM 1, TUNNELS LEAD TO 2 5 8, ARROWS LEFT: 5, SHOOT OR MOVE (S-M)? |
      | 10          | 20          | 13, 14    | 16, 17    | 4      | YOU ARE IN ROOM 10, TUNNELS LEAD TO 2 9 11, ARROWS LEFT: 4, SHOOT OR MOVE (S-M)? |

  # Interactive Game Loop 002
  Scenario Outline: Interactive Game Loop 002 move command is case insensitive
    Given an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player is in room <to_room>
    And the game is still in progress

    Examples:
      | from_room | to_room | wumpus_room | pit_rooms | bat_rooms | arrows | command |
      | 1         | 2       | 20          | 13, 14    | 16, 17    | 5      | m 2     |
      | 1         | 2       | 20          | 13, 14    | 16, 17    | 5      | M 2     |

  # Interactive Game Loop 003
  Scenario Outline: Interactive Game Loop 003 shoot command is case insensitive
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player wins
    And the displayed lines include <message>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message                                                                 |
      | 1           | 2           | 13, 14    | 16, 17    | 5      | s 2     | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |
      | 1           | 2           | 13, 14    | 16, 17    | 5      | S 2     | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Interactive Game Loop 004
  Scenario Outline: Interactive Game Loop 004 invalid command reports error and does not advance state
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the displayed lines include <message>
    And the player is in room <player_room>
    And the player has <arrows> arrows
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message            |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | x       | X IS NOT A COMMAND |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | jump 2  | JUMP IS NOT A COMMAND |

  # Interactive Game Loop 005
  Scenario Outline: Interactive Game Loop 005 invalid move reports error and reprompts without advancing
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the displayed lines include <message>
    And the player is in room <player_room>
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message          |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | m 20    | CAN'T MOVE THERE |
      | 10          | 20          | 13, 14    | 16, 17    | 5      | m 14    | CAN'T MOVE THERE |

  # Interactive Game Loop 006
  Scenario Outline: Interactive Game Loop 006 invalid arrow command reports error and preserves arrows
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the displayed lines include <message>
    And the player has <arrows> arrows
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command       | message           |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | s             | CAN'T SHOOT THERE |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | s 2 3 4 5 6 7 | CAN'T SHOOT THERE |
      | 1           | 20          | 13, 14    | 16, 17    | 5      | s 21          | CAN'T SHOOT THERE |

  # Interactive Game Loop 007
  Scenario Outline: Interactive Game Loop 007 winning shot ends the game with victory text
    Given an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player wins
    And the displayed lines include <message>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | message                                                                 |
      | 1           | 2           | 13, 14    | 16, 17    | 5      | s 2     | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Interactive Game Loop 008
  Scenario Outline: Interactive Game Loop 008 losing move displays loss text and same setup prompt
    Given an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows
    When the player enters command <command>
    Then the player loses
    And the displayed lines include <messages>

    Examples:
      | from_room | wumpus_room | pit_rooms | bat_rooms | arrows | command | messages                                                                      |
      | 1         | 20          | 2, 14     | 16, 17    | 5      | m 2     | YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!, SAME SET UP (Y-N)?       |

  # Interactive Game Loop 009
  Scenario Outline: Interactive Game Loop 009 same setup replay preserves the cave placement
    Given an interactive game setup with seed <seed>
    And the player has lost
    When the player answers same setup prompt with <answer>
    Then the next game setup is <setup_relation> to the lost game setup

    Examples:
      | seed | answer | setup_relation |
      | 1973 | y      | identical      |
      | 1973 | Y      | identical      |
      | 1973 | n      | different      |

  # Interactive Game Loop 010
  Scenario Outline: Interactive Game Loop 010 instruction prompt can show or skip instructions
    Given a new interactive session
    When the player answers instructions prompt with <answer>
    Then the displayed lines are <lines>

    Examples:
      | answer | lines                                |
      | y      | WELCOME TO 'HUNT THE WUMPUS'         |
      | n      | none                                 |
