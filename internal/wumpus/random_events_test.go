package wumpus

import (
	"math/rand"
	"testing"
)

func TestRandomBatRelocationCanUseSeededRNG(t *testing.T) {
	room, found := firstRandomBatRoomNot(1)

	if !found {
		t.Fatal("no seeded random bat relocation differed from room 1")
	}
	if room < 1 || room > 20 {
		t.Fatalf("bat relocation room = %d, want 1..20", room)
	}
}

func TestRandomWumpusWakeCanMoveToAdjacentRoom(t *testing.T) {
	game, choice, found := firstRandomWumpusWakeMove(t)

	if !found {
		t.Fatal("no seeded random Wumpus wake moved")
	}
	if choice.Stay {
		t.Fatal("wake choice stayed, want move")
	}
	if !NewCave().HasTunnel(game.Setup().Wumpus, choice.Destination) {
		t.Fatalf("Wumpus wake destination %d is not adjacent to %d", choice.Destination, game.Setup().Wumpus)
	}
}

func TestRandomArrowDeviationCanUseSeededRNG(t *testing.T) {
	room, found := firstRandomArrowDeviationNot(t, 1, 2)

	if !found {
		t.Fatal("no seeded random arrow deviation differed from first exit")
	}
	if !NewCave().HasTunnel(1, room) {
		t.Fatalf("arrow deviation room = %d, want tunnel from room 1", room)
	}
}

func TestRandomSleepyObservationCanPutWumpusToSleep(t *testing.T) {
	found := false
	for seed := int64(1); seed <= 100; seed++ {
		game := randomEventGame(t, seed)
		if game.nextSleepyObservation() {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("no seeded random sleepy observation put Wumpus to sleep")
	}
}

func TestRandomSleepingEntryCanLeaveWumpusAsleep(t *testing.T) {
	found := false
	for seed := int64(1); seed <= 100; seed++ {
		game := randomEventGame(t, seed)
		if game.nextSleepingEntryOutcome() == SleepingWumpusStaysAsleep {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("no seeded random sleeping entry left Wumpus asleep")
	}
}

func TestRandomJumpingWumpusCanTriggerWhump(t *testing.T) {
	game := randomEventGame(t, 1)
	found := false
	for turn := 0; turn < 400 && game.Status() == StatusInProgress; turn++ {
		result := game.ResolveJumpingWumpusTurn()
		if len(result.Messages) > 0 && result.Messages[0] == "YOU HEAR WHUMP, WHUMP." {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("seeded random jumping Wumpus did not trigger")
	}
}

func firstRandomBatRoomNot(defaultRoom int) (int, bool) {
	for seed := int64(1); seed <= 100; seed++ {
		game, _ := NewGameWithSetup(Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{2, 17}})
		game.eventRandom = rand.New(rand.NewSource(seed))
		if room := game.nextBatRoom(); room != defaultRoom {
			return room, true
		}
	}
	return 0, false
}

func firstRandomWumpusWakeMove(t *testing.T) (*Game, WumpusWakeChoice, bool) {
	t.Helper()
	for seed := int64(1); seed <= 100; seed++ {
		game := randomEventGame(t, seed)
		if choice := game.nextWumpusChoice(); !choice.Stay {
			return game, choice, true
		}
	}
	return nil, WumpusWakeChoice{}, false
}

func firstRandomArrowDeviationNot(t *testing.T, from, defaultRoom int) (int, bool) {
	t.Helper()
	for seed := int64(1); seed <= 100; seed++ {
		game := randomEventGame(t, seed)
		if room := game.nextArrowDeviationRoom(from); room != defaultRoom {
			return room, true
		}
	}
	return 0, false
}

func randomEventGame(t *testing.T, seed int64) *Game {
	t.Helper()
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.eventRandom = rand.New(rand.NewSource(seed))
	return game
}
