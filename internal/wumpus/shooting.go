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
	if err != nil {
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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T08:25:38-05:00","module_hash":"a72da6656d6da8eb10f1bb9432efe9f8e48d752baf317f519282fe80ae3d3662","functions":[{"id":"func/Game.SetArrows","name":"Game.SetArrows","line":9,"end_line":11,"hash":"f1523cf1c2dcf8b6283f7ed21f73d28b2c1d0dfb01dbba500b853fca7a232ffa"},{"id":"func/Game.Arrows","name":"Game.Arrows","line":13,"end_line":15,"hash":"4c28284ed29d6b60dfe74ec8eef6c448225cc2c63a42ec54b8b4a4aa3597bf15"},{"id":"func/Game.SetNextArrowDeviation","name":"Game.SetNextArrowDeviation","line":17,"end_line":19,"hash":"0502a95e3d637ae8bcd11fed6173050ef6f2ee2ef95e183763a72a2c23fa2464"},{"id":"func/Game.Shoot","name":"Game.Shoot","line":21,"end_line":30,"hash":"e364fa427251ac0f01398e96eea81c26391b15317cd93b3d0812f4e87b324174"},{"id":"func/Game.followArrow","name":"Game.followArrow","line":32,"end_line":47,"hash":"8b1595212130e5df6970843ee6d4761953d560d78e52253266b51791bd9f81f0"},{"id":"func/Game.nextArrowRoom","name":"Game.nextArrowRoom","line":49,"end_line":54,"hash":"e33ee5c416e7d342a882a0fad59152d5f9bf4efac83249dbd4a7b7f4268de5ee"},{"id":"func/Game.nextArrowDeviationRoom","name":"Game.nextArrowDeviationRoom","line":56,"end_line":67,"hash":"0d73ab5cb7c2efe779026f1d664e1216b4f480e36ef8a240674835fe3ea1c27e"},{"id":"func/Game.missedArrowResult","name":"Game.missedArrowResult","line":69,"end_line":76,"hash":"aee1bc91e58ac4580b33ca0432e07dd47dfc545c20c67595136ac9712e921137"}]}
// mutate4go-manifest-end
