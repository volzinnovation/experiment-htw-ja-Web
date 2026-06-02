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
// {"version":1,"tested_at":"2026-06-02T10:04:29-05:00","module_hash":"1ad015b637be6fcb88660b2d826cd68bf31be3f28f381014c6959367250ef67d","functions":[{"id":"func/NewGame","name":"NewGame","line":34,"end_line":48,"hash":"07cc79edc0d8821b5fb1704ec51b84ecc2ed6a93b834e68e789d23af07f77514"},{"id":"func/NewGameWithSetup","name":"NewGameWithSetup","line":50,"end_line":55,"hash":"847c826d1fbfad99bb3df2fed4d7f17c087b78694b37a6ffacb37bbb45977abc"},{"id":"func/Game.Setup","name":"Game.Setup","line":57,"end_line":59,"hash":"f7cc96cad6b70f176db1ddb96c64a695c5bd288b753e0571fd8bbf63cdd9e47e"},{"id":"func/Game.ReplaySameSetup","name":"Game.ReplaySameSetup","line":61,"end_line":73,"hash":"b30aed25a41768030e29a9df572b40ca43d6b61e8548c58d0670d66f4b94ff30"},{"id":"func/Setup.OccupiedRooms","name":"Setup.OccupiedRooms","line":75,"end_line":80,"hash":"d59c1c96bcb74ca2f6ed1f8836a847f34402944c023947d7aa17d9e5647c4c95"},{"id":"func/validateSetup","name":"validateSetup","line":82,"end_line":87,"hash":"d64d4dac5265eafa42344468b4f35891956f66b96d11194b9564d81ebcafb6bd"},{"id":"func/validateHazardCounts","name":"validateHazardCounts","line":89,"end_line":97,"hash":"c2199a6a3399289560f5d60731668f19997faddca240a185e25981b0a5c0efd6"},{"id":"func/validateOccupiedRooms","name":"validateOccupiedRooms","line":99,"end_line":111,"hash":"038ee86ca133d4ad803d73e36a8f1bba5edd01fd83baa8896d4ade8e6ca9e44d"},{"id":"func/copySetup","name":"copySetup","line":113,"end_line":120,"hash":"9ec02e10c8f782b35c31584022513df5a955a9b6183b66943452fbaa4c08f879"},{"id":"func/copyRoomPtr","name":"copyRoomPtr","line":122,"end_line":128,"hash":"7f6bb99634aefd2a97efa2fdd2c25422cd3a08f5c56d74a5bc17719cefa358e5"}]}
// mutate4go-manifest-end
