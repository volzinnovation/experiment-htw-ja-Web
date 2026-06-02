package wumpus

import (
	"reflect"
	"testing"
)

func TestAdjacentSleepyObservationCanPutWumpusToSleep(t *testing.T) {
	game := mustGame(t, Setup{Player: 6, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetNextSleepyWumpusObservation(true)

	result := game.Move(5)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	if !game.WumpusAsleep() {
		t.Fatal("Wumpus should be asleep")
	}
	wantWarnings := []string{"I SMELL A WUMPUS", "YOU HEAR HORRIBLE SNORING"}
	if !reflect.DeepEqual(game.TurnWarnings(), wantWarnings) {
		t.Fatalf("warnings = %v, want %v", game.TurnWarnings(), wantWarnings)
	}
}

func TestSnoringWarningDoesNotIncludeDistantBats(t *testing.T) {
	game := mustGame(t, Setup{Player: 3, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{2, 17}})
	game.SetPlayerRoom(5)
	game.SetWumpusAsleep(true)

	want := []string{"I SMELL A WUMPUS", "YOU HEAR HORRIBLE SNORING"}
	if !reflect.DeepEqual(game.TurnWarnings(), want) {
		t.Fatalf("warnings = %v, want %v", game.TurnWarnings(), want)
	}
}

func TestMovingAwayFromSleepingWumpusAdjacencyAwakensIt(t *testing.T) {
	game := mustGame(t, Setup{Player: 5, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)

	result := game.Move(6)

	requireWumpusAwake(t, game)
	want := []string{"YOU HEAR A SNORT AND \"HUH?\""}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
}

func TestMovingAwayFromAwakeWumpusAdjacencyDoesNotWakeIt(t *testing.T) {
	game := mustGame(t, Setup{Player: 5, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{16, 17}})

	result := game.Move(6)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	requireWumpusAwake(t, game)
}

func TestNonAdjacentMoveDoesNotWakeSleepingWumpus(t *testing.T) {
	game := mustGame(t, Setup{Player: 6, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)

	result := game.Move(7)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	if !game.WumpusAsleep() {
		t.Fatal("Wumpus should remain asleep")
	}
}

func TestEnteringSleepingWumpusRoomCanLeaveItAsleep(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)
	game.SetNextSleepingWumpusEntryOutcome(SleepingWumpusStaysAsleep)

	result := game.Move(2)

	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
	if !game.WumpusAsleep() {
		t.Fatal("Wumpus should remain asleep")
	}
	want := []string{"YOU SEE THE HUDDLED HORRIBLE SHAPE OF THE SLEEPING WUMPUS"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}

	leaving := game.Move(1)

	requireWumpusAwake(t, game)
	wantLeaving := []string{"YOU HEAR A PETULANT SCREAM!"}
	if !reflect.DeepEqual(leaving.Messages, wantLeaving) {
		t.Fatalf("leaving messages = %v, want %v", leaving.Messages, wantLeaving)
	}
}

func TestEnteringSleepingWumpusRoomCanWakeAndLose(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)
	game.SetNextSleepingWumpusEntryOutcome(SleepingWumpusWakes)

	result := game.Move(2)

	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
	requireWumpusAwake(t, game)
	want := []string{"YOU HEAR THE WUMPUS SAY \"YUMMY BREAKFAST!\"", "HA HA HA - YOU LOSE!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
}

func TestLeavingAwakeSeenWumpusRoomDoesNotScream(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetPlayerRoom(2)
	game.SetSawSleepingWumpus(true)

	result := game.Move(1)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	requireWumpusAwake(t, game)
}

func TestLeavingUnseenSleepingWumpusRoomDoesNotScream(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetPlayerRoom(2)
	game.SetWumpusAsleep(true)
	game.SetNextSleepyWumpusObservation(true)

	result := game.Move(1)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	if !game.WumpusAsleep() {
		t.Fatal("Wumpus should remain asleep")
	}
}

func TestSawSleepingWumpusElsewhereDoesNotScream(t *testing.T) {
	game := mustGame(t, Setup{Player: 6, Wumpus: 1, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)
	game.SetSawSleepingWumpus(true)

	result := game.Move(7)

	if len(result.Messages) != 0 {
		t.Fatalf("messages = %v, want none", result.Messages)
	}
	if !game.WumpusAsleep() {
		t.Fatal("Wumpus should remain asleep")
	}
}

func TestLeavingAfterSeeingSleepingWumpusAwakensIt(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetPlayerRoom(2)
	game.SetWumpusAsleep(true)
	game.SetSawSleepingWumpus(true)

	result := game.Move(1)

	requireWumpusAwake(t, game)
	want := []string{"YOU HEAR A PETULANT SCREAM!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
}

func TestMissedArrowWakesSleepingWumpus(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetWumpusAsleep(true)
	game.SetNextWumpusWakeChoice(WumpusWakeChoice{Destination: 11})

	result := game.Shoot([]int{5})

	requireWumpusAwake(t, game)
	if game.Setup().Wumpus != 11 {
		t.Fatalf("Wumpus room = %d, want 11", game.Setup().Wumpus)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
	if !reflect.DeepEqual(result.Messages, []string{"MISSED"}) {
		t.Fatalf("messages = %v", result.Messages)
	}
}

func TestObserveSleepyWumpusBehaviorIsDeterministic(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 3, Pits: []int{13, 14}, Bats: []int{16, 17}})

	got := game.ObserveSleepyWumpusBehavior(4)

	want := []string{"asleep", "awake", "asleep", "awake"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("observations = %v, want %v", got, want)
	}
}

func requireWumpusAwake(t *testing.T, game *Game) {
	t.Helper()
	if game.WumpusAsleep() {
		t.Fatal("Wumpus should be awake")
	}
}
