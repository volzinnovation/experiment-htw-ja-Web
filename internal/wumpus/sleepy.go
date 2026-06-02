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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:38:31-05:00","module_hash":"1eabf7fc52cc646d8d671c2fc8c712023724bc056541083b078e223075a21e8b","functions":[{"id":"func/Game.SetWumpusAsleep","name":"Game.SetWumpusAsleep","line":10,"end_line":12,"hash":"8281acb508406ccd2b2018337d07117cf9dbe31534fa2bd25beb152c00f63352"},{"id":"func/Game.WumpusAsleep","name":"Game.WumpusAsleep","line":14,"end_line":16,"hash":"97fc9a4214ce6bef013a8cd63e17ebbc3e040bbb604214b5af17de4019d84f23"},{"id":"func/Game.SetNextSleepyWumpusObservation","name":"Game.SetNextSleepyWumpusObservation","line":18,"end_line":20,"hash":"44b68532e3cec1e2febafc2d780425a522b545875803963b1bbfd2f8cadb484b"},{"id":"func/Game.SetNextSleepingWumpusEntryOutcome","name":"Game.SetNextSleepingWumpusEntryOutcome","line":22,"end_line":24,"hash":"dac6941c0627854b467a76debc246e6fd6dcc7d4885f720f4d6b142fbc4b30db"},{"id":"func/Game.SetSawSleepingWumpus","name":"Game.SetSawSleepingWumpus","line":26,"end_line":28,"hash":"c9f41a6efd9f050fd55f8bdbbcdd3b1271de7ae2fbe8e35986994f117e49b781"},{"id":"func/Game.SetPlayerRoom","name":"Game.SetPlayerRoom","line":30,"end_line":32,"hash":"d3fd868a383b8a40bc4b66c1d65d051c5a5600b74ae5eef3b27c5fcc5ec4d409"},{"id":"func/Game.ObserveSleepyWumpusBehavior","name":"Game.ObserveSleepyWumpusBehavior","line":34,"end_line":36,"hash":"5e5d5486f997642a27215eb283006b01c429ff20268236eb43c4ae89a6ebfe54"},{"id":"func/Game.observeAdjacentWumpus","name":"Game.observeAdjacentWumpus","line":38,"end_line":43,"hash":"372c3cab1345c01e6af604742c9b4334d67828da548fb2893056c356056d1f82"},{"id":"func/Game.nextSleepyObservation","name":"Game.nextSleepyObservation","line":45,"end_line":47,"hash":"8922f522addc0229fa2fb5395161be9429941f56e21da385edb96a404384792e"},{"id":"func/Game.resolveSleepingWumpusEntry","name":"Game.resolveSleepingWumpusEntry","line":49,"end_line":56,"hash":"d8f51c1e29b7be1303097dc322ebc7408fe33150ef307351f4b6d72b065ba831"},{"id":"func/Game.nextSleepingEntryOutcome","name":"Game.nextSleepingEntryOutcome","line":58,"end_line":60,"hash":"0830d1559612d42bd2a8090ddfde4b148d78290ff6461330f71e8499e26c0d20"}]}
// mutate4go-manifest-end
