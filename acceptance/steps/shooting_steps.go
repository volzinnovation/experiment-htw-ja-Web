package steps

import (
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
	return assertRoomList("arrow traversed rooms", world.State["shoot_result"].(wumpus.ShootResult).TraversedRooms, example["traversed_rooms"])
}

func thenRequestedShotPathIs(world *runtime.World, example map[string]string) error {
	return assertOptionalRoomList("requested shot path", world.State["requested_shot_path"].([]int), example["expected_path"])
}

func thenShotRejectedWithMessage(world *runtime.World, example map[string]string) error {
	result := world.State["shoot_result"].(wumpus.ShootResult)
	return assertRejectedMessage("shot rejection", result.RejectedMessage, example["message"])
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T10:46:22-05:00","module_hash":"d2ea898b5aa044aa3ef18f2f658d8421a0b521c8909247edc96f1782dace02ba","functions":[{"id":"func/givenShootingSetup","name":"givenShootingSetup","line":8,"end_line":19,"hash":"8bdeed84eae626862372e24575ec192cf0feb36d5e95adf72fa5cba482ed1fd3"},{"id":"func/givenPlayerStartsWithArrows","name":"givenPlayerStartsWithArrows","line":21,"end_line":28,"hash":"1b923158fb159a9cfda965cb7f76f4a6df155cf51808c555249b90ecf7e35101"},{"id":"func/givenNextArrowDeviationRoom","name":"givenNextArrowDeviationRoom","line":30,"end_line":40,"hash":"31c1845db267929944f4d05712c88ffad48dd20c26f59fb4151e02b56d349c39"},{"id":"func/whenPlayerShootsPath","name":"whenPlayerShootsPath","line":42,"end_line":53,"hash":"8d6b4ba59fe052ba7fa5b4fc7d560f239b2d8146aba2350f452e783bfa2454bb"},{"id":"func/thenPlayerWins","name":"thenPlayerWins","line":55,"end_line":57,"hash":"36cca97def267964154d464c411a742df12a989815bee6fc5aa6e3e4bfc8ed27"},{"id":"func/thenArrowTraversedRoomsAre","name":"thenArrowTraversedRoomsAre","line":59,"end_line":61,"hash":"b33a3b2ab81f09c7c6128def6df92f3c087df1d87670c7585e6d123435034863"},{"id":"func/thenRequestedShotPathIs","name":"thenRequestedShotPathIs","line":63,"end_line":65,"hash":"6aba28cd8a037c0c65c7032fb2f49ac7a6359967b30b00501611e9835f39c7aa"},{"id":"func/thenShotRejectedWithMessage","name":"thenShotRejectedWithMessage","line":67,"end_line":70,"hash":"630ddbc78d74c3d4e038f253e337fddabb339c92fabc83313e68698605ed89ab"}]}
// mutate4go-manifest-end
