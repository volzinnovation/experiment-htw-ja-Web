package wumpus

import (
	"reflect"
	"testing"
)

func TestNewGamePlacesRequiredOccupantsInDistinctRooms(t *testing.T) {
	game, err := NewGame(1973)
	if err != nil {
		t.Fatal(err)
	}

	setup := game.Setup()
	if setup.Player == 0 {
		t.Fatal("player was not placed")
	}
	if setup.Wumpus == 0 {
		t.Fatal("Wumpus was not placed")
	}
	if len(setup.Pits) != 2 {
		t.Fatalf("pit count = %d, want 2", len(setup.Pits))
	}
	if len(setup.Bats) != 2 {
		t.Fatalf("bat count = %d, want 2", len(setup.Bats))
	}
	if distinctCount(setup.OccupiedRooms()) != 6 {
		t.Fatalf("occupied rooms = %v, want 6 distinct rooms", setup.OccupiedRooms())
	}
}

func TestNewGameIsReproducibleBySeed(t *testing.T) {
	first, err := NewGame(2026)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewGame(2026)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(first.Setup(), second.Setup()) {
		t.Fatalf("same seed produced %v and %v", first.Setup(), second.Setup())
	}
}

func TestReplaySameSetupPreservesPlacement(t *testing.T) {
	game, err := NewGame(1)
	if err != nil {
		t.Fatal(err)
	}

	replay := game.ReplaySameSetup()

	if !reflect.DeepEqual(game.Setup(), replay.Setup()) {
		t.Fatalf("replay setup = %v, want %v", replay.Setup(), game.Setup())
	}
}

func TestConfiguredSetupRequiresValidDistinctRooms(t *testing.T) {
	_, err := NewGameWithSetup(Setup{Player: 1, Wumpus: 1, Pits: []int{2, 3}, Bats: []int{4, 5}})
	if err == nil {
		t.Fatal("expected duplicate-room setup error")
	}
}

func TestConfiguredSetupRequiresTwoPits(t *testing.T) {
	_, err := NewGameWithSetup(Setup{Player: 1, Wumpus: 2, Pits: []int{3}, Bats: []int{4, 5}})
	if err == nil {
		t.Fatal("expected pit-count setup error")
	}
}

func TestConfiguredSetupRequiresTwoBats(t *testing.T) {
	_, err := NewGameWithSetup(Setup{Player: 1, Wumpus: 2, Pits: []int{3, 4}, Bats: []int{5}})
	if err == nil {
		t.Fatal("expected bat-count setup error")
	}
}

func TestConfiguredSetupRequiresValidRooms(t *testing.T) {
	_, err := NewGameWithSetup(Setup{Player: 0, Wumpus: 2, Pits: []int{3, 4}, Bats: []int{5, 6}})
	if err == nil {
		t.Fatal("expected invalid-room setup error")
	}
}

func distinctCount(rooms []int) int {
	distinct := map[int]bool{}
	for _, room := range rooms {
		distinct[room] = true
	}
	return len(distinct)
}
