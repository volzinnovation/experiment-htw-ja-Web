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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T08:22:09-05:00","module_hash":"b6eb6c1fcab7cfea3d19c40319b1098dbadc4e69cb906f51b79248af7e19daf2","functions":[{"id":"func/Game.Move","name":"Game.Move","line":23,"end_line":32,"hash":"1f511b984b7320f4fbe80f2513a5dafe33eb8c5c251bc781084033a56c9d998c"},{"id":"func/Game.Status","name":"Game.Status","line":34,"end_line":36,"hash":"afed41261db0e9ac5dde549d3958f58f70cd99cd0cbd616008656311e8600d56"},{"id":"func/Game.SetNextBatRelocation","name":"Game.SetNextBatRelocation","line":38,"end_line":40,"hash":"3f4687bec66a92d1e67f433e2891eed0b02b1fc48bc899a73100460799a3873a"},{"id":"func/Game.SetNextWumpusWakeChoice","name":"Game.SetNextWumpusWakeChoice","line":42,"end_line":44,"hash":"dd392a58aefa58b7e19172145fb428645934e1b49ed65d1215540cfb125c77ec"},{"id":"func/Game.TurnWarnings","name":"Game.TurnWarnings","line":46,"end_line":66,"hash":"a7e15abf4f9f10dfe5f0afb7e44e3d59507b53f730c5dc3814fd87a830b353b5"},{"id":"func/Game.batsNearbyForWarning","name":"Game.batsNearbyForWarning","line":68,"end_line":81,"hash":"861c0a7cac7a1e6dc72bc7360b632d6684a4166f018a05eeaffeb8a04bd919c9"},{"id":"func/Game.mustExits","name":"Game.mustExits","line":83,"end_line":89,"hash":"3f39ccd8c60d96e65fe68e76e40e487879d84be24c38a4c68fa026cda1151460"},{"id":"func/Game.resolveArrival","name":"Game.resolveArrival","line":91,"end_line":102,"hash":"5415dbacf4f10f609da7db19bd59e4d470b033e80e5d15e7d45b0e97b7e7c33e"},{"id":"func/Game.resolveBatArrival","name":"Game.resolveBatArrival","line":104,"end_line":108,"hash":"4afa263b6db8e474fb2a69bc2dfd05574529d27dfac01ba9f3b9f2dec2de7b43"},{"id":"func/Game.nextBatRoom","name":"Game.nextBatRoom","line":110,"end_line":117,"hash":"3403bc095af58c5538454dbdae2c204d998870f11b8710f0131273a138a5ed05"},{"id":"func/Game.wakeWumpus","name":"Game.wakeWumpus","line":119,"end_line":128,"hash":"63de59dc2c61478da55b2bbe851084bb1ee0fe3c4026b012f8b2a3542cb299af"},{"id":"func/Game.nextWumpusChoice","name":"Game.nextWumpusChoice","line":130,"end_line":137,"hash":"1246e6c314697c625651e7811569057a80555a4ffdf14bf3a06c9706c8c5a745"},{"id":"func/Game.lose","name":"Game.lose","line":139,"end_line":142,"hash":"76351950526b8a93d43d0a1cce7cb580a4c37c99f5e04dcc31cd32c133f75530"}]}
// mutate4go-manifest-end
