package wumpus

type GrenadeResult struct {
	RejectedMessage string
	Messages        []string
}

func (g *Game) SetGrenadeRoom(room int) {
	g.setGrenadeState(&g.grenadeRoom, room)
}

func (g *Game) ClearGrenadeRoom() {
	g.grenadeRoom = nil
	g.carriesGrenade = false
	g.pendingGrenadeRoom = nil
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
	if g.inertHazards && (blast[g.setup.Wumpus] || blast[g.setup.Player]) {
		return append(messages, "QA INERT: BLAST IGNORED")
	}
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
	return blast
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:29:53-05:00","module_hash":"d85b85eebce9f17ef4993cc189765d9f2336d3f62dee516aab4c1ea5b75ad84b","functions":[{"id":"func/Game.SetGrenadeRoom","name":"Game.SetGrenadeRoom","line":8,"end_line":10,"hash":"ef5245398b2efd2618cd37da8bbc2ec23b925b9e8eccc635f6ca8edd650c236d"},{"id":"func/Game.GrenadeRoom","name":"Game.GrenadeRoom","line":12,"end_line":14,"hash":"28be4315c1025c9300aeb8758757929f2777a883cc20ce739f8b9efa41bf37e0"},{"id":"func/Game.GiveGrenade","name":"Game.GiveGrenade","line":16,"end_line":19,"hash":"c3554fbffdd96c815a6237dcfb519ba7e24cc9cfd375289639c0a5627d7138a2"},{"id":"func/Game.CarriesGrenade","name":"Game.CarriesGrenade","line":21,"end_line":23,"hash":"e68d89572c8481dba654ff88357d578776d15dd9d27f37db68a37911009cb468"},{"id":"func/Game.SetPendingGrenade","name":"Game.SetPendingGrenade","line":25,"end_line":27,"hash":"319c16036884eef6521f9d24285dc52058a20ff19586cdbae2360825c8d6aa91"},{"id":"func/Game.PendingGrenadeRoom","name":"Game.PendingGrenadeRoom","line":29,"end_line":31,"hash":"d2599784bb63a154271ba279cea65cef856091e5f1b2fc262013ade75b447112"},{"id":"func/intPtr","name":"intPtr","line":33,"end_line":35,"hash":"7a645be30de7ed50c987acd6e91b15724331c319961ebde21aebae1748a417ac"},{"id":"func/roomPtrValue","name":"roomPtrValue","line":37,"end_line":42,"hash":"0ce97497263759428a604c5feaca58cb8aba9ec9c9eea0486c4debbcc56ddd7e"},{"id":"func/Game.setGrenadeState","name":"Game.setGrenadeState","line":44,"end_line":47,"hash":"a526837593ec43be380f45b3f20031bf1dd3967ca91607e892dca84540b6d04b"},{"id":"func/Game.ThrowGrenade","name":"Game.ThrowGrenade","line":49,"end_line":56,"hash":"988a183e34eca8397b1ad720d39d48a22a60e126fe0a87ef5f93107566d3081b"},{"id":"func/Game.DetonateGrenade","name":"Game.DetonateGrenade","line":58,"end_line":75,"hash":"c8d64904471f57884da890aa3893a3c96f51c30f61da88ca987d21edd957c325"},{"id":"func/Game.WumpusAlive","name":"Game.WumpusAlive","line":77,"end_line":79,"hash":"ef7b82f1297ebbf45a557f355d78c8c4678c1b899cb633b95f3adbe0601c92e8"},{"id":"func/Game.collectGrenadeIfPresent","name":"Game.collectGrenadeIfPresent","line":81,"end_line":87,"hash":"93e0c738e0d013618da670dab65bba68a81bf2c5cb74d32dfc96e2550f974e21"},{"id":"func/Game.destroyBatsIn","name":"Game.destroyBatsIn","line":89,"end_line":97,"hash":"5c332c6c0781ae8b9f904f4c980843a8d20d7a43d7a83cf1354d94f7b5bd44b9"},{"id":"func/Cave.blastRooms","name":"Cave.blastRooms","line":99,"end_line":115,"hash":"d1b82a970456c983609beda40a37cce65dfbf240eb3134319dd69e3a2929d48a"}]}
// mutate4go-manifest-end
