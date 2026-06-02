//go:build property

package wumpus

import (
	"testing"
	"testing/quick"
)

func TestNewGameProperties(t *testing.T) {
	property := func(seed int64) bool {
		game, err := NewGame(seed)
		if err != nil {
			return false
		}
		setup := game.Setup()
		return allRoomsValid(setup.OccupiedRooms()) && distinctCount(setup.OccupiedRooms()) == 6
	}
	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestCaveTunnelProperties(t *testing.T) {
	cave := NewCave()
	property := func(room uint8) bool {
		roomNumber := int(room%20) + 1
		exits, err := cave.Exits(roomNumber)
		if err != nil || len(exits) != 3 {
			return false
		}
		for _, exit := range exits {
			if exit == roomNumber || !cave.HasTunnel(exit, roomNumber) {
				return false
			}
		}
		return true
	}
	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestLegalEmptyMoveProperties(t *testing.T) {
	cave := NewCave()
	property := func(room uint8, exitIndex uint8) bool {
		from := int(room%20) + 1
		exits, err := cave.Exits(from)
		if err != nil {
			return false
		}
		to := exits[int(exitIndex)%len(exits)]
		setup := safeEmptyMoveSetup(from, to)
		game, err := NewGameWithSetup(setup)
		if err != nil {
			return false
		}

		result := game.Move(to)
		got := game.Setup()
		return result.RejectedMessage == "" &&
			len(result.Messages) == 0 &&
			game.Status() == StatusInProgress &&
			got.Player == to &&
			got.Wumpus == setup.Wumpus &&
			distinctCount(got.OccupiedRooms()) == 6
	}
	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestShootingPathProperties(t *testing.T) {
	property := func(pathLength uint8) bool {
		game, err := NewGameWithSetup(Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
		if err != nil {
			return false
		}
		game.SetArrows(5)
		game.SetNextWumpusWakeChoice(WumpusWakeChoice{Stay: true})

		path := safeArrowPath(int(pathLength % 7))
		result := game.Shoot(path)
		if len(path) < 1 || len(path) > 5 {
			return result.RejectedMessage == "CAN'T SHOOT THERE" &&
				game.Arrows() == 5 &&
				game.Status() == StatusInProgress
		}
		return result.RejectedMessage == "" &&
			len(result.TraversedRooms) == len(path) &&
			game.Arrows() == 4
	}
	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func safeArrowPath(count int) []int {
	rooms := []int{2, 3, 4, 5, 6}
	if count > len(rooms) {
		return append(rooms, 7)
	}
	return rooms[:count]
}

func safeEmptyMoveSetup(from, to int) Setup {
	excluded := map[int]bool{from: true, to: true}
	rooms := roomsOutside(excluded, 5)
	return Setup{
		Player: from,
		Wumpus: rooms[0],
		Pits:   []int{rooms[1], rooms[2]},
		Bats:   []int{rooms[3], rooms[4]},
	}
}

func roomsOutside(excluded map[int]bool, count int) []int {
	rooms := make([]int, 0, count)
	for room := 1; room <= 20 && len(rooms) < count; room++ {
		if !excluded[room] {
			rooms = append(rooms, room)
		}
	}
	return rooms
}

func allRoomsValid(rooms []int) bool {
	for _, room := range rooms {
		if room < 1 || room > 20 {
			return false
		}
	}
	return true
}
