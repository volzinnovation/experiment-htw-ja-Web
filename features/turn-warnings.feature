Feature: Turn warnings
  At the start of a turn, only hazards in rooms adjacent to the player produce warnings.

  # Turn Warnings 001
  Scenario Outline: Turn Warnings 001 warnings appear for adjacent hazards
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the turn warnings are requested
    Then the warning messages are <warnings>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | warnings                                      |
      | 1           | 2           | 3, 4      | 5, 6      | I SMELL A WUMPUS, BATS NEARBY                |
      | 10          | 2           | 9, 18     | 6, 7      | I SMELL A WUMPUS, I FEEL A DRAFT             |
      | 13          | 7           | 12, 14    | 20, 1     | BATS NEARBY, I FEEL A DRAFT                  |
      | 6           | 5           | 7, 15     | 1, 2      | I SMELL A WUMPUS, BATS NEARBY, I FEEL A DRAFT |

  # Turn Warnings 002
  Scenario Outline: Turn Warnings 002 non-adjacent hazards produce no warnings
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the turn warnings are requested
    Then the warning messages are <warnings>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | warnings |
      | 1           | 20          | 13, 14    | 16, 17    | none     |
      | 10          | 5           | 13, 14    | 16, 17    | none     |

  # Turn Warnings 003
  Scenario Outline: Turn Warnings 003 repeated adjacent hazard types produce one warning per type
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    When the turn warnings are requested
    Then the warning messages are <warnings>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | warnings                         |
      | 1           | 20          | 2, 5      | 8, 17     | BATS NEARBY, I FEEL A DRAFT      |
      | 10          | 20          | 2, 9      | 11, 17    | BATS NEARBY, I FEEL A DRAFT      |
