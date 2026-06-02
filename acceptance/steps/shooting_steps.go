package steps

import (
	"fmt"
	"reflect"

	"htwgo/acceptance/runtime"
	"htwgo/internal/wumpus"
)

func givenShootingSetup(world *runtime.World, example map[string]string) error {
	player, err := intExample(example, "player_room")
	if err != nil {
		return err
	}
	wumpusRoom, err := intExample(example, "wumpus_room")
	if err != nil {
		return err
	}
	setup := wumpus.Setup{Player: player, Wumpus: wumpusRoom, Pits: []int{13, 14}, Bats: []int{16, 17}}
	return setConfiguredSetup(world, setup)
}

func givenPlayerStartsWithArrows(world *runtime.World, example map[string]string) error {
	arrows, err := intExample(example, "initial_arrows")
	if err != nil {
		return err
	}
	gameFrom(world, "game").SetArrows(arrows)
	return nil
}

func givenNextArrowDeviationRoom(world *runtime.World, example map[string]string) error {
	if example["deviation_room"] == "none" {
		return nil
	}
	room, err := intExample(example, "deviation_room")
	if err != nil {
		return err
	}
	gameFrom(world, "game").SetNextArrowDeviation(room)
	return nil
}

func whenPlayerShootsPath(world *runtime.World, example map[string]string) error {
	path, err := optionalRoomList(example["path"])
	if err != nil {
		return err
	}
	world.State["requested_shot_path"] = path
	result := gameFrom(world, "game").Shoot(path)
	world.State["shoot_result"] = result
	world.State["turn_messages"] = result.Messages
	world.State["action_taken"] = true
	return nil
}

func thenPlayerWins(world *runtime.World, _ map[string]string) error {
	return assertGameStatus(world, wumpus.StatusWon)
}

func thenArrowTraversedRoomsAre(world *runtime.World, example map[string]string) error {
	want, err := roomList(example["traversed_rooms"])
	if err != nil {
		return err
	}
	got := world.State["shoot_result"].(wumpus.ShootResult).TraversedRooms
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("arrow traversed rooms %v, want %v", got, want)
	}
	return nil
}

func thenRequestedShotPathIs(world *runtime.World, example map[string]string) error {
	want, err := optionalRoomList(example["expected_path"])
	if err != nil {
		return err
	}
	got := world.State["requested_shot_path"].([]int)
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("requested shot path %v, want %v", got, want)
	}
	return nil
}

func thenShotRejectedWithMessage(world *runtime.World, example map[string]string) error {
	result := world.State["shoot_result"].(wumpus.ShootResult)
	return assertRejectedMessage("shot rejection", result.RejectedMessage, example["message"])
}
