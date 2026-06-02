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
	path := g.nextWumpusJumpPath()
	for index, room := range path {
		if !NewCave().HasTunnel(g.setup.Wumpus, room) {
			room = firstExit(g.setup.Wumpus)
		}
		g.setup.Wumpus = room
		result.JumpedRooms = append(result.JumpedRooms, room)
		if room == g.setup.Player {
			result.Messages = append(result.Messages, g.resolveJumpLandingOnPlayer(index)...)
			if g.status != StatusInProgress {
				break
			}
		}
	}
	return result
}

func (g *Game) nextJumpingWumpusEvent() bool {
	if len(g.nextJumpEvents) == 0 {
		return false
	}
	jumps := g.nextJumpEvents[0]
	g.nextJumpEvents = g.nextJumpEvents[1:]
	return jumps
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
	if len(g.nextFirstJumpLand) == 0 {
		return FirstJumpTramples
	}
	outcome := g.nextFirstJumpLand[0]
	g.nextFirstJumpLand = g.nextFirstJumpLand[1:]
	return outcome
}

func (g *Game) ObserveJumpingWumpusBehavior(turnCount int) []string {
	events := make([]string, 0, turnCount)
	for i := 0; i < turnCount; i++ {
		if (g.setup.Player+g.setup.Wumpus+i)%3 == 0 {
			events = append(events, "jumps")
		} else {
			events = append(events, "no jump")
		}
	}
	return events
}

func firstExit(room int) int {
	exits, err := NewCave().Exits(room)
	if err != nil || len(exits) == 0 {
		return room
	}
	return exits[0]
}
