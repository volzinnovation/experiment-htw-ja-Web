package wumpus

type GrenadeResult struct {
	RejectedMessage string
	Messages        []string
}

func (g *Game) SetGrenadeRoom(room int) {
	g.grenadeRoom = &room
	g.carriesGrenade = false
}

func (g Game) GrenadeRoom() (int, bool) {
	if g.grenadeRoom == nil {
		return 0, false
	}
	return *g.grenadeRoom, true
}

func (g *Game) GiveGrenade() {
	g.carriesGrenade = true
	g.grenadeRoom = nil
}

func (g Game) CarriesGrenade() bool {
	return g.carriesGrenade
}

func (g *Game) SetPendingGrenade(room int) {
	g.pendingGrenadeRoom = &room
	g.carriesGrenade = false
}

func (g Game) PendingGrenadeRoom() (int, bool) {
	if g.pendingGrenadeRoom == nil {
		return 0, false
	}
	return *g.pendingGrenadeRoom, true
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
	return blast
}
