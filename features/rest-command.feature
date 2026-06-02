# mutation-stamp: sha256=08e20b712e6f2e9d2b2e98ec3fe80cdec9977aa5fe4c2aff064e3239d8258356
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Rest command","feature_path":"features/rest-command.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:e8ec88473271f539053ad1c8572f8419f7d138d40aea77c2c2337caee344cd79","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Rest command
  Rest consumes a full turn without moving, shooting, or throwing, while normal turn effects still occur.

  # Rest Command 001
  Scenario: Rest Command 001 short rest leaves player in place and preserves arrows
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command r
    Then the player is in room 1
    And the player has 5 arrows
    And the game is still in progress

  # Rest Command 002
  Scenario: Rest Command 002 uppercase short rest leaves player in place and preserves arrows
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command R
    Then the player is in room 1
    And the player has 5 arrows
    And the game is still in progress

  # Rest Command 003
  Scenario: Rest Command 003 long rest leaves player in place and preserves arrows
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command rest
    Then the player is in room 1
    And the player has 5 arrows
    And the game is still in progress

  # Rest Command 004
  Scenario: Rest Command 004 mixed-case long rest leaves player in place and preserves arrows
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command ReSt
    Then the player is in room 1
    And the player has 5 arrows
    And the game is still in progress

  # Rest Command 005
  Scenario: Rest Command 005 short rest consumes a full turn
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the turn count starts at 1
    When the player enters command r
    Then the turn count is 2

  # Rest Command 006
  Scenario: Rest Command 006 long rest consumes a full turn
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the turn count starts at 7
    When the player enters command rest
    Then the turn count is 8

  # Rest Command 007
  Scenario: Rest Command 007 rest turn still displays ordinary warnings before command
    Given an interactive game setup with the player in room 1, the Wumpus in room 2, pits in rooms 5, 14, bats in rooms 8, 17, and 5 arrows
    When the next turn is displayed
    Then the displayed lines include I SMELL A WUMPUS
    And the displayed lines include BATS NEARBY
    And the displayed lines include I FEEL A DRAFT

  # Rest Command 008
  Scenario: Rest Command 008 rest can trigger Jumping Wumpus behavior
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    When the player enters command r
    Then the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the Wumpus is in room 12

  # Rest Command 009
  Scenario: Rest Command 009 rest after throwing grenade detonates at end of turn
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 14, 15, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 13
    When the player enters command r
    Then the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And no Holy Hand Grenade detonation is pending

  # Rest Command 010
  Scenario: Rest Command 010 long rest outside grenade blast survives
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 14, 15, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 13
    When the player enters command rest
    Then the game is still in progress
    And the player is in room 1
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!

  # Rest Command 011
  Scenario: Rest Command 011 rest moving into grenade blast loses
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 2
    When the player enters command r
    Then the player loses
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Rest Command 012
  Scenario: Rest Command 012 rest already inside grenade blast loses
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 1
    When the player enters command rest
    Then the player loses
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Rest Command 013
  Scenario: Rest Command 013 rest alone does not resolve room hazards
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command r
    Then the player is in room 1
    And the game is still in progress
    And the turn messages are none

  # Rest Command 014
  Scenario: Rest Command 014 invalid short rest syntax does not advance the game
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the turn count starts at 1
    When the player enters command r 5
    Then the displayed lines include R IS NOT A COMMAND
    And the turn count starts at 1
    And the player is in room 1

  # Rest Command 015
  Scenario: Rest Command 015 invalid long rest syntax does not advance the game
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the turn count starts at 1
    When the player enters command rest 12
    Then the displayed lines include REST IS NOT A COMMAND
    And the turn count is 1
    And the player is in room 1
