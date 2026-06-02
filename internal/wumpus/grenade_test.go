package wumpus

import (
	"reflect"
	"slices"
	"testing"
)

func TestNewGamePlacesGrenadeInUnoccupiedRoom(t *testing.T) {
	game, err := NewGame(1973)
	if err != nil {
		t.Fatal(err)
	}

	room, ok := game.GrenadeRoom()
	if !ok {
		t.Fatal("grenade was not placed")
	}
	if room < 1 || room > 20 {
		t.Fatalf("grenade room = %d, want valid room", room)
	}
	if slices.Contains(game.Setup().OccupiedRooms(), room) {
		t.Fatalf("grenade room %d overlaps occupied rooms %v", room, game.Setup().OccupiedRooms())
	}
}

func TestMovingIntoGrenadeRoomAcquiresIt(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetGrenadeRoom(2)

	result := game.Move(2)

	if !game.CarriesGrenade() {
		t.Fatal("player should carry grenade")
	}
	if _, ok := game.GrenadeRoom(); ok {
		t.Fatal("unclaimed grenade should be gone")
	}
	want := []string{"YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!"}
	if !reflect.DeepEqual(result.Messages, want) {
		t.Fatalf("messages = %v, want %v", result.Messages, want)
	}
}

func TestThrowGrenadeStartsPendingDetonation(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.GiveGrenade()

	result := game.ThrowGrenade(2)

	if result.RejectedMessage != "" {
		t.Fatalf("throw rejected: %s", result.RejectedMessage)
	}
	if game.CarriesGrenade() {
		t.Fatal("grenade should be consumed from inventory")
	}
	room, ok := game.PendingGrenadeRoom()
	if !ok || room != 2 {
		t.Fatalf("pending grenade room = %d %v, want 2 true", room, ok)
	}
	if !reflect.DeepEqual(result.Messages, []string{"YOU HEAR TIC...TIC..."}) {
		t.Fatalf("messages = %v", result.Messages)
	}
}

func TestThrowGrenadeRequiresCarriedGrenadeAndAdjacentTarget(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})

	if result := game.ThrowGrenade(2); result.RejectedMessage != "CAN'T THROW THERE" {
		t.Fatalf("throw without grenade = %#v", result)
	}
	game.GiveGrenade()
	if result := game.ThrowGrenade(20); result.RejectedMessage != "CAN'T THROW THERE" {
		t.Fatalf("throw to non-adjacent room = %#v", result)
	}
	if !game.CarriesGrenade() {
		t.Fatal("invalid throw should not consume grenade")
	}
}

func TestNoPendingGrenadeReportsNoRoomAndDetonatesQuietly(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})

	if room, ok := game.PendingGrenadeRoom(); ok || room != 0 {
		t.Fatalf("pending grenade room = %d/%v, want 0/false", room, ok)
	}
	if messages := game.DetonateGrenade(); len(messages) != 0 {
		t.Fatalf("detonation messages = %v, want none", messages)
	}
	if !game.WumpusAlive() {
		t.Fatal("Wumpus should be alive")
	}
}

func TestGrenadeDetonationCanBeHarmless(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetPendingGrenade(10)

	messages := game.DetonateGrenade()

	if !reflect.DeepEqual(messages, []string{"YOU HEAR A HORRENDOUS EXPLOSION!"}) {
		t.Fatalf("messages = %v", messages)
	}
	if game.Status() != StatusInProgress {
		t.Fatalf("status = %s, want %s", game.Status(), StatusInProgress)
	}
}

func TestGrenadeDetonationKillsWumpusDestroysBatsAndLeavesPits(t *testing.T) {
	game := mustGame(t, Setup{Player: 1, Wumpus: 10, Pits: []int{9, 14}, Bats: []int{2, 16}})
	game.SetPendingGrenade(10)

	messages := game.DetonateGrenade()

	wantMessages := []string{
		"YOU HEAR A HORRENDOUS EXPLOSION!",
		"AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!",
	}
	if !reflect.DeepEqual(messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", messages, wantMessages)
	}
	if game.Status() != StatusWon {
		t.Fatalf("status = %s, want %s", game.Status(), StatusWon)
	}
	if !reflect.DeepEqual(game.Setup().Bats, []int{16}) {
		t.Fatalf("bats = %v, want [16]", game.Setup().Bats)
	}
	if !reflect.DeepEqual(game.Setup().Pits, []int{9, 14}) {
		t.Fatalf("pits = %v, want unchanged", game.Setup().Pits)
	}
}

func TestGrenadeDetonationCanKillPlayer(t *testing.T) {
	game := mustGame(t, Setup{Player: 2, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	game.SetPendingGrenade(10)

	messages := game.DetonateGrenade()

	wantMessages := []string{
		"YOU HEAR A HORRENDOUS EXPLOSION!",
		"YOU ARE BLOWN UP BY YOUR OWN HOLY HAND GRENADE!",
		"HA HA HA - YOU LOSE!",
	}
	if !reflect.DeepEqual(messages, wantMessages) {
		t.Fatalf("messages = %v, want %v", messages, wantMessages)
	}
	if game.Status() != StatusLost {
		t.Fatalf("status = %s, want %s", game.Status(), StatusLost)
	}
}
