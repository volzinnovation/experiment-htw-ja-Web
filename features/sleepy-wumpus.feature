# mutation-stamp: sha256=eeaca055be765aa05d48b66f7f018359fcacad95a328e31178e967c82f607184
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Sleepy Wumpus","feature_path":"features/sleepy-wumpus.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:54aa2a0ab07be21fba3ce621e7456e484c32d0e545073fe735c16000876d1a04","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Sleepy Wumpus
  The Wumpus can be asleep, allowing noisy warnings, risky close encounters, and explicit wake transitions.

  # Sleepy Wumpus 001
  Scenario: Sleepy Wumpus 001 adjacent room can include snoring with normal smell
    Given a game setup with the player in room 6, the Wumpus in room 1, pits in rooms 13, 14, and bats in rooms 16, 17
    And the next sleepy Wumpus adjacent observation is asleep
    When the player moves to room 5
    Then the Wumpus sleep state is asleep
    And the turn warnings are I SMELL A WUMPUS, YOU HEAR HORRIBLE SNORING

  # Sleepy Wumpus 002
  Scenario: Sleepy Wumpus 002 adjacent room can leave Wumpus awake
    Given a game setup with the player in room 6, the Wumpus in room 1, pits in rooms 13, 14, and bats in rooms 16, 17
    And the next sleepy Wumpus adjacent observation is awake
    When the player moves to room 5
    Then the Wumpus sleep state is awake
    And the turn warnings are I SMELL A WUMPUS

  # Sleepy Wumpus 003
  Scenario: Sleepy Wumpus 003 moving away from sleeping Wumpus adjacency awakens it
    Given a game setup with the player in room 5, the Wumpus in room 1, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is asleep
    When the player moves to room 6
    Then the Wumpus sleep state is awake
    And the turn messages are YOU HEAR A SNORT AND "HUH?"

  # Sleepy Wumpus 004
  Scenario: Sleepy Wumpus 004 entering sleeping Wumpus room can wake and kill
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is asleep
    And the next sleeping Wumpus room entry outcome is wakes
    When the player moves to room 2
    Then the game is lost
    And the Wumpus sleep state is awake
    And the turn messages are YOU HEAR THE WUMPUS SAY "YUMMY BREAKFAST!", HA HA HA - YOU LOSE!

  # Sleepy Wumpus 005
  Scenario: Sleepy Wumpus 005 entering sleeping Wumpus room can leave it asleep
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is asleep
    And the next sleeping Wumpus room entry outcome is stays asleep
    When the player moves to room 2
    Then the game is in progress
    And the Wumpus sleep state is asleep
    And the turn messages are YOU SEE THE HUDDLED HORRIBLE SHAPE OF THE SLEEPING WUMPUS

  # Sleepy Wumpus 006
  Scenario: Sleepy Wumpus 006 leaving after seeing sleeping Wumpus awakens it
    Given a game setup with the player in room 2, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is asleep
    And the player has seen the sleeping Wumpus shape
    When the player moves to room 1
    Then the Wumpus sleep state is awake
    And the player is in room 1
    And the turn messages are YOU HEAR A PETULANT SCREAM!

  # Sleepy Wumpus 007
  Scenario: Sleepy Wumpus 007 awake Wumpus entry uses original wake behavior
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is awake
    And the next Wumpus wake choice is stay
    When the player moves to room 2
    Then the Wumpus is in room 2
    And the game is lost
    And the turn messages are TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE!

  # Sleepy Wumpus 008
  Scenario: Sleepy Wumpus 008 shooting an arrow wakes sleeping Wumpus after a miss
    Given a game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, and bats in rooms 16, 17
    And the Wumpus is asleep
    And the player has 5 arrows
    And the next Wumpus wake choice is move to 11
    When the player shoots the path 5
    Then the Wumpus sleep state is awake
    And the Wumpus is in room 11
    And the game is in progress

  # Sleepy Wumpus 009
  Scenario: Sleepy Wumpus 009 seeded sleepy observations are reproducible
    Given a new game created with seed 1973
    And another new game created with seed 1973
    When both games observe sleepy Wumpus behavior for 10 turns
    Then both games produce identical sleepy Wumpus observations

  # Sleepy Wumpus 010
  Scenario: Sleepy Wumpus 010 another seed produces reproducible sleepy observations
    Given a new game created with seed 2026
    And another new game created with seed 2026
    When both games observe sleepy Wumpus behavior for 10 turns
    Then both games produce identical sleepy Wumpus observations

  # Sleepy Wumpus 011
  Scenario: Sleepy Wumpus 011 bat transport into Wumpus adjacency can reveal snoring
    Given a game setup with the player in room 1, the Wumpus in room 1, pits in rooms 13, 14, and bats in rooms 2, 17
    And the next bat relocation room is 5
    And the next sleepy Wumpus adjacent observation is asleep
    When the player moves to room 2
    Then the player is in room 5
    And the Wumpus sleep state is asleep
    And the turn warnings are I SMELL A WUMPUS, YOU HEAR HORRIBLE SNORING
