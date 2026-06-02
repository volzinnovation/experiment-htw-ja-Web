package wumpus

import (
	"reflect"
	"testing"
)

func TestArrowPathThatReachesWumpusWins(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)

	result := game.Shoot([]int{2, 10})

	wantMessages := []string{"AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"}
	if !reflect.DeepEqual(result.Messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", result.Messages, wantMessages)
	}
	if game.Status() != StatusWon {
		t.Fatalf("status = %s, want %s", game.Status(), StatusWon)
	}
	if game.Arrows() != 4 {
		t.Fatalf("arrows = %d, want 4", game.Arrows())
	}
}

func TestInvalidArrowSegmentUsesNextDeviationRoom(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)
	game.SetNextArrowDeviation(5)
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Stay: true})

	result := game.Shoot([]int{3, 4})

	wantTraversed := []int{5, 4}
	if !reflect.DeepEqual(result.TraversedRooms, wantTraversed) {
		t.Fatalf("traversed rooms = %v, want %v", result.TraversedRooms, wantTraversed)
	}
	if game.Arrows() != 4 {
		t.Fatalf("arrows = %d, want 4", game.Arrows())
	}
}

func TestInvalidArrowSegmentDefaultsToFirstExit(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Stay: true})

	result := game.Shoot([]int{3})

	wantTraversed := []int{2}
	if !reflect.DeepEqual(result.TraversedRooms, wantTraversed) {
		t.Fatalf("traversed rooms = %v, want %v", result.TraversedRooms, wantTraversed)
	}
}

func TestArrowCanHitPlayer(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)

	result := game.Shoot([]int{2, 1})

	wantMessages := []string{"OUCH! ARROW GOT YOU!", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", result.Messages, wantMessages)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
	if game.Arrows() != 4 {
		t.Fatalf("arrows = %d, want 4", game.Arrows())
	}
}

func TestShootingAfterGameOverDoesNothing(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)
	game.Shoot([]int{2})

	result := game.Shoot([]int{5})

	if result.RejectedMessage != "" || len(result.Messages) != 0 || len(result.TraversedRooms) != 0 {
		t.Fatalf("shoot after game over result = %#v, want empty result", result)
	}
	if game.Arrows() != 4 {
		t.Fatalf("arrows = %d, want 4", game.Arrows())
	}
}

func TestMissedArrowWakesWumpus(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Destination: 1})

	result := game.Shoot([]int{5})

	wantMessages := []string{"MISSED", "TSK TSK TSK - WUMPUS GOT YOU!", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", result.Messages, wantMessages)
	}
	if game.Setup().Wumpus != 1 {
		t.Fatalf("Wumpus room = %d, want 1", game.Setup().Wumpus)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestMissingWithLastArrowLoses(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(1)
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Stay: true})

	result := game.Shoot([]int{5})

	wantMessages := []string{"MISSED", "YOU RAN OUT OF ARROWS", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", result.Messages, wantMessages)
	}
	if game.Arrows() != 0 {
		t.Fatalf("arrows = %d, want 0", game.Arrows())
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestShootingPathMustContainOneToFiveRooms(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetArrows(5)

	result := game.Shoot(nil)

	if result.RejectedMessage != "CAN'T SHOOT THERE" {
		t.Fatalf("rejection = %q, want CAN'T SHOOT THERE", result.RejectedMessage)
	}
	if game.Arrows() != 5 {
		t.Fatalf("arrows = %d, want 5", game.Arrows())
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
}
