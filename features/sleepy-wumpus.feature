Feature: Sleepy Wumpus
  The Wumpus can be asleep, allowing noisy warnings, risky close encounters, and explicit wake transitions.

  # Sleepy Wumpus 001
  Scenario Outline: Sleepy Wumpus 001 adjacent room can include snoring with normal smell
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next sleepy Wumpus adjacent observation is <sleepy_observation>
    When the player moves to room <to_room>
    Then the Wumpus sleep state is <sleep_state>
    And the turn warnings are <warnings>

    Examples:
      | from_room | to_room | wumpus_room | pit_rooms | bat_rooms | sleepy_observation | sleep_state | warnings                                           |
      | 6         | 5       | 1           | 13, 14    | 16, 17    | asleep             | asleep      | I SMELL A WUMPUS, YOU HEAR HORRIBLE SNORING       |
      | 6         | 5       | 1           | 13, 14    | 16, 17    | awake              | awake       | I SMELL A WUMPUS                                  |

  # Sleepy Wumpus 002
  Scenario Outline: Sleepy Wumpus 002 moving away from sleeping Wumpus adjacency awakens it
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the Wumpus is asleep
    When the player moves to room <to_room>
    Then the Wumpus sleep state is awake
    And the turn messages are <messages>

    Examples:
      | from_room | to_room | wumpus_room | pit_rooms | bat_rooms | messages                       |
      | 5         | 6       | 1           | 13, 14    | 16, 17    | YOU HEAR A SNORT AND "HUH?"    |

  # Sleepy Wumpus 003
  Scenario Outline: Sleepy Wumpus 003 entering sleeping Wumpus room can wake and kill
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the Wumpus is asleep
    And the next sleeping Wumpus room entry outcome is <entry_outcome>
    When the player moves to room <wumpus_room>
    Then the game is <game_status>
    And the Wumpus sleep state is <sleep_state>
    And the turn messages are <messages>

    Examples:
      | from_room | wumpus_room | pit_rooms | bat_rooms | entry_outcome | game_status | sleep_state | messages                                                                      |
      | 1         | 2           | 13, 14    | 16, 17    | wakes         | lost        | awake       | YOU HEAR THE WUMPUS SAY "YUMMY BREAKFAST!", HA HA HA - YOU LOSE!              |
      | 1         | 2           | 13, 14    | 16, 17    | stays asleep  | in progress | asleep      | YOU SEE THE HUDDLED HORRIBLE SHAPE OF THE SLEEPING WUMPUS                    |

  # Sleepy Wumpus 004
  Scenario Outline: Sleepy Wumpus 004 leaving after seeing sleeping Wumpus awakens it
    Given a game setup with the player in room <wumpus_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the Wumpus is asleep
    And the player has seen the sleeping Wumpus shape
    When the player moves to room <to_room>
    Then the Wumpus sleep state is awake
    And the player is in room <to_room>
    And the turn messages are <messages>

    Examples:
      | wumpus_room | to_room | pit_rooms | bat_rooms | messages                         |
      | 2           | 1       | 13, 14    | 16, 17    | YOU HEAR A PETULANT SCREAM!      |

  # Sleepy Wumpus 005
  Scenario Outline: Sleepy Wumpus 005 awake Wumpus entry uses original wake behavior
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the Wumpus is awake
    And the next Wumpus wake choice is <wake_choice>
    When the player moves to room <wumpus_room>
    Then the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>
    And the turn messages are <messages>

    Examples:
      | from_room | wumpus_room | wake_choice | expected_wumpus_room | pit_rooms | bat_rooms | game_status | messages                                            |
      | 1         | 2           | stay        | 2                    | 13, 14    | 16, 17    | lost        | TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE! |

  # Sleepy Wumpus 006
  Scenario Outline: Sleepy Wumpus 006 shooting an arrow wakes sleeping Wumpus after a miss
    Given a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the Wumpus is asleep
    And the player has <arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the Wumpus sleep state is awake
    And the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>

    Examples:
      | player_room | wumpus_room | pit_rooms | bat_rooms | arrows | path | wake_choice | expected_wumpus_room | game_status |
      | 1           | 10          | 13, 14    | 16, 17    | 5      | 5    | move to 11  | 11                   | in progress |

  # Sleepy Wumpus 007
  Scenario Outline: Sleepy Wumpus 007 seeded sleepy observations are reproducible
    Given a new game created with seed <seed>
    And another new game created with seed <seed>
    When both games observe sleepy Wumpus behavior for <turn_count> turns
    Then both games produce identical sleepy Wumpus observations

    Examples:
      | seed | turn_count |
      | 1973 | 10         |
      | 2026 | 10         |

  # Sleepy Wumpus 008
  Scenario Outline: Sleepy Wumpus 008 bat transport into Wumpus adjacency can reveal snoring
    Given a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>
    And the next bat relocation room is <relocation_room>
    And the next sleepy Wumpus adjacent observation is <sleepy_observation>
    When the player moves to room <bat_room>
    Then the player is in room <relocation_room>
    And the Wumpus sleep state is <sleep_state>
    And the turn warnings are <warnings>

    Examples:
      | from_room | bat_room | relocation_room | wumpus_room | pit_rooms | bat_rooms | sleepy_observation | sleep_state | warnings                                     |
      | 1         | 2        | 5               | 1           | 13, 14    | 2, 17     | asleep             | asleep      | I SMELL A WUMPUS, YOU HEAR HORRIBLE SNORING |
