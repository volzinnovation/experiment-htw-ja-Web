package wumpus

type ShootResult struct {
	RejectedMessage string
	Messages        []string
	TraversedRooms  []int
}

func (g *Game) SetArrows(arrows int) {
	g.arrows = arrows
}

func (g Game) Arrows() int {
	return g.arrows
}

func (g *Game) SetNextArrowDeviation(room int) {
	g.nextArrowDeviation = append(g.nextArrowDeviation, room)
}

func (g *Game) Shoot(path []int) ShootResult {
	if g.status != StatusInProgress {
		return ShootResult{}
	}
	if len(path) < 1 || len(path) > 5 {
		return ShootResult{RejectedMessage: "CAN'T SHOOT THERE"}
	}
	g.arrows--
	return g.followArrow(path)
}

func (g *Game) followArrow(path []int) ShootResult {
	arrowRoom := g.setup.Player
	var traversed []int
	for _, requestedRoom := range path {
		arrowRoom = g.nextArrowRoom(arrowRoom, requestedRoom)
		traversed = append(traversed, arrowRoom)
		if arrowRoom == g.setup.Wumpus {
			g.status = StatusWon
			return ShootResult{Messages: []string{"AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"}, TraversedRooms: traversed}
		}
		if arrowRoom == g.setup.Player {
			return ShootResult{Messages: g.lose("OUCH! ARROW GOT YOU!"), TraversedRooms: traversed}
		}
	}
	return g.missedArrowResult(traversed)
}

func (g *Game) nextArrowRoom(from, requested int) int {
	if NewCave().HasTunnel(from, requested) {
		return requested
	}
	return g.nextArrowDeviationRoom(from)
}

func (g *Game) nextArrowDeviationRoom(from int) int {
	if len(g.nextArrowDeviation) > 0 {
		room := g.nextArrowDeviation[0]
		g.nextArrowDeviation = g.nextArrowDeviation[1:]
		return room
	}
	exits, err := NewCave().Exits(from)
	if err != nil || len(exits) == 0 {
		return from
	}
	return exits[0]
}

func (g *Game) missedArrowResult(traversed []int) ShootResult {
	messages := []string{"MISSED"}
	messages = append(messages, g.wakeWumpus()...)
	if g.status == StatusInProgress && g.arrows == 0 {
		messages = append(messages, g.lose("YOU RAN OUT OF ARROWS")...)
	}
	return ShootResult{Messages: messages, TraversedRooms: traversed}
}
