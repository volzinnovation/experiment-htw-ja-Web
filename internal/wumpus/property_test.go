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

func allRoomsValid(rooms []int) bool {
	for _, room := range rooms {
		if room < 1 || room > 20 {
			return false
		}
	}
	return true
}
