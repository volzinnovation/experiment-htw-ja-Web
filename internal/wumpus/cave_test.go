package wumpus

import (
	"reflect"
	"slices"
	"testing"
)

func TestCaveExitsUseCanonicalDodecahedron(t *testing.T) {
	cave := NewCave()

	got, err := cave.Exits(1)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{2, 5, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("exits for room 1 = %v, want %v", got, want)
	}
}

func TestCaveRejectsInvalidRoom(t *testing.T) {
	_, err := NewCave().Exits(21)
	if err == nil {
		t.Fatal("expected invalid room error")
	}
}

func TestCaveReachabilityIncludesAllRooms(t *testing.T) {
	reachable := NewCave().ReachableFrom(1)

	if len(reachable) != 20 {
		t.Fatalf("reachable rooms = %v, want 20 rooms", reachable)
	}
	for room := 1; room <= 20; room++ {
		if !slices.Contains(reachable, room) {
			t.Fatalf("room %d was not reachable from room 1", room)
		}
	}
}

func TestTunnelQueriesReportExistingAndMissingLinks(t *testing.T) {
	cave := NewCave()

	if !cave.HasTunnel(1, 2) {
		t.Fatal("expected tunnel from 1 to 2")
	}
	if cave.HasTunnel(1, 20) {
		t.Fatal("did not expect tunnel from 1 to 20")
	}
	if cave.HasTunnel(21, 1) {
		t.Fatal("did not expect tunnel from invalid room")
	}
}

func TestReachableFromInvalidRoomReturnsNoRooms(t *testing.T) {
	reachable := NewCave().ReachableFrom(21)

	if len(reachable) != 0 {
		t.Fatalf("reachable rooms from invalid room = %v, want none", reachable)
	}
}

func TestAdjacentHazardsReportsOnlyNeighboringHazardTypes(t *testing.T) {
	setup := Setup{Player: 1, Wumpus: 2, Pits: []int{3, 4}, Bats: []int{5, 6}}

	got := NewCave().AdjacentHazards(setup.Player, setup)
	want := []Hazard{HazardWumpus, HazardBats}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("adjacent hazards = %v, want %v", got, want)
	}
}

func TestAdjacentHazardsCanReportPitAndIgnoreInvalidRoom(t *testing.T) {
	setup := Setup{Player: 10, Wumpus: 1, Pits: []int{9, 18}, Bats: []int{6, 7}}

	got := NewCave().AdjacentHazards(setup.Player, setup)
	want := []Hazard{HazardPit}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("adjacent hazards = %v, want %v", got, want)
	}

	if hazards := NewCave().AdjacentHazards(21, setup); len(hazards) != 0 {
		t.Fatalf("adjacent hazards from invalid room = %v, want none", hazards)
	}
}
