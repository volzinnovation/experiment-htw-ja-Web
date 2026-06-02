# mutation-stamp: sha256=8b7cf047de37c50e568cab5f13702b150f6456231d079959c05f5b5aac06c761
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:47Z","feature_name":"Crooked arrow shooting","feature_path":"features/crooked-arrow.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:477fbe8b9d3c7867cc8cc5f3bc64199445e4116ce9c1b607562d6a75c5310f4c","scenarios":[{"index":2,"name":"Crooked Arrow 003 arrow can hit the player after a legal or deviated path","scenario_hash":"618b904b4f42a63a6d728389efb6f69fc7545be10950c491110a6909a7d76401","mutation_count":20,"result":{"Total":20,"Killed":20,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:46:28Z"},{"index":4,"name":"Crooked Arrow 005 shooting spends one arrow","scenario_hash":"393bc84ae76dba2a82c9bafe8d1f2be3413ef7f2ede77c6b808b3d3e24f83546","mutation_count":20,"result":{"Total":20,"Killed":20,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:46:28Z"},{"index":0,"name":"Crooked Arrow 001 arrow path that reaches the Wumpus wins immediately","scenario_hash":"95f8f22a652af1300a88fcb192a213dad58f0278f4c53d4ba02cb6e62adac4af","mutation_count":16,"result":{"Total":16,"Killed":16,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:45:13Z"},{"index":1,"name":"Crooked Arrow 002 invalid arrow segment deviates to a random adjacent room","scenario_hash":"75b454a93441f149ca5c96a2e26ece255a9e9ecfc9f90e6c92dedde81e4926b6","mutation_count":10,"result":{"Total":10,"Killed":10,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:45:13Z"},{"index":3,"name":"Crooked Arrow 004 missed arrow wakes the Wumpus","scenario_hash":"c97b9d5e587c1b49586fa8bbba46a8bc098c05d9f3d9254720049fbec9d9dcea","mutation_count":20,"result":{"Total":20,"Killed":20,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:45:13Z"},{"index":5,"name":"Crooked Arrow 006 missing with the last arrow loses","scenario_hash":"8a3353b5bd74c644080e638f4e03b557cdcad1cb97d3157d26b08a6b33eee23c","mutation_count":9,"result":{"Total":9,"Killed":9,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:45:13Z"},{"index":6,"name":"Crooked Arrow 007 shooting path must contain one to five rooms","scenario_hash":"b878ab6b52807afd90fc54958c8711fd815253a549134fdd5a08f15f27a376d2","mutation_count":18,"result":{"Total":18,"Killed":18,"Survived":0,"Errors":0},"tested_at":"2026-06-02T13:45:13Z"}]}
# acceptance-mutation-manifest-end

Feature: Crooked arrow shooting
  The player shoots one to five-room crooked arrows that follow legal tunnels or deviate randomly on invalid segments.

  # Crooked Arrow 001
  Scenario Outline: Crooked Arrow 001 arrow path that reaches the Wumpus wins immediately
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    When the player shoots the path <path>
    Then the player wins
    And the requested shot path is <expected_path>
    And the arrow traversed rooms are <traversed_rooms>
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | initial_arrows | path       | expected_path | traversed_rooms | remaining_arrows | messages                                                                 |
      | 1           | 5           | 5              | 5          | 5             | 5               | 4                | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |
      | 6           | 18          | 5              | 7, 8, 9, 18 | 7, 8, 9, 18  | 7, 8, 9, 18    | 4                | AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!! |

  # Crooked Arrow 002
  Scenario Outline: Crooked Arrow 002 invalid arrow segment deviates to a random adjacent room
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    And the next arrow deviation room is <deviation_room>
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the arrow traversed rooms are <traversed_rooms>
    And the requested shot path is <expected_path>
    And the game is <game_status>
    And the player has <remaining_arrows> arrows

    Examples:
      | player_room | wumpus_room | initial_arrows | path | expected_path | deviation_room | traversed_rooms | wake_choice | game_status | remaining_arrows |
      | 1           | 20          | 5              | 3, 4 | 3, 4          | 5              | 5, 4            | stay        | in progress | 4                |

  # Crooked Arrow 003
  Scenario Outline: Crooked Arrow 003 arrow can hit the player after a legal or deviated path
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    And the next arrow deviation room is <deviation_room>
    When the player shoots the path <path>
    Then the player loses
    And the requested shot path is <expected_path>
    And the Wumpus is in room <expected_wumpus_room>
    And the arrow traversed rooms are <traversed_rooms>
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | initial_arrows | path | expected_path | deviation_room | expected_wumpus_room | traversed_rooms | remaining_arrows | messages                                    |
      | 6           | 20          | 5              | 7, 6 | 7, 6          | none           | 20                   | 7, 6            | 4                | OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE! |
      | 1           | 20          | 5              | 3   | 3             | 1              | 20                   | 1               | 4                | OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE! |

  # Crooked Arrow 004
  Scenario Outline: Crooked Arrow 004 missed arrow wakes the Wumpus
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the Wumpus is in room <expected_wumpus_room>
    And the requested shot path is <expected_path>
    And the game is <game_status>
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | initial_arrows | path | expected_path | wake_choice | expected_wumpus_room | game_status | remaining_arrows | messages                                                    |
      | 1           | 10          | 5              | 5    | 5             | stay        | 10                   | in progress | 4                | MISSED                                                      |
      | 1           | 2           | 5              | 5    | 5             | move to 1   | 1                    | lost        | 4                | MISSED, TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE! |

  # Crooked Arrow 005
  Scenario Outline: Crooked Arrow 005 shooting spends one arrow
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the player has <remaining_arrows> arrows
    And the requested shot path is <expected_path>
    And the player is in room <expected_player_room>
    And the Wumpus is in room <expected_wumpus_room>
    And the game is <game_status>

    Examples:
      | player_room | wumpus_room | initial_arrows | path | expected_path | wake_choice | remaining_arrows | expected_player_room | expected_wumpus_room | game_status |
      | 1           | 10          | 5              | 5    | 5             | stay        | 4                | 1                    | 10                   | in progress |
      | 1           | 10          | 2              | 5    | 5             | stay        | 1                | 1                    | 10                   | in progress |

  # Crooked Arrow 006
  Scenario Outline: Crooked Arrow 006 missing with the last arrow loses
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    And the next Wumpus wake choice is <wake_choice>
    When the player shoots the path <path>
    Then the player loses
    And the requested shot path is <expected_path>
    And the Wumpus is in room <expected_wumpus_room>
    And the player has <remaining_arrows> arrows
    And the turn messages are <messages>

    Examples:
      | player_room | wumpus_room | initial_arrows | path | expected_path | wake_choice | expected_wumpus_room | remaining_arrows | messages                                      |
      | 1           | 10          | 1              | 5    | 5             | stay        | 10                   | 0                | MISSED, YOU RAN OUT OF ARROWS, HA HA HA - YOU LOSE! |

  # Crooked Arrow 007
  Scenario Outline: Crooked Arrow 007 shooting path must contain one to five rooms
    Given a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>
    And the player starts with <initial_arrows> arrows
    When the player shoots the path <path>
    Then the shot is rejected with message <message>
    And the requested shot path is <expected_path>
    And the player is in room <expected_player_room>
    And the Wumpus is in room <expected_wumpus_room>
    And the player has <remaining_arrows> arrows
    And the game is still in progress

    Examples:
      | player_room | wumpus_room | initial_arrows | path             | expected_path     | expected_player_room | expected_wumpus_room | remaining_arrows | message           |
      | 1           | 10          | 5              | none             | none              | 1                    | 10                   | 5                | CAN'T SHOOT THERE |
      | 1           | 10          | 5              | 2, 3, 4, 5, 6, 7 | 2, 3, 4, 5, 6, 7 | 1                    | 10                   | 5                | CAN'T SHOOT THERE |
