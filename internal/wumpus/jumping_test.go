package wumpus

import (
	"reflect"
	"testing"
)

func TestJumpingWumpusMovesAlongTwoRoomPath(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{11, 12})

	result := game.ResolveJumpingWumpusTurn()

	if game.Setup().Wumpus != 12 {
		t.Fatalf("Wumpus room = %d, want 12", game.Setup().Wumpus)
	}
	if !reflect.DeepEqual(result.Messages, []string{"YOU HEAR WHUMP, WHUMP."}) {
		t.Fatalf("messages = %v", result.Messages)
	}
	if !reflect.DeepEqual(result.JumpedRooms, []int{11, 12}) {
		t.Fatalf("jumped rooms = %v", result.JumpedRooms)
	}
}

func TestNoJumpingWumpusEventLeavesWumpusInPlace(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextJumpingWumpusTurnEvent(false)

	result := game.ResolveJumpingWumpusTurn()

	if game.Setup().Wumpus != 10 {
		t.Fatalf("Wumpus room = %d, want 10", game.Setup().Wumpus)
	}
	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
}

func TestFirstJumpLandingOnPlayerCanTrample(t *testing.T) {
	game := mustGame(t, Setup{Player: 2, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{2, 1})
	game.SetNextFirstJumpPlayerLandingOutcome(FirstJumpTramples)

	result := game.ResolveJumpingWumpusTurn()

	want := []string{"YOU HEAR WHUMP, WHUMP.", "THE WUMPUS TRAMPLES YOU TO DEATH!", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestFirstJumpLandingOnPlayerCanSlamAndContinue(t *testing.T) {
	game := mustGame(t, Setup{Player: 2, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{2, 1})
	game.SetNextFirstJumpPlayerLandingOutcome(FirstJumpSlams)

	result := game.ResolveJumpingWumpusTurn()

	want := []string{"YOU HEAR WHUMP, WHUMP.", "YOU ARE SLAMMED AGAINST THE CAVE WALL BY THE SNARLING WUMPUS!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
	if game.Setup().Wumpus != 1 {
		t.Fatalf("Wumpus room = %d, want 1", game.Setup().Wumpus)
	}
}

func TestSecondJumpLandingOnPlayerAllowsEscapeTurn(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{2, 1})

	result := game.ResolveJumpingWumpusTurn()

	want := []string{"YOU HEAR WHUMP, WHUMP.", "YOU SEE THE BLOODSTAINED EYES OF THE WUMPUS APPRAISING YOU!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
}

func TestJumpingWumpusEventsAreReproducibleForSameSetup(t *testing.T) {
	first := mustGame(t, Setup{Player: 3, Wumpus: 13, Pits: []int{15, 5}, Bats: []int{1, 18}})
	second := mustGame(t, Setup{Player: 3, Wumpus: 13, Pits: []int{15, 5}, Bats: []int{1, 18}})

	left := first.ObserveJumpingWumpusBehavior(20)
	right := second.ObserveJumpingWumpusBehavior(20)

	if !reflect.DeepEqual(left, right) {
		t.Fatalf("jump events differ: %v and %v", left, right)
	}
}
