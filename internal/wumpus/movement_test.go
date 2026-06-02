package wumpus

import (
	"reflect"
	"testing"
)

func TestLegalMoveRelocatesPlayer(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})

	result := game.Move(2)

	if result.RejectedMessage != "" {
		t.Fatalf("move was rejected: %s", result.RejectedMessage)
	}
	if game.Setup().Player != 2 {
		t.Fatalf("player room = %d, want 2", game.Setup().Player)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
}

func TestIllegalMoveIsRejectedWithoutRelocating(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 19, Pits: []int{13, 14}, Bats: []int{16, 17}})

	result := game.Move(20)

	if result.RejectedMessage != "CAN'T MOVE THERE" {
		t.Fatalf("rejection = %q, want CAN'T MOVE THERE", result.RejectedMessage)
	}
	if game.Setup().Player != 1 {
		t.Fatalf("player room = %d, want 1", game.Setup().Player)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
}

func TestMoveAfterLossDoesNothing(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{2, 14}, Bats: []int{16, 17}})
	game.Move(2)

	result := game.Move(3)

	if result.RejectedMessage != "" || len(result.Messages) != 0 {
		t.Fatalf("move after loss result = %#v, want empty result", result)
	}
	if game.Setup().Player != 2 {
		t.Fatalf("player room = %d, want 2", game.Setup().Player)
	}
}

func TestMoveIntoPitLosesImmediately(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{2, 14}, Bats: []int{16, 17}})

	result := game.Move(2)

	want := []string{"YYYIIIIEEEE . . . FELL IN PIT", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestMoveIntoBatsRelocatesAndResolvesDestinationHazards(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{2, 17}})
	game.SetNextBatRelocation(13)

	result := game.Move(2)

	want := []string{
		"ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!",
		"YYYIIIIEEEE . . . FELL IN PIT",
		"HA HA HA - YOU LOSE!",
	}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Setup().Player != 13 {
		t.Fatalf("player room = %d, want 13", game.Setup().Player)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestBatRelocationDefaultsToRoomOne(t *testing.T) {
	game := mustGame(t, Setup{Player: 3, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{2, 17}})

	game.Move(2)

	if game.Setup().Player != 1 {
		t.Fatalf("player room = %d, want default bat relocation to room 1", game.Setup().Player)
	}
}

func TestMoveIntoWumpusRoomWakesWumpus(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Destination: 3})

	result := game.Move(2)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	if game.Setup().Wumpus != 3 {
		t.Fatalf("Wumpus room = %d, want 3", game.Setup().Wumpus)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
}

func TestDefaultWumpusWakeChoiceStays(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})

	game.Move(2)

	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestWumpusWakeStayingWithPlayerLoses(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Stay: true})

	result := game.Move(2)

	want := []string{"TSK TSK TSK - WUMPUS GOT YOU!", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}

func TestTurnWarningsAreEmptyWhenNoHazardsNearby(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})

	if warnings := game.TurnWarnings(); len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none", warnings)
	}
}

func TestTurnWarningsUseOriginalWarningOrder(t *testing.T) {
	game := mustGame(t, Setup{Player: 6, Wumpus: 5, Pits: []int{7, 15}, Bats: []int{1, 2}})

	got := game.TurnWarnings()
	want := []string{"I SMELL A WUMPUS", "BATS NEARBY", "I FEEL A DRAFT"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("warnings = %v, want %v", got, want)
	}
}

func mustGame(t *testing.T, setup Setup) *Game {
	t.Helper()
	game, err := NewGameWithSetup(setup)
	if err != nil {
		t.Fatal(err)
	}
	return &game
}
