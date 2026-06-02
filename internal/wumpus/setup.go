package wumpus

import (
	"fmt"
	"math/rand"
)

type Setup struct {
	Player int
	Wumpus int
	Pits   []int
	Bats   []int
}

type Game struct {
	setup              Setup
	status             Status
	arrows             int
	grenadeRoom        *int
	carriesGrenade     bool
	pendingGrenadeRoom *int
	wumpusAsleep       bool
	sawSleepingWumpus  bool
	nextBatRelocation  []int
	nextWumpusWake     []WumpusWakeChoice
	nextArrowDeviation []int
	nextSleepyObserve  []bool
	nextSleepingEntry  []SleepingWumpusEntryOutcome
	nextJumpEvents     []bool
	nextJumpPaths      [][]int
	nextFirstJumpLand  []FirstJumpLandingOutcome
}

func NewGame(seed int64) (Game, error) {
	rooms := rand.New(rand.NewSource(seed)).Perm(20)
	setup := Setup{
		Player: rooms[0] + 1,
		Wumpus: rooms[1] + 1,
		Pits:   []int{rooms[2] + 1, rooms[3] + 1},
		Bats:   []int{rooms[4] + 1, rooms[5] + 1},
	}
	game, err := NewGameWithSetup(setup)
	if err != nil {
		return Game{}, err
	}
	game.SetGrenadeRoom(rooms[6] + 1)
	return game, nil
}

func NewGameWithSetup(setup Setup) (Game, error) {
	if err := validateSetup(setup); err != nil {
		return Game{}, err
	}
	return Game{setup: copySetup(setup), status: StatusInProgress, arrows: 5}, nil
}

func (g Game) Setup() Setup {
	return copySetup(g.setup)
}

func (g Game) ReplaySameSetup() Game {
	replay := Game{
		setup:              copySetup(g.setup),
		status:             StatusInProgress,
		arrows:             g.arrows,
		carriesGrenade:     g.carriesGrenade,
		pendingGrenadeRoom: copyRoomPtr(g.pendingGrenadeRoom),
		grenadeRoom:        copyRoomPtr(g.grenadeRoom),
		wumpusAsleep:       g.wumpusAsleep,
		sawSleepingWumpus:  g.sawSleepingWumpus,
	}
	return replay
}

func (s Setup) OccupiedRooms() []int {
	rooms := []int{s.Player, s.Wumpus}
	rooms = append(rooms, s.Pits...)
	rooms = append(rooms, s.Bats...)
	return rooms
}

func validateSetup(setup Setup) error {
	if err := validateHazardCounts(setup); err != nil {
		return err
	}
	return validateOccupiedRooms(setup.OccupiedRooms())
}

func validateHazardCounts(setup Setup) error {
	if len(setup.Pits) != 2 {
		return fmt.Errorf("setup has %d pits, want 2", len(setup.Pits))
	}
	if len(setup.Bats) != 2 {
		return fmt.Errorf("setup has %d bats, want 2", len(setup.Bats))
	}
	return nil
}

func validateOccupiedRooms(rooms []int) error {
	seen := map[int]bool{}
	for _, room := range rooms {
		if room < 1 || room > 20 {
			return fmt.Errorf("invalid room %d", room)
		}
		if seen[room] {
			return fmt.Errorf("room %d is occupied more than once", room)
		}
		seen[room] = true
	}
	return nil
}

func copySetup(setup Setup) Setup {
	return Setup{
		Player: setup.Player,
		Wumpus: setup.Wumpus,
		Pits:   append([]int(nil), setup.Pits...),
		Bats:   append([]int(nil), setup.Bats...),
	}
}

func copyRoomPtr(room *int) *int {
	if room == nil {
		return nil
	}
	copy := *room
	return &copy
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T07:58:43-05:00","module_hash":"70be69d501ffd4046df9dce3799766bd2c964fce063bfc8be3c86eba2e0bac23","functions":[{"id":"func/NewGame","name":"NewGame","line":19,"end_line":28,"hash":"30ab8b9033ddb38767a89a9b98a93995d80e16571f92cff897cdadea771becc9"},{"id":"func/NewGameWithSetup","name":"NewGameWithSetup","line":30,"end_line":35,"hash":"921cb2171388bd7a04c89b9f74b7e1a56525729ee8afd5cd105923c89f09322e"},{"id":"func/Game.Setup","name":"Game.Setup","line":37,"end_line":39,"hash":"f7cc96cad6b70f176db1ddb96c64a695c5bd288b753e0571fd8bbf63cdd9e47e"},{"id":"func/Game.ReplaySameSetup","name":"Game.ReplaySameSetup","line":41,"end_line":43,"hash":"5b77f23f676e573a1b15abb356de9f86db5435639b65099d47454a4192540364"},{"id":"func/Setup.OccupiedRooms","name":"Setup.OccupiedRooms","line":45,"end_line":50,"hash":"d59c1c96bcb74ca2f6ed1f8836a847f34402944c023947d7aa17d9e5647c4c95"},{"id":"func/validateSetup","name":"validateSetup","line":52,"end_line":57,"hash":"d64d4dac5265eafa42344468b4f35891956f66b96d11194b9564d81ebcafb6bd"},{"id":"func/validateHazardCounts","name":"validateHazardCounts","line":59,"end_line":67,"hash":"c2199a6a3399289560f5d60731668f19997faddca240a185e25981b0a5c0efd6"},{"id":"func/validateOccupiedRooms","name":"validateOccupiedRooms","line":69,"end_line":81,"hash":"038ee86ca133d4ad803d73e36a8f1bba5edd01fd83baa8896d4ade8e6ca9e44d"},{"id":"func/copySetup","name":"copySetup","line":83,"end_line":90,"hash":"9ec02e10c8f782b35c31584022513df5a955a9b6183b66943452fbaa4c08f879"}]}
// mutate4go-manifest-end
