package steps

import (
	"fmt"
	"reflect"
	"slices"

	"htwgo/acceptance/runtime"
	"htwgo/internal/interactive"
	"htwgo/internal/wumpus"
)

type roomState struct {
	room int
	ok   bool
}

func givenConfiguredSetupWithGrenade(world *runtime.World, example map[string]string) error {
	if err := givenConfiguredSetup(world, example); err != nil {
		return err
	}
	room, err := intExample(example, "grenade_room")
	if err != nil {
		return err
	}
	gameFrom(world, "game").SetGrenadeRoom(room)
	return nil
}

func whenPlayerMovesToGrenadeRoom(world *runtime.World, example map[string]string) error {
	room, err := intExample(example, "grenade_room")
	if err != nil {
		return err
	}
	return movePlayerToRoom(world, room)
}

func givenInteractiveSetupCarryingGrenade(world *runtime.World, example map[string]string) error {
	arrows, err := intExample(example, "arrows")
	if err != nil {
		return err
	}
	game, err := wumpus.NewGameWithSetup(wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}})
	if err != nil {
		return err
	}
	game.SetArrows(arrows)
	game.GiveGrenade()
	storeSession(world, interactive.NewSessionWithGame(game))
	return nil
}

func thenOneHolyHandGrenade(world *runtime.World, _ map[string]string) error {
	if _, ok := gameFrom(world, "game").GrenadeRoom(); !ok {
		return fmt.Errorf("Holy Hand Grenade not placed")
	}
	return nil
}

func thenGrenadeRoomValid(world *runtime.World, _ map[string]string) error {
	room, ok := gameFrom(world, "game").GrenadeRoom()
	if !ok || room < 1 || room > 20 {
		return fmt.Errorf("grenade room = %d/%v, want valid", room, ok)
	}
	return nil
}

func thenGrenadeRoomUnoccupied(world *runtime.World, _ map[string]string) error {
	room, _ := gameFrom(world, "game").GrenadeRoom()
	if slices.Contains(gameFrom(world, "game").Setup().OccupiedRooms(), room) {
		return fmt.Errorf("grenade room %d overlaps occupied rooms", room)
	}
	return nil
}

func thenBothGrenadeRoomsIdentical(world *runtime.World, _ map[string]string) error {
	leftRoom, leftOK := gameFrom(world, "game").GrenadeRoom()
	rightRoom, rightOK := gameFrom(world, "another_game").GrenadeRoom()
	if leftOK != rightOK || leftRoom != rightRoom {
		return fmt.Errorf("grenade rooms differ: %d/%v and %d/%v", leftRoom, leftOK, rightRoom, rightOK)
	}
	return nil
}

func thenOrGivenPlayerCarriesGrenade(world *runtime.World, _ map[string]string) error {
	if _, actionTaken := world.State["action_taken"]; !actionTaken {
		gameFrom(world, "game").GiveGrenade()
		return nil
	}
	if messages, ok := world.State["turn_messages"].([]string); ok {
		if slices.Contains(messages, "YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!") {
			gameFrom(world, "game").GiveGrenade()
			return nil
		}
	}
	if !gameFrom(world, "game").CarriesGrenade() {
		return fmt.Errorf("player does not carry grenade")
	}
	return nil
}

func thenPlayerDoesNotCarryGrenade(world *runtime.World, _ map[string]string) error {
	if gameFrom(world, "game").CarriesGrenade() {
		return fmt.Errorf("player carries grenade")
	}
	return nil
}

func thenNoUnclaimedGrenade(world *runtime.World, _ map[string]string) error {
	room, ok := gameFrom(world, "game").GrenadeRoom()
	return assertNoRoom("unclaimed grenade remains", roomState{room: room, ok: ok})
}

func thenNoGrenadePending(world *runtime.World, _ map[string]string) error {
	room, ok := gameFrom(world, "game").PendingGrenadeRoom()
	return assertNoRoom("grenade pending", roomState{room: room, ok: ok})
}

func assertNoRoom(label string, state roomState) error {
	if state.ok {
		return fmt.Errorf("%s in room %d", label, state.room)
	}
	return nil
}

func givenOrThenGrenadePending(world *runtime.World, example map[string]string) error {
	room, err := intExample(example, "target_room")
	if err != nil {
		return err
	}
	if _, actionTaken := world.State["action_taken"]; !actionTaken {
		gameFrom(world, "game").SetPendingGrenade(room)
		return nil
	}
	got, ok := gameFrom(world, "game").PendingGrenadeRoom()
	if !ok || got != room {
		return fmt.Errorf("pending grenade room = %d/%v, want %d/true", got, ok, room)
	}
	return nil
}

func thenWumpusAlive(world *runtime.World, _ map[string]string) error {
	if !gameFrom(world, "game").WumpusAlive() {
		return fmt.Errorf("Wumpus is not alive")
	}
	return nil
}

func thenRemainingBatRooms(world *runtime.World, example map[string]string) error {
	return assertSetupRooms("remaining_bat_rooms", gameFrom(world, "game").Setup().Bats, "bat rooms", example)
}

func thenPitRooms(world *runtime.World, example map[string]string) error {
	return assertSetupRooms("pit_rooms", gameFrom(world, "game").Setup().Pits, "pit rooms", example)
}

func assertSetupRooms(key string, got []int, label string, example map[string]string) error {
	want, err := roomList(example[key])
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("%s = %v, want %v", label, got, want)
	}
	return nil
}

func thenReplaySetupIncludingGrenadeIdentical(world *runtime.World, _ map[string]string) error {
	session := sessionFrom(world)
	lost := world.State["lost_setup"].(wumpus.Setup)
	current := session.Game().Setup()
	room, ok := session.Game().GrenadeRoom()
	if ok != world.State["lost_grenade_ok"].(bool) || room != world.State["lost_grenade_room"].(int) {
		return fmt.Errorf("replay grenade room = %d/%v", room, ok)
	}
	if !reflect.DeepEqual(lost, current) {
		return fmt.Errorf("setup is %v, want identical to %v", current, lost)
	}
	return nil
}

func thenReplayPendingGrenadeRoom(world *runtime.World, example map[string]string) error {
	want, err := intExample(example, "target_room")
	if err != nil {
		return err
	}
	got, ok := sessionFrom(world).Game().PendingGrenadeRoom()
	if !ok || got != want {
		return fmt.Errorf("replay pending grenade room = %d/%v, want %d/true", got, ok, want)
	}
	return nil
}
