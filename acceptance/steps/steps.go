package steps

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"htwgo/acceptance/runtime"
	"htwgo/internal/interactive"
	"htwgo/internal/wumpus"
)

type inspectedSetup struct {
	Player int
	Wumpus int
	Pits   []int
	Bats   []int
}

func NewHandlers() runtime.Handlers {
	return runtime.Handlers{
		"a new cave":                                                    givenNewCave,
		"the exits for room <room> are queried":                         whenExitsQueried,
		"the exits are <exits>":                                         thenExitsAre,
		"the exit count is <exit_count>":                                thenExitCountIs,
		"the cave is traversed from room 1":                             whenCaveTraversed,
		"every room from 1 through 20 is reachable":                     thenEveryRoomReachable,
		"the cave invariants are inspected":                             whenCaveInvariantsInspected,
		"every room has exactly three exits":                            thenEveryRoomHasThreeExits,
		"no room is one of its own exits":                               thenNoRoomIsItsOwnExit,
		"the tunnel from room <from_room> to room <to_room> is queried": whenTunnelQueried,
		"the reverse tunnel also exists":                                thenReverseTunnelExists,
		"room <room> is not one of the exits":                           thenRoomIsNotExit,
		"a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>": givenConfiguredSetup,
		"a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>":   givenConfiguredSetup,
		"a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 3 and 4, and bats in rooms 5 and 6":                                 givenSetupRoom1WumpusBats,
		"a game setup with the player in room 10, the Wumpus in room 2, pits in rooms 9 and 18, and bats in rooms 6 and 7":                               givenSetupRoom10WumpusPit,
		"a game setup with the player in room 6, the Wumpus in room 20, pits in rooms 1 and 2, and bats in rooms 3 and 4":                                givenSetupRoom6NoHazards,
		"adjacent hazards are queried from the player room":                                                                                              whenAdjacentHazardsQueried,
		"the adjacent hazard types are <hazards>":         thenAdjacentHazardsAre,
		"the adjacent hazard types are Wumpus and Bats":   thenAdjacentHazardsWumpusBats,
		"the adjacent hazard types are Wumpus and Pit":    thenAdjacentHazardsWumpusPit,
		"there are no adjacent hazard types":              thenNoAdjacentHazards,
		"a new game created with seed 1973":               givenNewGameSeed1973,
		"a new game created with seed <seed>":             givenNewGameSeed,
		"the setup is inspected":                          whenSetupInspected,
		"there is 1 player":                               thenOnePlayer,
		"there is 1 Wumpus":                               thenOneWumpus,
		"there are 2 pits":                                thenTwoPits,
		"there are 2 bats":                                thenTwoBats,
		"the occupied rooms are inspected":                whenOccupiedRoomsInspected,
		"every occupied room number is from 1 through 20": thenOccupiedRoomsValid,
		"exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats": thenDistinctOccupiedCount,
		"exactly 6 distinct rooms are occupied by the player, Wumpus, pits, and bats":                thenSixDistinctOccupiedRooms,
		"another new game created with seed 1973":                                                    givenAnotherNewGameSeed1973,
		"another new game created with seed <seed>":                                                  givenAnotherNewGameSeed,
		"both setups are inspected":                                                                  whenBothSetupsInspected,
		"both setups have identical player, Wumpus, pit, and bat rooms":                              thenBothSetupsIdentical,
		"a completed game created with seed 1973":                                                    givenCompletedGameSeed1973,
		"a completed game created with seed <seed>":                                                  givenCompletedGameSeed,
		"a same setup replay is started":                                                             whenSameSetupReplayStarted,
		"the replay setup has identical player, Wumpus, pit, and bat rooms":                          thenReplaySetupIdentical,
		"the player moves to room <to_room>":                                                         whenPlayerMoves,
		"the player moves to room <pit_room>":                                                        whenPlayerMoves,
		"the player moves to room <bat_room>":                                                        whenPlayerMoves,
		"the player moves to room <wumpus_room>":                                                     whenPlayerMoves,
		"the player moves to room <grenade_room>":                                                    whenPlayerMovesToGrenadeRoom,
		"the player is in room <to_room>":                                                            thenPlayerInToRoom,
		"the player is in room <from_room>":                                                          thenPlayerInFromRoom,
		"the player is in room <player_room>":                                                        thenPlayerInPlayerRoom,
		"the player is in room <relocation_room>":                                                    thenPlayerInRelocationRoom,
		"the game is still in progress":                                                              thenGameStillInProgress,
		"the game is in progress":                                                                    thenGameStillInProgress,
		"the game is <game_status>":                                                                  thenGameStatus,
		"the player loses":                                                                           thenPlayerLoses,
		"the turn messages are <messages>":                                                           thenTurnMessagesAre,
		"the move is rejected with message <message>":                                                thenMoveRejectedWithMessage,
		"the next bat relocation room is <relocation_room>":                                          givenNextBatRelocationRoom,
		"the next bat relocation room is <wumpus_room>":                                              givenNextBatRelocationRoom,
		"the next Wumpus wake choice is <wake_choice>":                                               givenNextWumpusWakeChoice,
		"the Wumpus is in room <expected_wumpus_room>":                                               thenWumpusInRoom,
		"the turn warnings are requested":                                                            whenTurnWarningsRequested,
		"the warning messages are <warnings>":                                                        thenWarningMessagesAre,
		"the player has <arrows> arrows":                                                             givenOrThenPlayerHasArrows,
		"the player has <remaining_arrows> arrows":                                                   thenPlayerHasRemainingArrows,
		"the next arrow deviation room is <deviation_room>":                                          givenNextArrowDeviationRoom,
		"the player shoots the path <path>":                                                          whenPlayerShootsPath,
		"the player wins":                                                                            thenPlayerWins,
		"the arrow traversed rooms are <traversed_rooms>":                                            thenArrowTraversedRoomsAre,
		"the shot is rejected with message <message>":                                                thenShotRejectedWithMessage,
		"an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows": givenInteractiveSetup,
		"an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows":   givenInteractiveSetup,
		"an interactive game setup with seed <seed>":                                      givenInteractiveSetupSeed,
		"a new interactive session":                                                       givenNewInteractiveSession,
		"the next turn is displayed":                                                      whenNextTurnDisplayed,
		"the player enters command <command>":                                             whenPlayerEntersCommand,
		"the displayed lines are <lines>":                                                 thenDisplayedLinesAre,
		"the displayed lines include <message>":                                           thenDisplayedLinesInclude,
		"the displayed lines include <messages>":                                          thenDisplayedLinesInclude,
		"the displayed lines include <prompt>":                                            thenDisplayedLinesInclude,
		"the player has lost":                                                             givenPlayerHasLost,
		"the player answers same setup prompt with <answer>":                              whenPlayerAnswersSameSetupPrompt,
		"the next game setup is <setup_relation> to the lost game setup":                  thenNextGameSetupRelation,
		"the player answers instructions prompt with <answer>":                            whenPlayerAnswersInstructionsPrompt,
		"there is 1 Holy Hand Grenade":                                                    thenOneHolyHandGrenade,
		"the Holy Hand Grenade room is from 1 through 20":                                 thenGrenadeRoomValid,
		"the Holy Hand Grenade room is not occupied by the player, Wumpus, pits, or bats": thenGrenadeRoomUnoccupied,
		"both setups have identical Holy Hand Grenade rooms":                              thenBothGrenadeRoomsIdentical,
		"a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and the Holy Hand Grenade in room <grenade_room>":   givenConfiguredSetupWithGrenade,
		"a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and the Holy Hand Grenade in room <grenade_room>": givenConfiguredSetupWithGrenade,
		"the player carries the Holy Hand Grenade":                                                         thenOrGivenPlayerCarriesGrenade,
		"the player does not carry the Holy Hand Grenade":                                                  thenPlayerDoesNotCarryGrenade,
		"the cave has no unclaimed Holy Hand Grenade":                                                      thenNoUnclaimedGrenade,
		"an interactive game setup where the player carries the Holy Hand Grenade and has <arrows> arrows": givenInteractiveSetupCarryingGrenade,
		"no Holy Hand Grenade detonation is pending":                                                       thenNoGrenadePending,
		"the Holy Hand Grenade is pending detonation in room <target_room>":                                givenOrThenGrenadePending,
		"the Wumpus is alive":                               thenWumpusAlive,
		"the remaining bat rooms are <remaining_bat_rooms>": thenRemainingBatRooms,
		"the pit rooms are <pit_rooms>":                     thenPitRooms,
		"the replay setup has identical player, Wumpus, pit, bat, and Holy Hand Grenade rooms": thenReplaySetupIncludingGrenadeIdentical,
		"the replay Holy Hand Grenade pending detonation room is <target_room>":                thenReplayPendingGrenadeRoom,
		"the next sleepy Wumpus adjacent observation is <sleepy_observation>":                  givenNextSleepyWumpusObservation,
		"the Wumpus is asleep":                                           givenWumpusAsleep,
		"the Wumpus is awake":                                            givenWumpusAwake,
		"the Wumpus sleep state is <sleep_state>":                        thenWumpusSleepState,
		"the Wumpus sleep state is asleep":                               thenWumpusSleepStateAsleep,
		"the Wumpus sleep state is awake":                                thenWumpusSleepStateAwake,
		"the turn warnings are <warnings>":                               thenTurnWarningsAre,
		"the next sleeping Wumpus room entry outcome is <entry_outcome>": givenNextSleepingWumpusEntryOutcome,
		"a game setup with the player in room <wumpus_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>": givenConfiguredSetupPlayerWithWumpus,
		"the player has seen the sleeping Wumpus shape":                    givenPlayerHasSeenSleepingWumpus,
		"both games observe sleepy Wumpus behavior for <turn_count> turns": whenBothGamesObserveSleepyWumpus,
		"both games produce identical sleepy Wumpus observations":          thenBothSleepyObservationsIdentical,
	}
}

func givenNewCave(world *runtime.World, _ map[string]string) error {
	world.State["cave"] = wumpus.NewCave()
	return nil
}

func whenExitsQueried(world *runtime.World, example map[string]string) error {
	cave := world.State["cave"].(wumpus.Cave)
	room, err := intExample(example, "room")
	if err != nil {
		return err
	}
	exits, err := cave.Exits(room)
	if err != nil {
		return err
	}
	world.State["exits"] = exits
	return nil
}

func thenExitsAre(world *runtime.World, example map[string]string) error {
	want, err := roomList(example["exits"])
	if err != nil {
		return err
	}
	got := world.State["exits"].([]int)
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("got exits %v, want %v", got, want)
	}
	return nil
}

func thenExitCountIs(world *runtime.World, example map[string]string) error {
	want, err := intExample(example, "exit_count")
	if err != nil {
		return err
	}
	got := len(world.State["exits"].([]int))
	if got != want {
		return fmt.Errorf("got exit count %d, want %d", got, want)
	}
	return nil
}

func whenCaveTraversed(world *runtime.World, _ map[string]string) error {
	world.State["reachable"] = world.State["cave"].(wumpus.Cave).ReachableFrom(1)
	return nil
}

func thenEveryRoomReachable(world *runtime.World, _ map[string]string) error {
	reachable := world.State["reachable"].([]int)
	if len(reachable) != 20 {
		return fmt.Errorf("got reachable rooms %v", reachable)
	}
	for room := 1; room <= 20; room++ {
		if !contains(reachable, room) {
			return fmt.Errorf("room %d is not reachable", room)
		}
	}
	return nil
}

func whenCaveInvariantsInspected(world *runtime.World, _ map[string]string) error {
	cave := world.State["cave"].(wumpus.Cave)
	exitsByRoom := map[int][]int{}
	for room := 1; room <= 20; room++ {
		exits, err := cave.Exits(room)
		if err != nil {
			return err
		}
		exitsByRoom[room] = exits
	}
	world.State["exits_by_room"] = exitsByRoom
	return nil
}

func thenEveryRoomHasThreeExits(world *runtime.World, _ map[string]string) error {
	for room, exits := range world.State["exits_by_room"].(map[int][]int) {
		if len(exits) != 3 {
			return fmt.Errorf("room %d has exits %v", room, exits)
		}
	}
	return nil
}

func thenNoRoomIsItsOwnExit(world *runtime.World, _ map[string]string) error {
	for room, exits := range world.State["exits_by_room"].(map[int][]int) {
		if contains(exits, room) {
			return fmt.Errorf("room %d connects to itself", room)
		}
	}
	return nil
}

func whenTunnelQueried(world *runtime.World, example map[string]string) error {
	cave := world.State["cave"].(wumpus.Cave)
	from, err := intExample(example, "from_room")
	if err != nil {
		return err
	}
	to, err := intExample(example, "to_room")
	if err != nil {
		return err
	}
	world.State["from_room"] = from
	world.State["to_room"] = to
	world.State["tunnel_exists"] = cave.HasTunnel(from, to)
	return nil
}

func thenReverseTunnelExists(world *runtime.World, _ map[string]string) error {
	cave := world.State["cave"].(wumpus.Cave)
	from := world.State["from_room"].(int)
	to := world.State["to_room"].(int)
	if !world.State["tunnel_exists"].(bool) {
		return fmt.Errorf("tunnel %d to %d does not exist", from, to)
	}
	if !cave.HasTunnel(to, from) {
		return fmt.Errorf("reverse tunnel %d to %d does not exist", to, from)
	}
	return nil
}

func thenRoomIsNotExit(world *runtime.World, example map[string]string) error {
	room, err := intExample(example, "room")
	if err != nil {
		return err
	}
	for _, exit := range world.State["exits"].([]int) {
		if exit == room {
			return fmt.Errorf("room %d connects to itself", room)
		}
	}
	return nil
}

func givenConfiguredSetup(world *runtime.World, example map[string]string) error {
	player, err := intAnyExample(example, "player_room", "from_room")
	if err != nil {
		return err
	}
	wumpusRoom, err := intExample(example, "wumpus_room")
	if err != nil {
		return err
	}
	pits, err := roomList(example["pit_rooms"])
	if err != nil {
		return err
	}
	bats, err := roomList(example["bat_rooms"])
	if err != nil {
		return err
	}
	return setConfiguredSetup(world, wumpus.Setup{Player: player, Wumpus: wumpusRoom, Pits: pits, Bats: bats})
}

func givenSetupRoom1WumpusBats(world *runtime.World, _ map[string]string) error {
	return setConfiguredSetup(world, wumpus.Setup{Player: 1, Wumpus: 2, Pits: []int{3, 4}, Bats: []int{5, 6}})
}

func givenSetupRoom10WumpusPit(world *runtime.World, _ map[string]string) error {
	return setConfiguredSetup(world, wumpus.Setup{Player: 10, Wumpus: 2, Pits: []int{9, 18}, Bats: []int{6, 7}})
}

func givenSetupRoom6NoHazards(world *runtime.World, _ map[string]string) error {
	return setConfiguredSetup(world, wumpus.Setup{Player: 6, Wumpus: 20, Pits: []int{1, 2}, Bats: []int{3, 4}})
}

func setConfiguredSetup(world *runtime.World, setup wumpus.Setup) error {
	player := setup.Player
	if setup.Player == setup.Wumpus {
		setup.Player = firstUnoccupiedRoom(setup.Wumpus, setup.Pits, setup.Bats)
	}
	game, err := wumpus.NewGameWithSetup(setup)
	if err != nil {
		return err
	}
	game.SetPlayerRoom(player)
	world.State["setup"] = setup
	world.State["game"] = &game
	return nil
}

func givenConfiguredSetupPlayerWithWumpus(world *runtime.World, example map[string]string) error {
	return givenConfiguredSetup(world, map[string]string{
		"player_room": example["wumpus_room"],
		"wumpus_room": example["wumpus_room"],
		"pit_rooms":   example["pit_rooms"],
		"bat_rooms":   example["bat_rooms"],
	})
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

func whenAdjacentHazardsQueried(world *runtime.World, _ map[string]string) error {
	setup := world.State["setup"].(wumpus.Setup)
	world.State["hazards"] = wumpus.NewCave().AdjacentHazards(setup.Player, setup)
	return nil
}

func thenAdjacentHazardsAre(world *runtime.World, example map[string]string) error {
	return requireHazards(world, hazardList(example["hazards"]))
}

func thenAdjacentHazardsWumpusBats(world *runtime.World, _ map[string]string) error {
	return requireHazards(world, []string{"Wumpus", "Bats"})
}

func thenAdjacentHazardsWumpusPit(world *runtime.World, _ map[string]string) error {
	return requireHazards(world, []string{"Wumpus", "Pit"})
}

func thenNoAdjacentHazards(world *runtime.World, _ map[string]string) error {
	return requireHazards(world, nil)
}

func requireHazards(world *runtime.World, want []string) error {
	got := world.State["hazards"].([]wumpus.Hazard)
	var gotStrings []string
	for _, hazard := range got {
		gotStrings = append(gotStrings, string(hazard))
	}
	if !reflect.DeepEqual(gotStrings, want) {
		return fmt.Errorf("got hazards %v, want %v", gotStrings, want)
	}
	return nil
}

func givenNewGameSeed1973(world *runtime.World, _ map[string]string) error {
	return setGame(world, 1973, "game")
}

func givenNewGameSeed(world *runtime.World, example map[string]string) error {
	return setGameFromExample(world, example, "game")
}

func givenAnotherNewGameSeed(world *runtime.World, example map[string]string) error {
	return setGameFromExample(world, example, "another_game")
}

func givenAnotherNewGameSeed1973(world *runtime.World, _ map[string]string) error {
	return setGame(world, 1973, "another_game")
}

func givenCompletedGameSeed(world *runtime.World, example map[string]string) error {
	return setGameFromExample(world, example, "game")
}

func givenCompletedGameSeed1973(world *runtime.World, _ map[string]string) error {
	return setGame(world, 1973, "game")
}

func setGameFromExample(world *runtime.World, example map[string]string, key string) error {
	seed, err := int64Example(example, "seed")
	if err != nil {
		return err
	}
	return setGame(world, seed, key)
}

func setGame(world *runtime.World, seed int64, key string) error {
	game, err := wumpus.NewGame(seed)
	if err != nil {
		return err
	}
	world.State[key] = &game
	return nil
}

func whenSetupInspected(world *runtime.World, _ map[string]string) error {
	world.State["inspected_setup"] = setupSnapshot(gameFrom(world, "game").Setup())
	return nil
}

func thenOnePlayer(world *runtime.World, _ map[string]string) error {
	return requirePlaced(inspected(world).Player, "player not placed")
}

func thenOneWumpus(world *runtime.World, _ map[string]string) error {
	return requirePlaced(inspected(world).Wumpus, "Wumpus not placed")
}

func thenTwoPits(world *runtime.World, _ map[string]string) error {
	return requireCount(len(inspected(world).Pits), 2, "pits")
}

func thenTwoBats(world *runtime.World, _ map[string]string) error {
	return requireCount(len(inspected(world).Bats), 2, "bats")
}

func inspected(world *runtime.World) inspectedSetup {
	return world.State["inspected_setup"].(inspectedSetup)
}

func requirePlaced(room int, message string) error {
	if room == 0 {
		return fmt.Errorf("%s", message)
	}
	return nil
}

func requireCount(got, want int, name string) error {
	if got != want {
		return fmt.Errorf("got %d %s", got, name)
	}
	return nil
}

func whenOccupiedRoomsInspected(world *runtime.World, _ map[string]string) error {
	setup := gameFrom(world, "game").Setup()
	world.State["occupied_rooms"] = setup.OccupiedRooms()
	return nil
}

func thenOccupiedRoomsValid(world *runtime.World, _ map[string]string) error {
	for _, room := range world.State["occupied_rooms"].([]int) {
		if room < 1 || room > 20 {
			return fmt.Errorf("invalid room %d", room)
		}
	}
	return nil
}

func thenDistinctOccupiedCount(world *runtime.World, example map[string]string) error {
	want, err := intExample(example, "occupied_count")
	if err != nil {
		return err
	}
	return requireDistinctOccupiedCount(world, want)
}

func thenSixDistinctOccupiedRooms(world *runtime.World, _ map[string]string) error {
	return requireDistinctOccupiedCount(world, 6)
}

func requireDistinctOccupiedCount(world *runtime.World, want int) error {
	rooms := world.State["occupied_rooms"].([]int)
	distinct := map[int]bool{}
	for _, room := range rooms {
		distinct[room] = true
	}
	if len(distinct) != want {
		return fmt.Errorf("got %d distinct occupied rooms from %v, want %d", len(distinct), rooms, want)
	}
	return nil
}

func whenBothSetupsInspected(world *runtime.World, _ map[string]string) error {
	world.State["inspected_setup"] = setupSnapshot(gameFrom(world, "game").Setup())
	world.State["another_inspected_setup"] = setupSnapshot(gameFrom(world, "another_game").Setup())
	return nil
}

func thenBothSetupsIdentical(world *runtime.World, _ map[string]string) error {
	left := world.State["inspected_setup"].(inspectedSetup)
	right := world.State["another_inspected_setup"].(inspectedSetup)
	if !reflect.DeepEqual(left, right) {
		return fmt.Errorf("setups differ: %v and %v", left, right)
	}
	return nil
}

func whenSameSetupReplayStarted(world *runtime.World, _ map[string]string) error {
	replay := gameFrom(world, "game").ReplaySameSetup()
	world.State["replay_setup"] = setupSnapshot(replay.Setup())
	return nil
}

func thenReplaySetupIdentical(world *runtime.World, _ map[string]string) error {
	original := setupSnapshot(gameFrom(world, "game").Setup())
	replay := world.State["replay_setup"].(inspectedSetup)
	if !reflect.DeepEqual(original, replay) {
		return fmt.Errorf("replay setup %v differs from original %v", replay, original)
	}
	return nil
}

func whenPlayerMoves(world *runtime.World, example map[string]string) error {
	room, err := intAnyExample(example, "to_room", "pit_room", "bat_room", "wumpus_room")
	if err != nil {
		return err
	}
	return movePlayerToRoom(world, room)
}

func whenPlayerMovesToGrenadeRoom(world *runtime.World, example map[string]string) error {
	room, err := intExample(example, "grenade_room")
	if err != nil {
		return err
	}
	return movePlayerToRoom(world, room)
}

func movePlayerToRoom(world *runtime.World, room int) error {
	result := gameFrom(world, "game").Move(room)
	world.State["move_result"] = result
	world.State["turn_messages"] = result.Messages
	world.State["action_taken"] = true
	return nil
}

func thenPlayerInToRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "to_room")
}

func thenPlayerInFromRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "from_room")
}

func thenPlayerInPlayerRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "player_room")
}

func thenPlayerInRelocationRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "relocation_room")
}

func thenPlayerInExampleRoom(world *runtime.World, example map[string]string, key string) error {
	want, err := intExample(example, key)
	if err != nil {
		return err
	}
	got := gameFrom(world, "game").Setup().Player
	if got != want {
		return fmt.Errorf("player room is %d, want %d", got, want)
	}
	return nil
}

func thenGameStillInProgress(world *runtime.World, _ map[string]string) error {
	return assertGameStatus(world, wumpus.StatusInProgress)
}

func thenGameStatus(world *runtime.World, example map[string]string) error {
	return assertGameStatus(world, wumpus.Status(example["game_status"]))
}

func thenPlayerLoses(world *runtime.World, _ map[string]string) error {
	return assertGameStatus(world, wumpus.StatusLost)
}

func assertGameStatus(world *runtime.World, want wumpus.Status) error {
	got := gameFrom(world, "game").Status()
	if got != want {
		return fmt.Errorf("game status is %s, want %s", got, want)
	}
	return nil
}

func thenTurnMessagesAre(world *runtime.World, example map[string]string) error {
	want := stringList(example["messages"])
	got := world.State["turn_messages"].([]string)
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("turn messages are %v, want %v", got, want)
	}
	return nil
}

func thenMoveRejectedWithMessage(world *runtime.World, example map[string]string) error {
	result := world.State["move_result"].(wumpus.MoveResult)
	if result.RejectedMessage != example["message"] {
		return fmt.Errorf("rejection message is %q, want %q", result.RejectedMessage, example["message"])
	}
	return nil
}

func givenNextBatRelocationRoom(world *runtime.World, example map[string]string) error {
	room, err := intAnyExample(example, "relocation_room", "wumpus_room")
	if err != nil {
		return err
	}
	gameFrom(world, "game").SetNextBatRelocation(room)
	return nil
}

func givenNextWumpusWakeChoice(world *runtime.World, example map[string]string) error {
	choice, err := parseWakeChoice(example["wake_choice"])
	if err != nil {
		return err
	}
	gameFrom(world, "game").SetNextWumpusWakeChoice(choice)
	return nil
}

func thenWumpusInRoom(world *runtime.World, example map[string]string) error {
	want, err := intExample(example, "expected_wumpus_room")
	if err != nil {
		return err
	}
	got := gameFrom(world, "game").Setup().Wumpus
	if got != want {
		return fmt.Errorf("Wumpus room is %d, want %d", got, want)
	}
	return nil
}

func whenTurnWarningsRequested(world *runtime.World, _ map[string]string) error {
	world.State["warnings"] = gameFrom(world, "game").TurnWarnings()
	return nil
}

func thenWarningMessagesAre(world *runtime.World, example map[string]string) error {
	got := world.State["warnings"].([]string)
	want := stringList(example["warnings"])
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("warning messages are %v, want %v", got, want)
	}
	return nil
}

func givenOrThenPlayerHasArrows(world *runtime.World, example map[string]string) error {
	arrows, err := intExample(example, "arrows")
	if err != nil {
		return err
	}
	if _, actionTaken := world.State["action_taken"]; !actionTaken {
		gameFrom(world, "game").SetArrows(arrows)
		return nil
	}
	return assertArrows(world, arrows)
}

func thenPlayerHasRemainingArrows(world *runtime.World, example map[string]string) error {
	arrows, err := intExample(example, "remaining_arrows")
	if err != nil {
		return err
	}
	return assertArrows(world, arrows)
}

func assertArrows(world *runtime.World, want int) error {
	got := gameFrom(world, "game").Arrows()
	if got != want {
		return fmt.Errorf("player has %d arrows, want %d", got, want)
	}
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

func thenShotRejectedWithMessage(world *runtime.World, example map[string]string) error {
	result := world.State["shoot_result"].(wumpus.ShootResult)
	if result.RejectedMessage != example["message"] {
		return fmt.Errorf("shot rejection message is %q, want %q", result.RejectedMessage, example["message"])
	}
	return nil
}

func givenInteractiveSetup(world *runtime.World, example map[string]string) error {
	if err := givenConfiguredSetup(world, example); err != nil {
		return err
	}
	arrows, err := intExample(example, "arrows")
	if err != nil {
		return err
	}
	game := *gameFrom(world, "game")
	game.SetArrows(arrows)
	session := interactive.NewSessionWithGame(game)
	world.State["session"] = session
	world.State["game"] = session.Game()
	return nil
}

func givenInteractiveSetupCarryingGrenade(world *runtime.World, example map[string]string) error {
	setup := wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}
	game, err := wumpus.NewGameWithSetup(setup)
	if err != nil {
		return err
	}
	arrows, err := intExample(example, "arrows")
	if err != nil {
		return err
	}
	game.SetArrows(arrows)
	game.GiveGrenade()
	session := interactive.NewSessionWithGame(game)
	world.State["session"] = session
	world.State["game"] = session.Game()
	return nil
}

func givenInteractiveSetupSeed(world *runtime.World, example map[string]string) error {
	seed, err := int64Example(example, "seed")
	if err != nil {
		return err
	}
	session := interactive.NewSessionWithSeed(seed)
	world.State["session"] = session
	world.State["game"] = session.Game()
	return nil
}

func givenNewInteractiveSession(world *runtime.World, _ map[string]string) error {
	world.State["session"] = interactive.NewSession()
	return nil
}

func whenNextTurnDisplayed(world *runtime.World, _ map[string]string) error {
	world.State["displayed_lines"] = sessionFrom(world).DisplayTurn()
	return nil
}

func whenPlayerEntersCommand(world *runtime.World, example map[string]string) error {
	lines := sessionFrom(world).EnterCommand(example["command"])
	world.State["displayed_lines"] = lines
	world.State["game"] = sessionFrom(world).Game()
	world.State["action_taken"] = true
	return nil
}

func thenDisplayedLinesAre(world *runtime.World, example map[string]string) error {
	got := world.State["displayed_lines"].([]string)
	want := stringList(example["lines"])
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("displayed lines are %v, want %v", got, want)
	}
	return nil
}

func thenDisplayedLinesInclude(world *runtime.World, example map[string]string) error {
	lines := world.State["displayed_lines"].([]string)
	expected := stringList(firstPresent(example, "message", "messages"))
	if prompt, ok := example["prompt"]; ok {
		expected = []string{prompt}
	}
	for _, want := range expected {
		if !containsString(lines, want) {
			return fmt.Errorf("displayed lines %v do not include %q", lines, want)
		}
	}
	return nil
}

func givenPlayerHasLost(world *runtime.World, _ map[string]string) error {
	session := sessionFrom(world)
	world.State["lost_setup"] = session.Game().Setup()
	grenadeRoom, grenadeOK := session.Game().GrenadeRoom()
	pendingRoom, pendingOK := session.Game().PendingGrenadeRoom()
	world.State["lost_grenade_room"] = grenadeRoom
	world.State["lost_grenade_ok"] = grenadeOK
	world.State["lost_pending_room"] = pendingRoom
	world.State["lost_pending_ok"] = pendingOK
	session.MarkLostForTest()
	return nil
}

func whenPlayerAnswersSameSetupPrompt(world *runtime.World, example map[string]string) error {
	sessionFrom(world).AnswerSameSetup(example["answer"])
	return nil
}

func thenNextGameSetupRelation(world *runtime.World, example map[string]string) error {
	lost := world.State["lost_setup"].(wumpus.Setup)
	current := sessionFrom(world).Game().Setup()
	identical := reflect.DeepEqual(lost, current)
	if example["setup_relation"] == "identical" && !identical {
		return fmt.Errorf("setup is %v, want identical to %v", current, lost)
	}
	if example["setup_relation"] == "different" && identical {
		return fmt.Errorf("setup is identical to lost setup %v", lost)
	}
	return nil
}

func whenPlayerAnswersInstructionsPrompt(world *runtime.World, example map[string]string) error {
	world.State["displayed_lines"] = sessionFrom(world).AnswerInstructions(example["answer"])
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
	if contains(gameFrom(world, "game").Setup().OccupiedRooms(), room) {
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
		if containsString(messages, "YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!") {
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
	if room, ok := gameFrom(world, "game").GrenadeRoom(); ok {
		return fmt.Errorf("unclaimed grenade remains in room %d", room)
	}
	return nil
}

func thenNoGrenadePending(world *runtime.World, _ map[string]string) error {
	if room, ok := gameFrom(world, "game").PendingGrenadeRoom(); ok {
		return fmt.Errorf("grenade pending in room %d", room)
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
	want, err := roomList(example["remaining_bat_rooms"])
	if err != nil {
		return err
	}
	got := gameFrom(world, "game").Setup().Bats
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("bat rooms = %v, want %v", got, want)
	}
	return nil
}

func thenPitRooms(world *runtime.World, example map[string]string) error {
	want, err := roomList(example["pit_rooms"])
	if err != nil {
		return err
	}
	got := gameFrom(world, "game").Setup().Pits
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("pit rooms = %v, want %v", got, want)
	}
	return nil
}

func thenReplaySetupIncludingGrenadeIdentical(world *runtime.World, _ map[string]string) error {
	session := sessionFrom(world)
	if !reflect.DeepEqual(session.Game().Setup(), world.State["lost_setup"].(wumpus.Setup)) {
		return fmt.Errorf("replay setup differs")
	}
	room, ok := session.Game().GrenadeRoom()
	if ok != world.State["lost_grenade_ok"].(bool) || room != world.State["lost_grenade_room"].(int) {
		return fmt.Errorf("replay grenade room = %d/%v", room, ok)
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

func givenNextSleepyWumpusObservation(world *runtime.World, example map[string]string) error {
	switch example["sleepy_observation"] {
	case "asleep":
		gameFrom(world, "game").SetNextSleepyWumpusObservation(true)
	case "awake":
		gameFrom(world, "game").SetNextSleepyWumpusObservation(false)
	default:
		return fmt.Errorf("unsupported sleepy observation %q", example["sleepy_observation"])
	}
	return nil
}

func givenWumpusAsleep(world *runtime.World, _ map[string]string) error {
	gameFrom(world, "game").SetWumpusAsleep(true)
	return nil
}

func givenWumpusAwake(world *runtime.World, _ map[string]string) error {
	gameFrom(world, "game").SetWumpusAsleep(false)
	return nil
}

func thenWumpusSleepState(world *runtime.World, example map[string]string) error {
	wantAsleep := example["sleep_state"] == "asleep"
	if example["sleep_state"] != "asleep" && example["sleep_state"] != "awake" {
		return fmt.Errorf("unsupported sleep state %q", example["sleep_state"])
	}
	if got := gameFrom(world, "game").WumpusAsleep(); got != wantAsleep {
		return fmt.Errorf("Wumpus asleep = %v, want %v", got, wantAsleep)
	}
	return nil
}

func thenWumpusSleepStateAsleep(world *runtime.World, _ map[string]string) error {
	return requireWumpusSleepState(world, true)
}

func thenWumpusSleepStateAwake(world *runtime.World, _ map[string]string) error {
	return requireWumpusSleepState(world, false)
}

func requireWumpusSleepState(world *runtime.World, wantAsleep bool) error {
	if got := gameFrom(world, "game").WumpusAsleep(); got != wantAsleep {
		return fmt.Errorf("Wumpus asleep = %v, want %v", got, wantAsleep)
	}
	return nil
}

func thenTurnWarningsAre(world *runtime.World, example map[string]string) error {
	want := stringList(example["warnings"])
	got := gameFrom(world, "game").TurnWarnings()
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("turn warnings = %v, want %v", got, want)
	}
	return nil
}

func givenNextSleepingWumpusEntryOutcome(world *runtime.World, example map[string]string) error {
	outcome := wumpus.SleepingWumpusEntryOutcome(example["entry_outcome"])
	switch outcome {
	case wumpus.SleepingWumpusWakes, wumpus.SleepingWumpusStaysAsleep:
		gameFrom(world, "game").SetNextSleepingWumpusEntryOutcome(outcome)
		return nil
	default:
		return fmt.Errorf("unsupported sleeping Wumpus entry outcome %q", example["entry_outcome"])
	}
}

func givenPlayerHasSeenSleepingWumpus(world *runtime.World, _ map[string]string) error {
	gameFrom(world, "game").SetSawSleepingWumpus(true)
	return nil
}

func whenBothGamesObserveSleepyWumpus(world *runtime.World, example map[string]string) error {
	turnCount, err := intExample(example, "turn_count")
	if err != nil {
		return err
	}
	world.State["sleepy_observations"] = gameFrom(world, "game").ObserveSleepyWumpusBehavior(turnCount)
	world.State["another_sleepy_observations"] = gameFrom(world, "another_game").ObserveSleepyWumpusBehavior(turnCount)
	return nil
}

func thenBothSleepyObservationsIdentical(world *runtime.World, _ map[string]string) error {
	left := world.State["sleepy_observations"].([]string)
	right := world.State["another_sleepy_observations"].([]string)
	if !reflect.DeepEqual(left, right) {
		return fmt.Errorf("sleepy observations differ: %v and %v", left, right)
	}
	return nil
}

func setupSnapshot(setup wumpus.Setup) inspectedSetup {
	pits := append([]int(nil), setup.Pits...)
	bats := append([]int(nil), setup.Bats...)
	sort.Ints(pits)
	sort.Ints(bats)
	return inspectedSetup{Player: setup.Player, Wumpus: setup.Wumpus, Pits: pits, Bats: bats}
}

func gameFrom(world *runtime.World, key string) *wumpus.Game {
	return world.State[key].(*wumpus.Game)
}

func sessionFrom(world *runtime.World) *interactive.Session {
	return world.State["session"].(*interactive.Session)
}

func firstPresent(example map[string]string, keys ...string) string {
	for _, key := range keys {
		if value, ok := example[key]; ok {
			return value
		}
	}
	return ""
}

func intAnyExample(example map[string]string, keys ...string) (int, error) {
	for _, key := range keys {
		if _, ok := example[key]; ok {
			return intExample(example, key)
		}
	}
	return 0, fmt.Errorf("missing any example key %v", keys)
}

func intExample(example map[string]string, key string) (int, error) {
	value, ok := example[key]
	if !ok {
		return 0, fmt.Errorf("missing example %q", key)
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q for %s", value, key)
	}
	return parsed, nil
}

func parseWakeChoice(value string) (wumpus.WumpusWakeChoice, error) {
	if value == "stay" {
		return wumpus.WumpusWakeChoice{Stay: true}, nil
	}
	const prefix = "move to "
	if strings.HasPrefix(value, prefix) {
		room, err := strconv.Atoi(strings.TrimPrefix(value, prefix))
		if err != nil {
			return wumpus.WumpusWakeChoice{}, fmt.Errorf("invalid wake choice %q", value)
		}
		return wumpus.WumpusWakeChoice{Destination: room}, nil
	}
	return wumpus.WumpusWakeChoice{}, fmt.Errorf("unsupported wake choice %q", value)
}

func int64Example(example map[string]string, key string) (int64, error) {
	value, ok := example[key]
	if !ok {
		return 0, fmt.Errorf("missing example %q", key)
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q for %s", value, key)
	}
	return parsed, nil
}

func roomList(value string) ([]int, error) {
	var rooms []int
	for _, part := range strings.Split(value, ",") {
		room, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid room list %q", value)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func optionalRoomList(value string) ([]int, error) {
	if value == "none" {
		return nil, nil
	}
	return roomList(value)
}

func stringList(value string) []string {
	if value == "none" {
		return nil
	}
	var values []string
	for _, part := range strings.Split(value, ",") {
		values = append(values, strings.TrimSpace(part))
	}
	return values
}

func hazardList(value string) []string {
	if value == "none" {
		return nil
	}
	var hazards []string
	for _, part := range strings.Split(value, ",") {
		hazards = append(hazards, strings.TrimSpace(part))
	}
	return hazards
}

func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func firstUnoccupiedRoom(occupied ...any) int {
	seen := map[int]bool{}
	for _, value := range occupied {
		switch rooms := value.(type) {
		case int:
			seen[rooms] = true
		case []int:
			for _, room := range rooms {
				seen[room] = true
			}
		}
	}
	for room := 1; room <= 20; room++ {
		if !seen[room] {
			return room
		}
	}
	return 1
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
