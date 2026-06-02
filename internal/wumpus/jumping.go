package wumpus

type FirstJumpLandingOutcome string

const (
	FirstJumpTramples FirstJumpLandingOutcome = "trample"
	FirstJumpSlams    FirstJumpLandingOutcome = "slam"
)

type JumpResult struct {
	Messages    []string
	JumpedRooms []int
}

func (g *Game) SetNextJumpingWumpusTurnEvent(jumps bool) {
	g.nextJumpEvents = append(g.nextJumpEvents, jumps)
}

func (g *Game) SetNextWumpusJumpPath(path []int) {
	g.nextJumpPaths = append(g.nextJumpPaths, append([]int(nil), path...))
}

func (g *Game) SetNextFirstJumpPlayerLandingOutcome(outcome FirstJumpLandingOutcome) {
	g.nextFirstJumpLand = append(g.nextFirstJumpLand, outcome)
}

func (g *Game) ResolveJumpingWumpusTurn() JumpResult {
	if g.status != StatusInProgress || !g.nextJumpingWumpusEvent() {
		return JumpResult{}
	}
	result := JumpResult{Messages: []string{"YOU HEAR WHUMP, WHUMP."}}
	for index, room := range g.nextWumpusJumpPath() {
		if !g.resolveWumpusJump(index, room, &result) {
			break
		}
	}
	return result
}

func (g *Game) resolveWumpusJump(index, room int, result *JumpResult) bool {
	room = g.legalJumpDestination(room)
	g.setup.Wumpus = room
	result.JumpedRooms = append(result.JumpedRooms, room)
	if room != g.setup.Player {
		return true
	}
	result.Messages = append(result.Messages, g.resolveJumpLandingOnPlayer(index)...)
	return g.status == StatusInProgress
}

func (g Game) legalJumpDestination(room int) int {
	if NewCave().HasTunnel(g.setup.Wumpus, room) {
		return room
	}
	return firstExit(g.setup.Wumpus)
}

func (g *Game) nextJumpingWumpusEvent() bool {
	if len(g.nextJumpEvents) == 0 && g.random != nil {
		return g.random.Intn(20) == 0
	}
	return dequeueOr(&g.nextJumpEvents, false)
}

func (g *Game) nextWumpusJumpPath() []int {
	if len(g.nextJumpPaths) == 0 {
		first := g.randomExit(g.setup.Wumpus)
		return []int{first, g.randomExit(first)}
	}
	path := g.nextJumpPaths[0]
	g.nextJumpPaths = g.nextJumpPaths[1:]
	return append([]int(nil), path...)
}

func (g *Game) resolveJumpLandingOnPlayer(jumpIndex int) []string {
	if jumpIndex == 0 {
		if g.nextFirstJumpLandingOutcome() == FirstJumpSlams {
			return []string{"YOU ARE SLAMMED AGAINST THE CAVE WALL BY THE SNARLING WUMPUS!"}
		}
		return g.lose("THE WUMPUS TRAMPLES YOU TO DEATH!")
	}
	return []string{"YOU SEE THE BLOODSTAINED EYES OF THE WUMPUS APPRAISING YOU!"}
}

func (g *Game) nextFirstJumpLandingOutcome() FirstJumpLandingOutcome {
	if len(g.nextFirstJumpLand) == 0 && g.random != nil {
		if g.random.Intn(2) == 0 {
			return FirstJumpTramples
		}
		return FirstJumpSlams
	}
	return dequeueOr(&g.nextFirstJumpLand, FirstJumpTramples)
}

func (g *Game) ObserveJumpingWumpusBehavior(turnCount int) []string {
	return periodicBehavior(g.setup.Player, g.setup.Wumpus, turnCount, 3, "jumps", "no jump")
}

func firstExit(room int) int {
	exits, err := NewCave().Exits(room)
	if err != nil {
		return room
	}
	return exits[0]
}

func (g *Game) randomExit(room int) int {
	exits, err := NewCave().Exits(room)
	if err != nil {
		return room
	}
	if g.random == nil {
		return exits[0]
	}
	return exits[g.random.Intn(len(exits))]
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:34:52-05:00","module_hash":"6f3cd4231d9f7993b6a1ae6f3f3462e92ff2fab0cd82baf8bc805b960b9b6f01","functions":[{"id":"func/Game.SetNextJumpingWumpusTurnEvent","name":"Game.SetNextJumpingWumpusTurnEvent","line":15,"end_line":17,"hash":"321df785aa31f7e38a119f7b97a3353e710fccc490d96d408eae9166c3e83526"},{"id":"func/Game.SetNextWumpusJumpPath","name":"Game.SetNextWumpusJumpPath","line":19,"end_line":21,"hash":"038e240a13152eb3448149600231e885e8d851efac558a58d037b1fe7100eb64"},{"id":"func/Game.SetNextFirstJumpPlayerLandingOutcome","name":"Game.SetNextFirstJumpPlayerLandingOutcome","line":23,"end_line":25,"hash":"ecb4173a9fbda05dab783f5f014601a4a179048463c41c128429d9fd649793a7"},{"id":"func/Game.ResolveJumpingWumpusTurn","name":"Game.ResolveJumpingWumpusTurn","line":27,"end_line":38,"hash":"7b5fb23f6b81cd7aafdf7e2e20f0da2284dc35ba4498fa88a148a1bbcb7b7c2d"},{"id":"func/Game.resolveWumpusJump","name":"Game.resolveWumpusJump","line":40,"end_line":49,"hash":"f9d4752108033c91bb08dc0620923e3bb4c64b1c1762979e512b0f0bbb71a828"},{"id":"func/Game.legalJumpDestination","name":"Game.legalJumpDestination","line":51,"end_line":56,"hash":"abf6956b4b5c3e92b246aa87e5083c0681900163a73d18ed7437ad2caa3339a9"},{"id":"func/Game.nextJumpingWumpusEvent","name":"Game.nextJumpingWumpusEvent","line":58,"end_line":60,"hash":"7a74fa4036746cb580a299468ac06191933e0748b7679bf707c0939aa04daa9d"},{"id":"func/Game.nextWumpusJumpPath","name":"Game.nextWumpusJumpPath","line":62,"end_line":70,"hash":"40ba6177973d80830a3ff446171ca6faebf0b0cc857f35e1d59c0775a61486b3"},{"id":"func/Game.resolveJumpLandingOnPlayer","name":"Game.resolveJumpLandingOnPlayer","line":72,"end_line":80,"hash":"bbde4dfae5a61f341ae54e9b5e62f97c76a18a0797b42476db137b4ba3fa8ec1"},{"id":"func/Game.nextFirstJumpLandingOutcome","name":"Game.nextFirstJumpLandingOutcome","line":82,"end_line":84,"hash":"d77597064e40255b42296bb8727721592fd94e4b0b2ef5f4c0f7f4031b1cc81f"},{"id":"func/Game.ObserveJumpingWumpusBehavior","name":"Game.ObserveJumpingWumpusBehavior","line":86,"end_line":88,"hash":"b2965048a455fb0ce13fd239325a9e63ad50482f0fe3c3b21379d87c8c67efa5"},{"id":"func/firstExit","name":"firstExit","line":90,"end_line":96,"hash":"e68bc0f5404f3faca8a71379ce67976c74a7cd683f2ca43cc70e56ae45a521fd"}]}
// mutate4go-manifest-end
