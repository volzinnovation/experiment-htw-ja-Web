# mutation-stamp: sha256=d837a258f0da5e66e70d78a2b64800a195243420057c02dc9c2f1dde0c94015a
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Jumping Wumpus","feature_path":"features/jumping-wumpus.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:e3b618797b2b6f8becf15b078127a49e7e4658b5ba8099a2f6b627c5291ea78b","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Jumping Wumpus
  On a turn, the Wumpus can make two legal jumps with exact messages and landing outcomes.

  # Jumping Wumpus 001
  Scenario: Jumping Wumpus 001 seeded turn can trigger double jump event
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    When the next turn begins
    Then the Wumpus is in room 12
    And the displayed lines include YOU HEAR WHUMP, WHUMP.

  # Jumping Wumpus 002
  Scenario: Jumping Wumpus 002 no jump event leaves Wumpus in place
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is no jump
    When the next turn begins
    Then the Wumpus is in room 10
    And the displayed lines do not include YOU HEAR WHUMP, WHUMP.

  # Jumping Wumpus 003
  Scenario: Jumping Wumpus 003 first jump landing on player can trample
    Given an interactive game setup with the player in room 2, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 2, 1
    And the next first jump player landing outcome is trample
    When the next turn begins
    Then the game is lost
    And the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the displayed lines include THE WUMPUS TRAMPLES YOU TO DEATH!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Jumping Wumpus 004
  Scenario: Jumping Wumpus 004 first jump landing on player can slam
    Given an interactive game setup with the player in room 2, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 2, 1
    And the next first jump player landing outcome is slam
    When the next turn begins
    Then the game is in progress
    And the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the displayed lines include YOU ARE SLAMMED AGAINST THE CAVE WALL BY THE SNARLING WUMPUS!

  # Jumping Wumpus 005
  Scenario: Jumping Wumpus 005 second jump landing on player grants escape turn
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 2, 1
    When the next turn begins
    Then the game is still in progress
    And the player may take the next command
    And the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the displayed lines include YOU SEE THE BLOODSTAINED EYES OF THE WUMPUS APPRAISING YOU!

  # Jumping Wumpus 006
  Scenario: Jumping Wumpus 006 jumps only follow legal tunnels
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    When the next turn begins
    Then every Wumpus jump segment is a legal tunnel
    And the Wumpus is in room 12

  # Jumping Wumpus 007
  Scenario: Jumping Wumpus 007 another jump path only follows legal tunnels
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 19, 18
    When the next turn begins
    Then every Wumpus jump segment is a legal tunnel
    And the Wumpus is in room 18

  # Jumping Wumpus 008
  Scenario: Jumping Wumpus 008 jump can occur before a move command
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    And the next Wumpus wake choice is stay
    When the player enters command m 5
    Then the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the game is in progress

  # Jumping Wumpus 009
  Scenario: Jumping Wumpus 009 jump can occur before a shoot command
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    And the next Wumpus wake choice is stay
    When the player enters command s 5
    Then the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the game is in progress

  # Jumping Wumpus 010
  Scenario: Jumping Wumpus 010 pending grenade still detonates after a jump turn
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 14, 15, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 13
    And the next jumping Wumpus turn event is jumps
    And the next Wumpus jump path is 11, 12
    When the player enters command m 5
    Then the displayed lines include YOU HEAR WHUMP, WHUMP.
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And no Holy Hand Grenade detonation is pending

  # Jumping Wumpus 011
  Scenario: Jumping Wumpus 011 seeded jump events are reproducible
    Given a new game created with seed 1973
    And another new game created with seed 1973
    When both games evaluate jumping Wumpus behavior for 20 turns
    Then both games produce identical jumping Wumpus events

  # Jumping Wumpus 012
  Scenario: Jumping Wumpus 012 another seed produces reproducible jump events
    Given a new game created with seed 2026
    And another new game created with seed 2026
    When both games evaluate jumping Wumpus behavior for 20 turns
    Then both games produce identical jumping Wumpus events
