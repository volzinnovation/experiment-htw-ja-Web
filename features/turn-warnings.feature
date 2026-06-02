# mutation-stamp: sha256=9a1a43e5e552e329eb8437b0e4789ca092d965217893a17dd5fcbf1c5ac41c98
# acceptance-mutation-manifest-begin
# {"version":1,"tested_at":"2026-06-02T15:23:48Z","feature_name":"Turn warnings","feature_path":"features/turn-warnings.feature","background_hash":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b","implementation_hash":"sha256:84eae0be346f4f4cc5d119ab61fcecb12264725995d9624acc03f2b1b3808d09","scenarios":[]}
# acceptance-mutation-manifest-end

Feature: Turn warnings
  At the start of a turn, only hazards in rooms adjacent to the player produce warnings.

  Scenario: Turn Warnings 001 warnings appear for adjacent Wumpus and bats
    Given a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 3, 4, and bats in rooms 5, 6
    When the turn warnings are requested
    Then the warning messages are I SMELL A WUMPUS, BATS NEARBY

  Scenario: Turn Warnings 002 warnings appear for adjacent Wumpus and pits
    Given a game setup with the player in room 10, the Wumpus in room 2, pits in rooms 9, 18, and bats in rooms 6, 7
    When the turn warnings are requested
    Then the warning messages are I SMELL A WUMPUS, I FEEL A DRAFT

  Scenario: Turn Warnings 003 warnings appear for adjacent bats and pits
    Given a game setup with the player in room 13, the Wumpus in room 7, pits in rooms 12, 14, and bats in rooms 20, 1
    When the turn warnings are requested
    Then the warning messages are BATS NEARBY, I FEEL A DRAFT

  Scenario: Turn Warnings 004 warnings appear for all adjacent hazard types
    Given a game setup with the player in room 6, the Wumpus in room 5, pits in rooms 7, 15, and bats in rooms 1, 2
    When the turn warnings are requested
    Then the warning messages are I SMELL A WUMPUS, BATS NEARBY, I FEEL A DRAFT

  Scenario: Turn Warnings 005 non-adjacent hazards produce no warnings near room one
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 13, 14, and bats in rooms 16, 17
    When the turn warnings are requested
    Then the warning messages are none

  Scenario: Turn Warnings 006 non-adjacent hazards produce no warnings near room ten
    Given a game setup with the player in room 10, the Wumpus in room 5, pits in rooms 13, 14, and bats in rooms 16, 17
    When the turn warnings are requested
    Then the warning messages are none

  Scenario: Turn Warnings 007 repeated adjacent hazard types produce one warning per type near room one
    Given a game setup with the player in room 1, the Wumpus in room 20, pits in rooms 2, 5, and bats in rooms 8, 17
    When the turn warnings are requested
    Then the warning messages are BATS NEARBY, I FEEL A DRAFT

  Scenario: Turn Warnings 008 repeated adjacent hazard types produce one warning per type near room ten
    Given a game setup with the player in room 10, the Wumpus in room 20, pits in rooms 2, 9, and bats in rooms 11, 17
    When the turn warnings are requested
    Then the warning messages are BATS NEARBY, I FEEL A DRAFT
