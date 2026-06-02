package wumpus

import "slices"

type Status string

const (
	StatusInProgress Status = "in progress"
	StatusLost       Status = "lost"
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
	g.setup.Player = to
	return MoveResult{Messages: g.resolveArrival()}
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
	hazards := NewCave().AdjacentHazards(g.setup.Player, g.setup)
	present := map[Hazard]bool{}
	for _, hazard := range hazards {
		present[hazard] = true
	}
	if g.batsNearbyForWarning() {
		present[HazardBats] = true
	}
	var warnings []string
	if present[HazardWumpus] {
		warnings = append(warnings, "I SMELL A WUMPUS")
	}
	if present[HazardBats] {
		warnings = append(warnings, "BATS NEARBY")
	}
	if present[HazardPit] {
		warnings = append(warnings, "I FEEL A DRAFT")
	}
	return warnings
}

func (g Game) batsNearbyForWarning() bool {
	cave := NewCave()
	for _, batRoom := range g.setup.Bats {
		if cave.HasTunnel(g.setup.Player, batRoom) {
			return true
		}
		for _, exit := range g.mustExits(g.setup.Player) {
			if cave.HasTunnel(exit, batRoom) {
				return true
			}
		}
	}
	return false
}

func (g Game) mustExits(room int) []int {
	exits, err := NewCave().Exits(room)
	if err != nil {
		return nil
	}
	return exits
}

func (g *Game) resolveArrival() []string {
	switch {
	case slices.Contains(g.setup.Pits, g.setup.Player):
		return g.lose("YYYIIIIEEEE . . . FELL IN PIT")
	case slices.Contains(g.setup.Bats, g.setup.Player):
		return g.resolveBatArrival()
	case g.setup.Wumpus == g.setup.Player:
		return g.wakeWumpus()
	default:
		return nil
	}
}

func (g *Game) resolveBatArrival() []string {
	messages := []string{"ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!"}
	g.setup.Player = g.nextBatRoom()
	return append(messages, g.resolveArrival()...)
}

func (g *Game) nextBatRoom() int {
	if len(g.nextBatRelocation) == 0 {
		return 1
	}
	room := g.nextBatRelocation[0]
	g.nextBatRelocation = g.nextBatRelocation[1:]
	return room
}

func (g *Game) wakeWumpus() []string {
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
