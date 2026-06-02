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
	nextBatRelocation  []int
	nextWumpusWake     []WumpusWakeChoice
	nextArrowDeviation []int
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
