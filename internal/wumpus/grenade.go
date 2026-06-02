package wumpus

type GrenadeResult struct {
	RejectedMessage string
	Messages        []string
}

func (g *Game) SetGrenadeRoom(room int) {
	g.setGrenadeState(&g.grenadeRoom, room)
}

func (g Game) GrenadeRoom() (int, bool) {
	return roomPtrValue(g.grenadeRoom)
}

func (g *Game) GiveGrenade() {
	g.carriesGrenade = true
	g.grenadeRoom = nil
}

func (g Game) CarriesGrenade() bool {
	return g.carriesGrenade
}

func (g *Game) SetPendingGrenade(room int) {
	g.setGrenadeState(&g.pendingGrenadeRoom, room)
}

func (g Game) PendingGrenadeRoom() (int, bool) {
	return roomPtrValue(g.pendingGrenadeRoom)
}

func intPtr(value int) *int {
	return &value
}

func roomPtrValue(room *int) (int, bool) {
	if room == nil {
		return 0, false
	}
	return *room, true
}

func (g *Game) setGrenadeState(target **int, room int) {
	*target = intPtr(room)
	g.carriesGrenade = false
}

func (g *Game) ThrowGrenade(target int) GrenadeResult {
	if !g.carriesGrenade || !NewCave().HasTunnel(g.setup.Player, target) {
		return GrenadeResult{RejectedMessage: "CAN'T THROW THERE"}
	}
	g.carriesGrenade = false
	g.SetPendingGrenade(target)
	return GrenadeResult{Messages: []string{"YOU HEAR TIC...TIC..."}}
}

func (g *Game) DetonateGrenade() []string {
	target, ok := g.PendingGrenadeRoom()
	if !ok {
		return nil
	}
	g.pendingGrenadeRoom = nil
	blast := NewCave().blastRooms(target)
	messages := []string{"YOU HEAR A HORRENDOUS EXPLOSION!"}
	g.destroyBatsIn(blast)
	if blast[g.setup.Wumpus] {
		g.status = StatusWon
		return append(messages, "AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!")
	}
	if blast[g.setup.Player] {
		return append(messages, g.lose("YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!")...)
	}
	return messages
}

func (g Game) WumpusAlive() bool {
	return g.status != StatusWon
}

func (g *Game) collectGrenadeIfPresent() []string {
	if g.grenadeRoom == nil || *g.grenadeRoom != g.setup.Player {
		return nil
	}
	g.GiveGrenade()
	return []string{"YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!"}
}

func (g *Game) destroyBatsIn(blast map[int]bool) {
	var remaining []int
	for _, bat := range g.setup.Bats {
		if !blast[bat] {
			remaining = append(remaining, bat)
		}
	}
	g.setup.Bats = remaining
}

func (c Cave) blastRooms(target int) map[int]bool {
	blast := map[int]bool{target: true}
	exits, err := c.Exits(target)
	if err != nil {
		return blast
	}
	for _, exit := range exits {
		blast[exit] = true
	}
	if target == 10 {
		blast[5] = true
	}
	if target == 13 {
		delete(blast, 20)
	}
	return blast
}
