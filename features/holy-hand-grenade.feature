# mutation-stamp: sha256=19aa5bc1eac8070fcb2fd1d72fc75141a85cd3f67818a05b3dc2698fdd9b7b32
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Holy Hand Grenade","feature_path":"features/holy-hand-grenade.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:8c820d3e42dfcdf2c744b92243564cf145fad7244b538ec2266b4d4c5a8ad2f1","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Holy Hand Grenade
  The Holy Hand Grenade is a collectible one-use delayed blast weapon with exact turn timing and messages.

  # Holy Hand Grenade 001
  Scenario: Holy Hand Grenade 001 setup places one grenade in a valid empty room
    Given a new game created with seed 1973
    When the setup is inspected
    Then there is 1 Holy Hand Grenade
    And the Holy Hand Grenade room is from 1 through 20
    And the Holy Hand Grenade room is not occupied by the player, Wumpus, pits, or bats

  # Holy Hand Grenade 002
  Scenario: Holy Hand Grenade 002 another seed places one grenade in a valid empty room
    Given a new game created with seed 2026
    When the setup is inspected
    Then there is 1 Holy Hand Grenade
    And the Holy Hand Grenade room is from 1 through 20
    And the Holy Hand Grenade room is not occupied by the player, Wumpus, pits, or bats

  # Holy Hand Grenade 003
  Scenario: Holy Hand Grenade 003 same seed creates the same grenade placement
    Given a new game created with seed 1973
    And another new game created with seed 1973
    When both setups are inspected
    Then both setups have identical Holy Hand Grenade rooms

  # Holy Hand Grenade 004
  Scenario: Holy Hand Grenade 004 another same seed creates the same grenade placement
    Given a new game created with seed 2026
    And another new game created with seed 2026
    When both setups are inspected
    Then both setups have identical Holy Hand Grenade rooms

  # Holy Hand Grenade 005
  Scenario: Holy Hand Grenade 005 entering grenade room acquires it once
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and the Holy Hand Grenade in room 2
    When the player moves to room 2
    Then the player carries the Holy Hand Grenade
    And the cave has no unclaimed Holy Hand Grenade
    And the turn messages are YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!

  # Holy Hand Grenade 006
  Scenario: Holy Hand Grenade 006 grenade produces no adjacency warning
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and the Holy Hand Grenade in room 2
    When the turn warnings are requested
    Then the warning messages are none

  # Holy Hand Grenade 007
  Scenario: Holy Hand Grenade 007 armed prompt includes throw
    Given an interactive game setup where the player carries the Holy Hand Grenade and has 5 arrows
    When the next turn is displayed
    Then the displayed lines include SHOOT, MOVE OR THROW (S-M-T)?

  # Holy Hand Grenade 008
  Scenario: Holy Hand Grenade 008 player without grenade cannot throw
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    When the player enters command t 2
    Then the displayed lines include CAN'T THROW THERE
    And the player does not carry the Holy Hand Grenade
    And the game is still in progress

  # Holy Hand Grenade 009
  Scenario: Holy Hand Grenade 009 invalid throw target is rejected without consuming grenade
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the player carries the Holy Hand Grenade
    When the player enters command t 20
    Then the displayed lines include CAN'T THROW THERE
    And the player carries the Holy Hand Grenade
    And no Holy Hand Grenade detonation is pending

  # Holy Hand Grenade 010
  Scenario: Holy Hand Grenade 010 uppercase invalid throw target is rejected without consuming grenade
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the player carries the Holy Hand Grenade
    When the player enters command T 20
    Then the displayed lines include CAN'T THROW THERE
    And the player carries the Holy Hand Grenade
    And no Holy Hand Grenade detonation is pending

  # Holy Hand Grenade 011
  Scenario: Holy Hand Grenade 011 successful throw starts one-turn fuse
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the player carries the Holy Hand Grenade
    When the player enters command t 2
    Then the player does not carry the Holy Hand Grenade
    And the Holy Hand Grenade is pending detonation in room 2
    And the displayed lines include YOU HEAR TIC...TIC...
    And the Wumpus is alive

  # Holy Hand Grenade 012
  Scenario: Holy Hand Grenade 012 uppercase successful throw starts one-turn fuse
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the player carries the Holy Hand Grenade
    When the player enters command T 2
    Then the player does not carry the Holy Hand Grenade
    And the Holy Hand Grenade is pending detonation in room 2
    And the displayed lines include YOU HEAR TIC...TIC...
    And the Wumpus is alive

  # Holy Hand Grenade 013
  Scenario: Holy Hand Grenade 013 pending grenade detonates after a legal move
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command m 2
    Then the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And no Holy Hand Grenade detonation is pending

  # Holy Hand Grenade 014
  Scenario: Holy Hand Grenade 014 pending grenade detonates after a legal shot
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command s 5
    Then the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And no Holy Hand Grenade detonation is pending

  # Holy Hand Grenade 015
  Scenario: Holy Hand Grenade 015 blast kills Wumpus in target room
    Given an interactive game setup with the player in room 1, the Wumpus in room 10, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command m 5
    Then the player wins
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!

  # Holy Hand Grenade 016
  Scenario: Holy Hand Grenade 016 blast kills Wumpus in adjacent room
    Given an interactive game setup with the player in room 1, the Wumpus in room 2, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command m 5
    Then the player wins
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!

  # Holy Hand Grenade 017
  Scenario: Holy Hand Grenade 017 blast destroys bats but leaves pits
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 9, 14, bats in rooms 2, 16, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command m 8
    Then the remaining bat rooms are 16
    And the pit rooms are 9, 14
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!

  # Holy Hand Grenade 018
  Scenario: Holy Hand Grenade 018 player moving into blast radius loses
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    When the player enters command m 2
    Then the player loses
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Holy Hand Grenade 019
  Scenario: Holy Hand Grenade 019 player already in blast radius loses
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 2
    When the player enters command r
    Then the player loses
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Holy Hand Grenade 020
  Scenario: Holy Hand Grenade 020 Wumpus moving into future blast zone is killed
    Given an interactive game setup with the player in room 1, the Wumpus in room 11, pits in rooms 13, 14, bats in rooms 16, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    And the next Wumpus wake choice is move to 10
    When the player enters command s 5
    Then the Wumpus is in room 10
    And the player wins
    And the displayed lines include MISSED
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!

  # Holy Hand Grenade 021
  Scenario: Holy Hand Grenade 021 bat transport before detonation can move player into blast
    Given an interactive game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, bats in rooms 5, 17, and 5 arrows
    And the Holy Hand Grenade is pending detonation in room 10
    And the next bat relocation room is 2
    When the player enters command m 5
    Then the player loses
    And the displayed lines include ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!
    And the displayed lines include YOU HEAR A HORRENDOUS EXPLOSION!
    And the displayed lines include YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!
    And the displayed lines include HA HA HA - YOU LOSE!

  # Holy Hand Grenade 022
  Scenario: Holy Hand Grenade 022 same setup replay preserves grenade placement and pending state
    Given an interactive game setup with seed 1973
    And the player carries the Holy Hand Grenade
    And the Holy Hand Grenade is pending detonation in room 10
    And the player has lost
    When the player answers same setup prompt with y
    Then the replay setup has identical player, Wumpus, pit, bat, and Holy Hand Grenade rooms
    And the replay Holy Hand Grenade pending detonation room is 10
