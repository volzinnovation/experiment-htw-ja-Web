package wumpus

import "math/rand"

const eventSeedOffset int64 = 0x5eed5eed

func newSetupRandom(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func newEventRandom(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed + eventSeedOffset))
}

func copyRandom(random *rand.Rand) *rand.Rand {
	if random == nil {
		return nil
	}
	copy := *random
	return &copy
}

func (g *Game) eventIntn(limit int) (int, bool) {
	if g.eventRandom == nil {
		return 0, false
	}
	return g.eventRandom.Intn(limit), true
}

func (g *Game) randomRoom() int {
	if room, ok := g.eventIntn(20); ok {
		return room + 1
	}
	return 1
}
