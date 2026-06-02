package wumpus

import "slices"

type Status string

const (
	StatusInProgress Status = "in progress"
	StatusLost       Status = "lost"
	StatusWon        Status = "won"
)

type MoveResult struct {
	RejectedMessage string
	Messages        []string
}

type WumpusWakeChoice struct {
	Stay        bool
	Destination int
}

func (g *Game) Move(to int) MoveResult {
	if g.status != StatusInProgress {
		return MoveResult{}
	}
	if !NewCave().HasTunnel(g.setup.Player, to) {
		return MoveResult{RejectedMessage: "CAN'T MOVE THERE"}
	}
	from := g.setup.Player
	g.setup.Player = to
	messages := g.sleepWakeMessagesOnMove(from, to)
	messages = append(messages, g.resolveArrival()...)
	if g.status == StatusInProgress && g.setup.Player != g.setup.Wumpus {
		g.observeAdjacentWumpus()
	}
	return MoveResult{Messages: messages}
}

func (g Game) Status() Status {
	return g.status
}

func (g *Game) SetNextBatRelocation(room int) {
	g.nextBatRelocation = append(g.nextBatRelocation, room)
}

func (g *Game) SetNextWumpusWakeChoice(choice WumpusWakeChoice) {
	g.nextWumpusWake = append(g.nextWumpusWake, choice)
}

func (g Game) TurnWarnings() []string {
	present := g.warningHazards()
	return warningMessages(present, g.wumpusAsleep)
}

func (g Game) warningHazards() map[Hazard]bool {
	hazards := NewCave().AdjacentHazards(g.setup.Player, g.setup)
	present := map[Hazard]bool{}
	for _, hazard := range hazards {
		present[hazard] = true
	}
	return present
}

func warningMessages(present map[Hazard]bool, wumpusAsleep bool) []string {
	var warnings []string
	if present[HazardWumpus] {
		warnings = append(warnings, "I SMELL A WUMPUS")
		if wumpusAsleep {
			warnings = append(warnings, "YOU HEAR HORRIBLE SNORING")
		}
	}
	if present[HazardBats] {
		warnings = append(warnings, "BATS NEARBY")
	}
	if present[HazardPit] {
		warnings = append(warnings, "I FEEL A DRAFT")
	}
	return warnings
}

func (g *Game) resolveArrival() []string {
	switch {
	case slices.Contains(g.setup.Pits, g.setup.Player):
		if g.inertHazards {
			return []string{"QA INERT: PIT IGNORED"}
		}
		return g.lose("YYYIIIIEEEE . . . FELL IN PIT")
	case slices.Contains(g.setup.Bats, g.setup.Player):
		if g.inertHazards {
			return []string{"QA INERT: BATS IGNORED"}
		}
		return g.resolveBatArrival()
	case g.setup.Wumpus == g.setup.Player:
		if g.inertHazards {
			return []string{"QA INERT: WUMPUS IGNORED"}
		}
		if g.wumpusAsleep {
			return g.resolveSleepingWumpusEntry()
		}
		return g.wakeWumpus()
	default:
		return g.collectGrenadeIfPresent()
	}
}

func (g *Game) resolveBatArrival() []string {
	messages := []string{"ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!"}
	g.setup.Player = g.nextBatRoom()
	return append(messages, g.resolveArrival()...)
}

func (g *Game) nextBatRoom() int {
	if len(g.nextBatRelocation) == 0 && g.random != nil {
		return g.randomRoom()
	}
	return dequeueOr(&g.nextBatRelocation, 1)
}

func (g *Game) wakeWumpus() []string {
	g.wumpusAsleep = false
	choice := g.nextWumpusChoice()
	if !choice.Stay {
		g.setup.Wumpus = choice.Destination
	}
	if g.setup.Wumpus == g.setup.Player {
		return g.lose("TSK TSK TSK - WUMPUS GOT YOU!")
	}
	return nil
}

func (g *Game) nextWumpusChoice() WumpusWakeChoice {
	if len(g.nextWumpusWake) == 0 && g.random != nil {
		if g.random.Intn(4) == 0 {
			return WumpusWakeChoice{Stay: true}
		}
		return WumpusWakeChoice{Destination: g.randomExit(g.setup.Wumpus)}
	}
	if len(g.nextWumpusWake) == 0 {
		return WumpusWakeChoice{Stay: true}
	}
	choice := g.nextWumpusWake[0]
	g.nextWumpusWake = g.nextWumpusWake[1:]
	return choice
}

func (g *Game) lose(reason string) []string {
	g.status = StatusLost
	return []string{reason, "HA HA HA - YOU LOSE!"}
}

func (g *Game) sleepWakeMessagesOnMove(from, to int) []string {
	if message, woke := g.wakeAfterLeavingSeenWumpus(from, to); woke {
		return []string{message}
	}
	if message, woke := g.wakeAfterLeavingWumpusAdjacency(from, to); woke {
		return []string{message}
	}
	return nil
}

func (g *Game) wakeAfterLeavingSeenWumpus(from, to int) (string, bool) {
	if !g.wumpusAsleep || !g.sawSleepingWumpus || from != g.setup.Wumpus || to == g.setup.Wumpus {
		return "", false
	}
	g.wumpusAsleep = false
	g.sawSleepingWumpus = false
	return "YOU HEAR A PETULANT SCREAM!", true
}

func (g *Game) wakeAfterLeavingWumpusAdjacency(from, to int) (string, bool) {
	cave := NewCave()
	if !g.wumpusAsleep || !cave.HasTunnel(from, g.setup.Wumpus) || to == g.setup.Wumpus || cave.HasTunnel(to, g.setup.Wumpus) {
		return "", false
	}
	g.wumpusAsleep = false
	return "YOU HEAR A SNORT AND \"HUH?\"", true
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:53:22-05:00","module_hash":"fa0136080ffa420258eef0a4d358703bf6d3406b4002c3f44d9958d86d131c12","functions":[{"id":"func/Game.Move","name":"Game.Move","line":23,"end_line":38,"hash":"f5c1e119fc3fb4979823b5fffb3ed58fb30d9b4002fb1b7835c8ee3ecdf6718b"},{"id":"func/Game.Status","name":"Game.Status","line":40,"end_line":42,"hash":"afed41261db0e9ac5dde549d3958f58f70cd99cd0cbd616008656311e8600d56"},{"id":"func/Game.SetNextBatRelocation","name":"Game.SetNextBatRelocation","line":44,"end_line":46,"hash":"3f4687bec66a92d1e67f433e2891eed0b02b1fc48bc899a73100460799a3873a"},{"id":"func/Game.SetNextWumpusWakeChoice","name":"Game.SetNextWumpusWakeChoice","line":48,"end_line":50,"hash":"dd392a58aefa58b7e19172145fb428645934e1b49ed65d1215540cfb125c77ec"},{"id":"func/Game.TurnWarnings","name":"Game.TurnWarnings","line":52,"end_line":55,"hash":"58209c4d47f68769d50d9ed5d355b9c2e1714eaeb9faf0bfe640d5a969f24a7c"},{"id":"func/Game.warningHazards","name":"Game.warningHazards","line":57,"end_line":67,"hash":"499822be32553ef221e834c2e8134e446944782a8f6463ec2aeb6f2b1220373c"},{"id":"func/warningMessages","name":"warningMessages","line":69,"end_line":84,"hash":"7c2857d584b6d642e5ee5fbe0b9c0b303e50cf8beb7e408ab308e04397887108"},{"id":"func/Game.batsNearbyForWarning","name":"Game.batsNearbyForWarning","line":86,"end_line":99,"hash":"861c0a7cac7a1e6dc72bc7360b632d6684a4166f018a05eeaffeb8a04bd919c9"},{"id":"func/Game.mustExits","name":"Game.mustExits","line":101,"end_line":107,"hash":"3f39ccd8c60d96e65fe68e76e40e487879d84be24c38a4c68fa026cda1151460"},{"id":"func/Game.resolveArrival","name":"Game.resolveArrival","line":109,"end_line":123,"hash":"b058de6d9902fbe7d6332b8c982b6fedd3d0cd46ba73c6314cc340e8a9ab0fa5"},{"id":"func/Game.resolveBatArrival","name":"Game.resolveBatArrival","line":125,"end_line":129,"hash":"4afa263b6db8e474fb2a69bc2dfd05574529d27dfac01ba9f3b9f2dec2de7b43"},{"id":"func/Game.nextBatRoom","name":"Game.nextBatRoom","line":131,"end_line":138,"hash":"3403bc095af58c5538454dbdae2c204d998870f11b8710f0131273a138a5ed05"},{"id":"func/Game.wakeWumpus","name":"Game.wakeWumpus","line":140,"end_line":150,"hash":"6cd7c98c81833a7c24a22c7d20abddb0e7e24c35408c8828c5db2d313336645e"},{"id":"func/Game.nextWumpusChoice","name":"Game.nextWumpusChoice","line":152,"end_line":159,"hash":"1246e6c314697c625651e7811569057a80555a4ffdf14bf3a06c9706c8c5a745"},{"id":"func/Game.lose","name":"Game.lose","line":161,"end_line":164,"hash":"76351950526b8a93d43d0a1cce7cb580a4c37c99f5e04dcc31cd32c133f75530"},{"id":"func/Game.sleepWakeMessagesOnMove","name":"Game.sleepWakeMessagesOnMove","line":166,"end_line":174,"hash":"1f241ade99d536383c9acf3aabd30f439b0b2fe7d52f0dc9091998ef6e8b4bb6"},{"id":"func/Game.wakeAfterLeavingSeenWumpus","name":"Game.wakeAfterLeavingSeenWumpus","line":176,"end_line":183,"hash":"043b6dfd2af7d8245af2125d03ab6b32cd5d3092e8105ef45202a8b1fbf24b51"},{"id":"func/Game.wakeAfterLeavingWumpusAdjacency","name":"Game.wakeAfterLeavingWumpusAdjacency","line":185,"end_line":192,"hash":"d55dafb8742efb8ee2206f1aa2657e1f685cb8d87674392928f5e763f9d25dcd"}]}
// mutate4go-manifest-end
