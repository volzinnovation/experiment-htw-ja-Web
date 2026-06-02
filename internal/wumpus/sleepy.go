package wumpus

type SleepingWumpusEntryOutcome string

const (
	SleepingWumpusWakes       SleepingWumpusEntryOutcome = "wakes"
	SleepingWumpusStaysAsleep SleepingWumpusEntryOutcome = "stays asleep"
)

func (g *Game) SetWumpusAsleep(asleep bool) {
	g.wumpusAsleep = asleep
}

func (g Game) WumpusAsleep() bool {
	return g.wumpusAsleep
}

func (g *Game) SetNextSleepyWumpusObservation(asleep bool) {
	g.nextSleepyObserve = append(g.nextSleepyObserve, asleep)
}

func (g *Game) SetNextSleepingWumpusEntryOutcome(outcome SleepingWumpusEntryOutcome) {
	g.nextSleepingEntry = append(g.nextSleepingEntry, outcome)
}

func (g *Game) SetSawSleepingWumpus(seen bool) {
	g.sawSleepingWumpus = seen
}

func (g *Game) SetPlayerRoom(room int) {
	g.setup.Player = room
}

func (g *Game) ObserveSleepyWumpusBehavior(turnCount int) []string {
	return periodicBehavior(g.setup.Player, g.setup.Wumpus, turnCount, 2, "asleep", "awake")
}

func (g *Game) observeAdjacentWumpus() {
	if !NewCave().HasTunnel(g.setup.Player, g.setup.Wumpus) {
		return
	}
	g.wumpusAsleep = g.nextSleepyObservation()
}

func (g *Game) nextSleepyObservation() bool {
	return dequeueOr(&g.nextSleepyObserve, false)
}

func (g *Game) resolveSleepingWumpusEntry() []string {
	if g.nextSleepingEntryOutcome() == SleepingWumpusStaysAsleep {
		g.sawSleepingWumpus = true
		return []string{"YOU SEE THE HUDDLED HORRIBLE SHAPE OF THE SLEEPING WUMPUS"}
	}
	g.wumpusAsleep = false
	return g.lose("YOU HEAR THE WUMPUS SAY \"YUMMY BREAKFAST!\"")
}

func (g *Game) nextSleepingEntryOutcome() SleepingWumpusEntryOutcome {
	return dequeueOr(&g.nextSleepingEntry, SleepingWumpusWakes)
}
