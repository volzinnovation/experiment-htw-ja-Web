# mutation-stamp: sha256=4335a12c8f521d9cdc4f8b56a28fc6fa9a57057bb7aad109fe03e901f4751c13
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Interactive game loop","feature_path":"features/interactive-game-loop.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:0273e640e74422b81b2f57d6733dd1b34fd98f6d4380c3db13f7769e820a312d","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Interactive game loop
  A human can play Hunt the Wumpus through compact single-line commands while the domain remains free of input and output.

  # Interactive Game Loop 001
  Scenario: Interactive Game Loop 001 each turn displays room, tunnels, warnings, arrows, and prompt
    Given an interactive game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, bats in rooms 5, 17, and 5 arrows
    When the next turn is displayed
    Then the displayed lines are I SMELL A WUMPUS, BATS NEARBY, YOU ARE IN ROOM 1, TUNNELS LEAD TO 2 5 8, ARROWS LEFT: 5, SHOOT OR MOVE (S-M)?

  # Interactive Game Loop 002
  Scenario: Interactive Game Loop 002 move command is case insensitive
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command M 2
    Then the player is in room 2
    And the game is still in progress

  # Interactive Game Loop 003
  Scenario: Interactive Game Loop 003 shoot command is case insensitive
    Given an interactive game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command S 2
    Then the player wins
    And the displayed lines include AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!

  # Interactive Game Loop 004
  Scenario: Interactive Game Loop 004 invalid command reports error and does not advance state
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command jump 2
    Then the displayed lines include JUMP IS NOT A COMMAND
    And the player is in room 1
    And the player has 5 arrows
    And the game is still in progress

  # Interactive Game Loop 005
  Scenario: Interactive Game Loop 005 invalid move reports error and reprompts without advancing
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command m 20
    Then the displayed lines include CAN'T MOVE THERE
    And the player is in room 1
    And the game is still in progress

  # Interactive Game Loop 006
  Scenario: Interactive Game Loop 006 invalid arrow command reports error and preserves arrows
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command s 2 3 4 5 6 7
    Then the displayed lines include CAN'T SHOOT THERE
    And the player has 5 arrows
    And the game is still in progress

  # Interactive Game Loop 007
  Scenario: Interactive Game Loop 007 winning shot ends the game with victory text
    Given an interactive game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command s 2
    Then the player wins
    And the displayed lines include AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!

  # Interactive Game Loop 008
  Scenario: Interactive Game Loop 008 losing move displays loss text and same setup prompt
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 2, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command m 2
    Then the player loses
    And the displayed lines include YYYIIIIEEEE . . . FELL IN PIT
    And the displayed lines include HA HA HA - YOU LOSE!
    And the displayed lines include SAME SET UP (Y-N)?

  # Interactive Game Loop 009
  Scenario: Interactive Game Loop 009 same setup replay preserves the cave placement
    Given an interactive game setup with seed 1973
    And the player has lost
    When the player answers same setup prompt with Y
    Then the next game setup is identical to the lost game setup

  # Interactive Game Loop 010
  Scenario: Interactive Game Loop 010 instruction prompt can show instructions
    Given a new interactive session
    When the player answers instructions prompt with y
    Then the displayed lines include WELCOME TO 'HUNT THE WUMPUS'
    And the displayed lines include THE WUMPUS LIVES IN A CAVE OF 20 ROOMS: EACH ROOM HAS 3 TUNNELS LEADING TO OTHER
    And the displayed lines include HAZARDS:
    And the displayed lines include WARNINGS:

  Scenario: Interactive Game Loop 011 instruction prompt can skip instructions
    Given a new interactive session
    When the player answers instructions prompt with n
    Then the displayed lines are none
