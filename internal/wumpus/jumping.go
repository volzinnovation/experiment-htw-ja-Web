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
	return dequeueOr(&g.nextJumpEvents, false)
}

func (g *Game) nextWumpusJumpPath() []int {
	if len(g.nextJumpPaths) == 0 {
		first := firstExit(g.setup.Wumpus)
		return []int{first, firstExit(first)}
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
	return dequeueOr(&g.nextFirstJumpLand, FirstJumpTramples)
}

func (g *Game) ObserveJumpingWumpusBehavior(turnCount int) []string {
	return periodicBehavior(g.setup.Player, g.setup.Wumpus, turnCount, 3, "jumps", "no jump")
}

func firstExit(room int) int {
	exits, err := NewCave().Exits(room)
	if err != nil || len(exits) == 0 {
		return room
	}
	return exits[0]
}
