package steps

import (
	"fmt"
	"reflect"
	"slices"
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
		"exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats":       thenDistinctOccupiedCount,
		"exactly 6 distinct rooms are occupied by the player, Wumpus, pits, and bats":                      thenSixDistinctOccupiedRooms,
		"another new game created with seed 1973":                                                          givenAnotherNewGameSeed1973,
		"another new game created with seed <seed>":                                                        givenAnotherNewGameSeed,
		"both setups are inspected":                                                                        whenBothSetupsInspected,
		"both setups have identical player, Wumpus, pit, and bat rooms":                                    thenBothSetupsIdentical,
		"a completed game created with seed 1973":                                                          givenCompletedGameSeed1973,
		"a completed game created with seed <seed>":                                                        givenCompletedGameSeed,
		"a same setup replay is started":                                                                   whenSameSetupReplayStarted,
		"the replay setup has identical player, Wumpus, pit, and bat rooms":                                thenReplaySetupIdentical,
		"the player moves to room <to_room>":                                                               whenPlayerMoves,
		"the player moves to room <pit_room>":                                                              whenPlayerMoves,
		"the player moves to room <bat_room>":                                                              whenPlayerMoves,
		"the player moves to room <wumpus_room>":                                                           whenPlayerMoves,
		"the player moves to room <grenade_room>":                                                          whenPlayerMovesToGrenadeRoom,
		"the player is in room <to_room>":                                                                  thenPlayerInToRoom,
		"the player is in room <from_room>":                                                                thenPlayerInFromRoom,
		"the player is in room <player_room>":                                                              thenPlayerInPlayerRoom,
		"the player is in room <expected_player_room>":                                                     thenPlayerInExpectedPlayerRoom,
		"the player is in room <relocation_room>":                                                          thenPlayerInRelocationRoom,
		"the game is still in progress":                                                                    thenGameStillInProgress,
		"the game is in progress":                                                                          thenGameStillInProgress,
		"the game is <game_status>":                                                                        thenGameStatus,
		"the player loses":                                                                                 thenPlayerLoses,
		"the turn messages are <messages>":                                                                 thenTurnMessagesAre,
		"the move is rejected with message <message>":                                                      thenMoveRejectedWithMessage,
		"the next bat relocation room is <relocation_room>":                                                givenNextBatRelocationRoom,
		"the next bat relocation room is <wumpus_room>":                                                    givenNextBatRelocationRoom,
		"the next Wumpus wake choice is <wake_choice>":                                                     givenNextWumpusWakeChoice,
		"the Wumpus is in room <expected_wumpus_room>":                                                     thenWumpusInRoom,
		"the Wumpus is in room <wumpus_room>":                                                              thenWumpusInWumpusRoom,
		"the turn warnings are requested":                                                                  whenTurnWarningsRequested,
		"the warning messages are <warnings>":                                                              thenWarningMessagesAre,
		"a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>": givenShootingSetup,
		"the player starts with <initial_arrows> arrows":                                                   givenPlayerStartsWithArrows,
		"the player has <arrows> arrows":                                                                   givenOrThenPlayerHasArrows,
		"the player has <remaining_arrows> arrows":                                                         thenPlayerHasRemainingArrows,
		"the next arrow deviation room is <deviation_room>":                                                givenNextArrowDeviationRoom,
		"the player shoots the path <path>":                                                                whenPlayerShootsPath,
		"the player wins":                                                                                  thenPlayerWins,
		"the arrow traversed rooms are <traversed_rooms>":                                                  thenArrowTraversedRoomsAre,
		"the requested shot path is <expected_path>":                                                       thenRequestedShotPathIs,
		"the shot is rejected with message <message>":                                                      thenShotRejectedWithMessage,
		"an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows": givenInteractiveSetup,
		"an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows":   givenInteractiveSetup,
		"an interactive game setup with seed <seed>":                                      givenInteractiveSetupSeed,
		"a new interactive session":                                                       givenNewInteractiveSession,
		"the next turn is displayed":                                                      whenNextTurnDisplayed,
		"the player enters command <command>":                                             whenPlayerEntersCommand,
		"the displayed lines are <lines>":                                                 thenDisplayedLinesAre,
		"the displayed lines include <message>":                                           thenDisplayedLinesInclude,
		"the displayed lines include <messages>":                                          thenDisplayedLinesInclude,
		"the displayed lines include <warnings>":                                          thenDisplayedLinesInclude,
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
		"the player has seen the sleeping Wumpus shape":                      givenPlayerHasSeenSleepingWumpus,
		"both games observe sleepy Wumpus behavior for <turn_count> turns":   whenBothGamesObserveSleepyWumpus,
		"both games produce identical sleepy Wumpus observations":            thenBothSleepyObservationsIdentical,
		"the next jumping Wumpus turn event is <jump_event>":                 givenNextJumpingWumpusTurnEvent,
		"the next jumping Wumpus turn event is jumps":                        givenNextJumpingWumpusTurnEventJumps,
		"the next jumping Wumpus turn event is no jump":                      givenNextJumpingWumpusTurnEventNoJump,
		"the next Wumpus jump path is <jump_path>":                           givenNextWumpusJumpPath,
		"the next first jump player landing outcome is <landing_outcome>":    givenNextFirstJumpLandingOutcome,
		"the next turn begins":                                               whenNextTurnBegins,
		"the displayed lines do not include <message>":                       thenDisplayedLinesDoNotInclude,
		"the player may take the next command":                               thenPlayerMayTakeNextCommand,
		"every Wumpus jump segment is a legal tunnel":                        thenEveryWumpusJumpSegmentLegal,
		"both games evaluate jumping Wumpus behavior for <turn_count> turns": whenBothGamesEvaluateJumpingWumpus,
		"both games produce identical jumping Wumpus events":                 thenBothJumpEventsIdentical,
		"the turn count starts at <turn_count>":                              givenTurnCount,
		"the turn count is <turn_count>":                                     givenOrThenTurnCount,
		"the turn count is <expected_turn_count>":                            thenExpectedTurnCount,
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
	return assertRoomList("exits", world.State["exits"].([]int), example["exits"])
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
		if !slices.Contains(reachable, room) {
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
		if slices.Contains(exits, room) {
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
	setup, err := setupFromExample(example)
	if err != nil {
		return err
	}
	return setConfiguredSetup(world, setup)
}

func setupFromExample(example map[string]string) (wumpus.Setup, error) {
	player, err := intAnyExample(example, "player_room", "from_room")
	if err != nil {
		return wumpus.Setup{}, err
	}
	wumpusRoom, err := intExample(example, "wumpus_room")
	if err != nil {
		return wumpus.Setup{}, err
	}
	pits, err := roomList(example["pit_rooms"])
	if err != nil {
		return wumpus.Setup{}, err
	}
	bats, err := roomList(example["bat_rooms"])
	if err != nil {
		return wumpus.Setup{}, err
	}
	return wumpus.Setup{Player: player, Wumpus: wumpusRoom, Pits: pits, Bats: bats}, nil
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
	setup, err := setupFromExample(map[string]string{
		"player_room": example["wumpus_room"],
		"wumpus_room": example["wumpus_room"],
		"pit_rooms":   example["pit_rooms"],
		"bat_rooms":   example["bat_rooms"],
	})
	if err != nil {
		return err
	}
	return setConfiguredSetup(world, setup)
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
	inspectSetup(world, "inspected_setup", "game")
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
	distinct := map[int]struct{}{}
	for _, room := range rooms {
		distinct[room] = struct{}{}
	}
	if len(distinct) != want {
		return fmt.Errorf("got %d distinct occupied rooms from %v, want %d", len(distinct), rooms, want)
	}
	return nil
}

func whenBothSetupsInspected(world *runtime.World, _ map[string]string) error {
	inspectSetup(world, "inspected_setup", "game")
	inspectSetup(world, "another_inspected_setup", "another_game")
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
	return movePlayerToAnyExampleRoom(world, example, "to_room", "pit_room", "bat_room", "wumpus_room")
}

func movePlayerToAnyExampleRoom(world *runtime.World, example map[string]string, keys ...string) error {
	room, err := intAnyExample(example, keys...)
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

func thenPlayerInExpectedPlayerRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "expected_player_room")
}

func thenPlayerInRelocationRoom(world *runtime.World, example map[string]string) error {
	return thenPlayerInExampleRoom(world, example, "relocation_room")
}

func thenPlayerInExampleRoom(world *runtime.World, example map[string]string, key string) error {
	return thenSetupRoom(world, example, key, "player", setupPlayerRoom)
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
	return assertStringState(world, "turn_messages", "turn messages", example["messages"])
}

func thenMoveRejectedWithMessage(world *runtime.World, example map[string]string) error {
	result := world.State["move_result"].(wumpus.MoveResult)
	return assertRejectedMessage("move rejection", result.RejectedMessage, example["message"])
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
	return thenWumpusInExampleRoom(world, example, "expected_wumpus_room")
}

func thenWumpusInWumpusRoom(world *runtime.World, example map[string]string) error {
	return thenWumpusInExampleRoom(world, example, "wumpus_room")
}

func thenWumpusInExampleRoom(world *runtime.World, example map[string]string, key string) error {
	return thenSetupRoom(world, example, key, "Wumpus", setupWumpusRoom)
}

func thenSetupRoom(world *runtime.World, example map[string]string, key, label string, roomOf func(wumpus.Setup) int) error {
	want, err := intExample(example, key)
	if err != nil {
		return err
	}
	got := roomOf(gameFrom(world, "game").Setup())
	if got != want {
		return fmt.Errorf("%s room is %d, want %d", label, got, want)
	}
	return nil
}

func setupPlayerRoom(setup wumpus.Setup) int {
	return setup.Player
}

func setupWumpusRoom(setup wumpus.Setup) int {
	return setup.Wumpus
}

func whenTurnWarningsRequested(world *runtime.World, _ map[string]string) error {
	world.State["warnings"] = gameFrom(world, "game").TurnWarnings()
	return nil
}

func thenWarningMessagesAre(world *runtime.World, example map[string]string) error {
	return assertStringState(world, "warnings", "warning messages", example["warnings"])
}

func givenOrThenPlayerHasArrows(world *runtime.World, example map[string]string) error {
	return setOrAssertParsedInt(world, example, "arrows", func(arrows int) {
		gameFrom(world, "game").SetArrows(arrows)
	}, func(arrows int) error {
		return assertArrows(world, arrows)
	})
}

func setOrAssertParsedInt(world *runtime.World, example map[string]string, key string, set func(int), assert func(int) error) error {
	value, err := intExample(example, key)
	if err != nil {
		return err
	}
	if _, actionTaken := world.State["action_taken"]; !actionTaken {
		set(value)
		return nil
	}
	return assert(value)
}

func setStringChoice(value, label string, allowed []string, set func(string)) error {
	if slices.Contains(allowed, value) {
		set(value)
		return nil
	}
	return fmt.Errorf("unsupported %s %q", label, value)
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

func assertRejectedMessage(label, got, want string) error {
	if got != want {
		return fmt.Errorf("%s message is %q, want %q", label, got, want)
	}
	return nil
}

func assertRoomList(label string, got []int, wantValue string) error {
	return assertParsedRoomListValue(label, got, wantValue, roomList)
}

func assertOptionalRoomList(label string, got []int, wantValue string) error {
	return assertParsedRoomListValue(label, got, wantValue, optionalRoomList)
}

func assertParsedRoomListValue(label string, got []int, wantValue string, parse func(string) ([]int, error)) error {
	want, err := parse(wantValue)
	if err != nil {
		return err
	}
	return assertParsedRoomList(label, got, want)
}

func assertParsedRoomList(label string, got, want []int) error {
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("%s %v, want %v", label, got, want)
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
	storeSession(world, interactive.NewSessionWithGame(game))
	return nil
}

func givenInteractiveSetupSeed(world *runtime.World, example map[string]string) error {
	seed, err := int64Example(example, "seed")
	if err != nil {
		return err
	}
	storeSession(world, interactive.NewSessionWithSeed(seed))
	return nil
}

func givenNewInteractiveSession(world *runtime.World, _ map[string]string) error {
	storeSession(world, interactive.NewSession())
	return nil
}

func whenNextTurnDisplayed(world *runtime.World, _ map[string]string) error {
	world.State["displayed_lines"] = sessionFrom(world).DisplayTurn()
	return nil
}

func whenPlayerEntersCommand(world *runtime.World, example map[string]string) error {
	lines := sessionFrom(world).EnterCommand(example["command"])
	world.State["displayed_lines"] = lines
	world.State["turn_messages"] = lines
	world.State["game"] = sessionFrom(world).Game()
	world.State["action_taken"] = true
	return nil
}

func thenDisplayedLinesAre(world *runtime.World, example map[string]string) error {
	return assertStringState(world, "displayed_lines", "displayed lines", example["lines"])
}

func assertStringState(world *runtime.World, key, label, wantValue string) error {
	return assertStringList(label, world.State[key].([]string), wantValue)
}

func assertStringList(label string, got []string, wantValue string) error {
	want := stringList(wantValue)
	if !slices.Equal(got, want) {
		return fmt.Errorf("%s are %v, want %v", label, got, want)
	}
	return nil
}

func thenDisplayedLinesInclude(world *runtime.World, example map[string]string) error {
	lines := world.State["displayed_lines"].([]string)
	value := firstPresent(example, "message", "messages", "warnings", "prompt")
	expected := stringList(value)
	if strings.HasPrefix(value, "SHOOT") {
		expected = []string{value}
	}
	for _, want := range expected {
		if !slices.Contains(lines, want) {
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
	return setStringChoice(example["entry_outcome"], "sleeping Wumpus entry outcome", []string{
		string(wumpus.SleepingWumpusWakes),
		string(wumpus.SleepingWumpusStaysAsleep),
	}, func(value string) {
		gameFrom(world, "game").SetNextSleepingWumpusEntryOutcome(wumpus.SleepingWumpusEntryOutcome(value))
	})
}

func givenPlayerHasSeenSleepingWumpus(world *runtime.World, _ map[string]string) error {
	gameFrom(world, "game").SetSawSleepingWumpus(true)
	return nil
}

func whenBothGamesObserveSleepyWumpus(world *runtime.World, example map[string]string) error {
	return recordPairedStringObservations(world, example, "sleepy_observations", "another_sleepy_observations", func(game *wumpus.Game, turnCount int) []string {
		return game.ObserveSleepyWumpusBehavior(turnCount)
	})
}

func thenBothSleepyObservationsIdentical(world *runtime.World, _ map[string]string) error {
	return assertStringObservationsIdentical(world, "sleepy observations", "sleepy_observations", "another_sleepy_observations")
}

func recordPairedStringObservations(world *runtime.World, example map[string]string, leftKey, rightKey string, observe func(*wumpus.Game, int) []string) error {
	turnCount, err := intExample(example, "turn_count")
	if err != nil {
		return err
	}
	world.State[leftKey] = observe(gameFrom(world, "game"), turnCount)
	world.State[rightKey] = observe(gameFrom(world, "another_game"), turnCount)
	return nil
}

func assertStringObservationsIdentical(world *runtime.World, label, leftKey, rightKey string) error {
	left := world.State[leftKey].([]string)
	right := world.State[rightKey].([]string)
	if !reflect.DeepEqual(left, right) {
		return fmt.Errorf("%s differ: %v and %v", label, left, right)
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

func inspectSetup(world *runtime.World, stateKey, gameKey string) {
	world.State[stateKey] = setupSnapshot(gameFrom(world, gameKey).Setup())
}

func storeSession(world *runtime.World, session *interactive.Session) {
	world.State["session"] = session
	world.State["game"] = session.Game()
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
	const whumpToken = "__WHUMP_MESSAGE__"
	const armedPromptToken = "__ARMED_PROMPT__"
	value = strings.ReplaceAll(value, "YOU HEAR WHUMP, WHUMP.", whumpToken)
	value = strings.ReplaceAll(value, "SHOOT, MOVE OR THROW (S-M-T)?", armedPromptToken)
	values := commaSeparatedStrings(value)
	for index, item := range values {
		values[index] = strings.ReplaceAll(item, whumpToken, "YOU HEAR WHUMP, WHUMP.")
		values[index] = strings.ReplaceAll(values[index], armedPromptToken, "SHOOT, MOVE OR THROW (S-M-T)?")
	}
	return values
}

func hazardList(value string) []string {
	return commaSeparatedStrings(value)
}

func commaSeparatedStrings(value string) []string {
	if value == "none" {
		return nil
	}
	var values []string
	for _, part := range strings.Split(value, ",") {
		values = append(values, strings.TrimSpace(part))
	}
	return values
}

func firstUnoccupiedRoom(wumpusRoom int, pits, bats []int) int {
	occupied := append([]int{wumpusRoom}, pits...)
	occupied = append(occupied, bats...)
	for room := 1; room <= 20; room++ {
		if !slices.Contains(occupied, room) {
			return room
		}
	}
	return 1
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T10:59:48-05:00","module_hash":"66dea6db3459ebc3ec2bd2447f1b912e4437e221ce624ffbbc4b841c6131d9d7","functions":[{"id":"func/NewHandlers","name":"NewHandlers","line":23,"end_line":158,"hash":"de654007bb7cfeafb2f1943b0dcda70c9cea0dd6b86fd50e46ae309bb04f67de"},{"id":"func/givenNewCave","name":"givenNewCave","line":160,"end_line":163,"hash":"b303885e9253964f1c89c93d452cf776e8e66c22e02159b1600d3b23d5355e19"},{"id":"func/whenExitsQueried","name":"whenExitsQueried","line":165,"end_line":177,"hash":"72829feb5fdf15929cf91dab92e9b7205ab21e530269d7a28d37bd58c4e48cc6"},{"id":"func/thenExitsAre","name":"thenExitsAre","line":179,"end_line":181,"hash":"54a51c9280ca1917c76b6b50556422bb3d4b1237a6efaee38d59370bec41326e"},{"id":"func/thenExitCountIs","name":"thenExitCountIs","line":183,"end_line":193,"hash":"1e85d5799c79c7af3efa7ee78bf2dc5b093c508518041a1d5a6ebfa830409f7f"},{"id":"func/whenCaveTraversed","name":"whenCaveTraversed","line":195,"end_line":198,"hash":"1da77aa27b97820fcb189f839fa21cf56d27a57be8eea178352c500a81aa1dc9"},{"id":"func/thenEveryRoomReachable","name":"thenEveryRoomReachable","line":200,"end_line":211,"hash":"0b3d230a1f7911643c77062ade9af8d7ae4806d6a12d96831b4bc46bf3d07995"},{"id":"func/whenCaveInvariantsInspected","name":"whenCaveInvariantsInspected","line":213,"end_line":225,"hash":"56b296da689d0bbaaf3c4a7d32548d9e8a86036c140ce528d098cbc50651c253"},{"id":"func/thenEveryRoomHasThreeExits","name":"thenEveryRoomHasThreeExits","line":227,"end_line":234,"hash":"0a9569a5a84bd8663f7331f7503ca19ba137b2eaf0419af9eacc1a15d7d673d8"},{"id":"func/thenNoRoomIsItsOwnExit","name":"thenNoRoomIsItsOwnExit","line":236,"end_line":243,"hash":"b34289e335d77d52f30c512171e3b2009f38560eff0951581cf84770c8e89392"},{"id":"func/whenTunnelQueried","name":"whenTunnelQueried","line":245,"end_line":259,"hash":"5b59358ccb11b0ebe29e77e8923cd0963253509107025c89218d4b411cbf872e"},{"id":"func/thenReverseTunnelExists","name":"thenReverseTunnelExists","line":261,"end_line":272,"hash":"f5f67c27db51fc2cd84c2bb11b78c328558e1edf568faf33411addfdeebe8967"},{"id":"func/thenRoomIsNotExit","name":"thenRoomIsNotExit","line":274,"end_line":285,"hash":"e9544d7b105d91000287651b70364b8245909ae0f27c7ce74fb0107e61c6bf5b"},{"id":"func/givenConfiguredSetup","name":"givenConfiguredSetup","line":287,"end_line":293,"hash":"b365e6aaeef6081cc1b6514314ae24974b1a533ee09db54499ca05fb3854b685"},{"id":"func/setupFromExample","name":"setupFromExample","line":295,"end_line":313,"hash":"d020a420000e0e55d93ce75381d92f6faabb81e8e103469be447180650f48749"},{"id":"func/givenSetupRoom1WumpusBats","name":"givenSetupRoom1WumpusBats","line":315,"end_line":317,"hash":"48947502a897bcc1f6b0a59903efa01e0ed9b9474dfc5bd7d962ca60d94eb06c"},{"id":"func/givenSetupRoom10WumpusPit","name":"givenSetupRoom10WumpusPit","line":319,"end_line":321,"hash":"553d1e6c8444552ee0975195ed8019f91a7ff7e923308ab726a12174032b0031"},{"id":"func/givenSetupRoom6NoHazards","name":"givenSetupRoom6NoHazards","line":323,"end_line":325,"hash":"efc60680c5f79c09354efb2436a2b43564879bfbec94b88c045c58814cc43eca"},{"id":"func/setConfiguredSetup","name":"setConfiguredSetup","line":327,"end_line":340,"hash":"b7bb005b5fa16382036aac5238cd696cedf84d7e5ef5061baedf6e9625ef6a25"},{"id":"func/givenConfiguredSetupPlayerWithWumpus","name":"givenConfiguredSetupPlayerWithWumpus","line":342,"end_line":353,"hash":"bd505da8e6a4c9e5d112c7f76cfe693f3fcd1f167995f5b34733a0c5f47ccc9b"},{"id":"func/whenAdjacentHazardsQueried","name":"whenAdjacentHazardsQueried","line":355,"end_line":359,"hash":"ad8ebc9034df26f198e45997834455f89c9b859850758a4782e2d67e503175bb"},{"id":"func/thenAdjacentHazardsAre","name":"thenAdjacentHazardsAre","line":361,"end_line":363,"hash":"82f0d8fe6f49bca84b58910f49f4108676c4d95ece6c11ba5d06d1355a459077"},{"id":"func/thenAdjacentHazardsWumpusBats","name":"thenAdjacentHazardsWumpusBats","line":365,"end_line":367,"hash":"252f9b4cb039c099264b57fd0c645069edeb8faf113d189fedfa75a44cf38f87"},{"id":"func/thenAdjacentHazardsWumpusPit","name":"thenAdjacentHazardsWumpusPit","line":369,"end_line":371,"hash":"7c93f88965ba15dc3e90b05ef059c57dcdfd704e2789ba6cbe9866981cf1a510"},{"id":"func/thenNoAdjacentHazards","name":"thenNoAdjacentHazards","line":373,"end_line":375,"hash":"4a2373582f346e206232f7b489beb792647c840aee2c01cbf935f713bfe797dd"},{"id":"func/requireHazards","name":"requireHazards","line":377,"end_line":387,"hash":"ad745660fd99798be8c557c0f58a151dc26df92f0d2b445dccf5c38a13b8a26d"},{"id":"func/givenNewGameSeed1973","name":"givenNewGameSeed1973","line":389,"end_line":391,"hash":"697a8abcf485d2adc597262328e2c4ffdf487d0176f24ff3d9d977bea0be314f"},{"id":"func/givenNewGameSeed","name":"givenNewGameSeed","line":393,"end_line":395,"hash":"1e85dd442f94cb2f603fad5c66c32cfb37babf4cb66180d8f754d76cf0df0f4e"},{"id":"func/givenAnotherNewGameSeed","name":"givenAnotherNewGameSeed","line":397,"end_line":399,"hash":"b67ad35f9532f1d2f16e417075a93381fc8e9b833a14ca4afc26227ec0b289e4"},{"id":"func/givenAnotherNewGameSeed1973","name":"givenAnotherNewGameSeed1973","line":401,"end_line":403,"hash":"28ef48e76e29f95355dc54de7114d6901ede30c3a58d4c0a846c0accd06bee99"},{"id":"func/givenCompletedGameSeed","name":"givenCompletedGameSeed","line":405,"end_line":407,"hash":"d2d983b667fa0ea6c2963f0b09df0044a1788009d635a47678e0026abe4ec338"},{"id":"func/givenCompletedGameSeed1973","name":"givenCompletedGameSeed1973","line":409,"end_line":411,"hash":"77a7b094d27dca37636f518b8334fd12e9f66a0039b045662629d686d5852085"},{"id":"func/setGameFromExample","name":"setGameFromExample","line":413,"end_line":419,"hash":"11194e6924ff1701273368a8a677c9605beab8f5005bd4dca265ca7b55d698ab"},{"id":"func/setGame","name":"setGame","line":421,"end_line":428,"hash":"bfea1554f7c90fb5d262dff95ffb51daffd63fdd49a94738145e7eced8fa9d40"},{"id":"func/whenSetupInspected","name":"whenSetupInspected","line":430,"end_line":433,"hash":"57cf8cda29c8633354437e09b8c34030b60ec408b7df97cd97b8c5c3d39e2377"},{"id":"func/thenOnePlayer","name":"thenOnePlayer","line":435,"end_line":437,"hash":"0109f2fb9e21b7264aeb2f9349c5a14cbca3fe17019b3cda58af3f4b879d7901"},{"id":"func/thenOneWumpus","name":"thenOneWumpus","line":439,"end_line":441,"hash":"4fbcd442325d17c16906bedc17c638b7870dd4209a4848f27262aa213e89efbf"},{"id":"func/thenTwoPits","name":"thenTwoPits","line":443,"end_line":445,"hash":"8123a9d5ac9b139683838d3a146d487a5dd2a702abb5669df66a9dca4f0ffdc7"},{"id":"func/thenTwoBats","name":"thenTwoBats","line":447,"end_line":449,"hash":"81bff643b7f216de15a07681ac1ee51d587c25f1437ab824b68f9943e1b406a8"},{"id":"func/inspected","name":"inspected","line":451,"end_line":453,"hash":"1c34327b0fec611fdbf723d8fb4d2b129b02d577c466fb1bb5affd23316435be"},{"id":"func/requirePlaced","name":"requirePlaced","line":455,"end_line":460,"hash":"6e17aa67b6a61da1599cc1e174aaf9c9da4c6dbb1699ab1b1d6dc0529d61ed32"},{"id":"func/requireCount","name":"requireCount","line":462,"end_line":467,"hash":"6d517c6577380a1185f0b32a3b65cc4c38a699cbe8449f1f6b8bdf8c02b7e68d"},{"id":"func/whenOccupiedRoomsInspected","name":"whenOccupiedRoomsInspected","line":469,"end_line":473,"hash":"dde907045be9e972d3c00d8a23dd70ebfb4ccdddc7cc107af584334a72c1aff4"},{"id":"func/thenOccupiedRoomsValid","name":"thenOccupiedRoomsValid","line":475,"end_line":482,"hash":"435e2417a03f5478138803a4f67dac7639e11164e388b01d0b170ea076cf7c31"},{"id":"func/thenDistinctOccupiedCount","name":"thenDistinctOccupiedCount","line":484,"end_line":490,"hash":"98e95d667992aa5d7a3b966b3f60ad09baea33903f8178fb9cf9c914a4c88eeb"},{"id":"func/thenSixDistinctOccupiedRooms","name":"thenSixDistinctOccupiedRooms","line":492,"end_line":494,"hash":"e112385db267657bdeffe80bd15a608ef044015013a53631c2dfd651b288317f"},{"id":"func/requireDistinctOccupiedCount","name":"requireDistinctOccupiedCount","line":496,"end_line":506,"hash":"6dc23f0ef66feff258042cb7443f8fa7d13ad6c71f03dd09ef4843880800eaa0"},{"id":"func/whenBothSetupsInspected","name":"whenBothSetupsInspected","line":508,"end_line":512,"hash":"bf52f279b8e8976e51a9c7d694bb9004bf8aaaaa78b077a26ee16f054aebff26"},{"id":"func/thenBothSetupsIdentical","name":"thenBothSetupsIdentical","line":514,"end_line":521,"hash":"ae04a139983d0e373b057cd30959090aeb807983cb60b983b40350583f17eeed"},{"id":"func/whenSameSetupReplayStarted","name":"whenSameSetupReplayStarted","line":523,"end_line":527,"hash":"ad655d7c84d7eaa890c24817e1fea0743c2c1e54f1a38e213570c3964738d414"},{"id":"func/thenReplaySetupIdentical","name":"thenReplaySetupIdentical","line":529,"end_line":536,"hash":"4d9a144a7e315ecaabf64071c508b392e9d8c90aa44506d42b9374b52a473db6"},{"id":"func/whenPlayerMoves","name":"whenPlayerMoves","line":538,"end_line":540,"hash":"deb4e2c298d1bf38833abd6c76c6641e82149efbf996a1dac1d01fafcc0bf393"},{"id":"func/movePlayerToAnyExampleRoom","name":"movePlayerToAnyExampleRoom","line":542,"end_line":548,"hash":"33bb63c049b93be3a6dfdf05e4afd2b1f390d2c3fafa2310f484e2ab176f3f6e"},{"id":"func/movePlayerToRoom","name":"movePlayerToRoom","line":550,"end_line":556,"hash":"bf6b63b29408981629d90bd225638917becc17416d36808d39f05fb5133cf7d2"},{"id":"func/thenPlayerInToRoom","name":"thenPlayerInToRoom","line":558,"end_line":560,"hash":"c2200eda6575e1a42ff7aae8db34a4da933e6cfceeac6769a83245617b86d59b"},{"id":"func/thenPlayerInFromRoom","name":"thenPlayerInFromRoom","line":562,"end_line":564,"hash":"a39df878629abd764e5bf2e5251bc9674b051ecde23a4da04745eaa3ec39af55"},{"id":"func/thenPlayerInPlayerRoom","name":"thenPlayerInPlayerRoom","line":566,"end_line":568,"hash":"c75adfae4fe1f6dc48e73d3af9844df51860e33e6c3f79580648339e3ec5a894"},{"id":"func/thenPlayerInExpectedPlayerRoom","name":"thenPlayerInExpectedPlayerRoom","line":570,"end_line":572,"hash":"7a6ddc3cb5f53a8c08dbc487cc8c7e3fd1e6a4a4036ae7d7bfc5af1542b7c58e"},{"id":"func/thenPlayerInRelocationRoom","name":"thenPlayerInRelocationRoom","line":574,"end_line":576,"hash":"c9318779ea3f46571b49e01caa65aeb0898758c9f276e616af0099c391f06078"},{"id":"func/thenPlayerInExampleRoom","name":"thenPlayerInExampleRoom","line":578,"end_line":580,"hash":"141e97a9d5ebbed1a7395e7f066389d39d2b17b9464777c5ffee362211f19833"},{"id":"func/thenGameStillInProgress","name":"thenGameStillInProgress","line":582,"end_line":584,"hash":"5402e966e59ffa3cd350c5b8277acbb9455f0c10cb65c6ce7854050621b2e117"},{"id":"func/thenGameStatus","name":"thenGameStatus","line":586,"end_line":588,"hash":"4c473a54f502f6a2022ff1302210af5a031e67f97701a157a1e196f522e8d226"},{"id":"func/thenPlayerLoses","name":"thenPlayerLoses","line":590,"end_line":592,"hash":"d519700803bfc1bebacb9fe652cd6fb300012edd168e008cc115791f7fb48988"},{"id":"func/assertGameStatus","name":"assertGameStatus","line":594,"end_line":600,"hash":"4dbe14fbfe822b22afd1fb1bb5e75f2bf4de40efcc72b3566f43fd57919b94e3"},{"id":"func/thenTurnMessagesAre","name":"thenTurnMessagesAre","line":602,"end_line":604,"hash":"f635a02eac8e7318d1a243c4dfe30c26799117910525caf9e8bcad3726e35d24"},{"id":"func/thenMoveRejectedWithMessage","name":"thenMoveRejectedWithMessage","line":606,"end_line":609,"hash":"5f0cce24a56cc9f5e8cadf5e420d5dfd065b812a19ecbe1ed53fa93d861e566e"},{"id":"func/givenNextBatRelocationRoom","name":"givenNextBatRelocationRoom","line":611,"end_line":618,"hash":"fba1a91320446c26360da0000789983d081969c088971d581ca0dbf781a1acc0"},{"id":"func/givenNextWumpusWakeChoice","name":"givenNextWumpusWakeChoice","line":620,"end_line":627,"hash":"376fb8d876819d319c5edf10c3d6ad5739d3cd33b2379285eb69bbd892618634"},{"id":"func/thenWumpusInRoom","name":"thenWumpusInRoom","line":629,"end_line":631,"hash":"858514c805e2baa0579a9273a83a64940fbb9a7771c29e58985dc82444be21ea"},{"id":"func/thenWumpusInWumpusRoom","name":"thenWumpusInWumpusRoom","line":633,"end_line":635,"hash":"1f93e75dc66db361a659716e0d1ea13e29223b98b193f343c23b0dfce24f43c1"},{"id":"func/thenWumpusInExampleRoom","name":"thenWumpusInExampleRoom","line":637,"end_line":639,"hash":"b17bc210696e3c380e285c11b6569e5e02fb914b0550c229ec07d21b2cff11f8"},{"id":"func/thenSetupRoom","name":"thenSetupRoom","line":641,"end_line":651,"hash":"6f5e57cd440123a54ce5582766730f07b7fc05a92ceae24bbf2a6354206a80ab"},{"id":"func/setupPlayerRoom","name":"setupPlayerRoom","line":653,"end_line":655,"hash":"c088cccf31329f6323617271b6b1d28d6fa642bf36ee17faac1f1b7cc025b946"},{"id":"func/setupWumpusRoom","name":"setupWumpusRoom","line":657,"end_line":659,"hash":"1648c66fa2ec5b4fe2a7cf77e7d3249bfc348f2960c198ed92fbb22f531a355a"},{"id":"func/whenTurnWarningsRequested","name":"whenTurnWarningsRequested","line":661,"end_line":664,"hash":"cadf2260938820cc39f30aee91c15249bc50f73005773ae611417a1c8675eb05"},{"id":"func/thenWarningMessagesAre","name":"thenWarningMessagesAre","line":666,"end_line":668,"hash":"f21b9c015381493277db481a6fa29ee64a38145ee258aa30f4aa2cf12d4f76a3"},{"id":"func/givenOrThenPlayerHasArrows","name":"givenOrThenPlayerHasArrows","line":670,"end_line":676,"hash":"2b84982cc273afe8cbc2159b62f4c48cce712a41069dfd640fc0843e67f9a7e5"},{"id":"func/setOrAssertParsedInt","name":"setOrAssertParsedInt","line":678,"end_line":688,"hash":"94919279879d012e994a0ef1e8602908cacf78c851372c0e7a701a16c6646d21"},{"id":"func/setStringChoice","name":"setStringChoice","line":690,"end_line":696,"hash":"aece07139f8d598d5692e05ae2da2166ea433fa7b0cb2aa4fc7646ac0cf61aff"},{"id":"func/thenPlayerHasRemainingArrows","name":"thenPlayerHasRemainingArrows","line":698,"end_line":704,"hash":"1c25e30f315689fc63e09c184eb03e99d3f07bcf923845126e044fa310bbc645"},{"id":"func/assertArrows","name":"assertArrows","line":706,"end_line":712,"hash":"41617524809e988dfe58f62840b5bc15e0578ab27f438aa4c8878e0a45bed2ad"},{"id":"func/assertRejectedMessage","name":"assertRejectedMessage","line":714,"end_line":719,"hash":"661bf672f5cce7a51285db5a5df74dc9503e1d5e3dd59fb5fb836389509c3958"},{"id":"func/assertRoomList","name":"assertRoomList","line":721,"end_line":723,"hash":"acc602fedbc05083bbfecd8d6cf116c954e2a7db74aeec697d114ae831f9a951"},{"id":"func/assertOptionalRoomList","name":"assertOptionalRoomList","line":725,"end_line":727,"hash":"a511b59b0bcb46f65e6b296ea35bdce7dddf10cb7cfbe6790760124296bd43e2"},{"id":"func/assertParsedRoomListValue","name":"assertParsedRoomListValue","line":729,"end_line":735,"hash":"31db10b3e2b3649d99846ef299e3f3e04672da66a8609b758e79c55414682850"},{"id":"func/assertParsedRoomList","name":"assertParsedRoomList","line":737,"end_line":742,"hash":"2cf47aa2bd6568f34264e7e22de68807fcd981326867d8f5b0b3bb75ceaf6a93"},{"id":"func/givenInteractiveSetup","name":"givenInteractiveSetup","line":744,"end_line":756,"hash":"32e01e2ad3d2e2f22d2ce292a65d964d0581def2bacb68e2a16c859580e62277"},{"id":"func/givenInteractiveSetupSeed","name":"givenInteractiveSetupSeed","line":758,"end_line":765,"hash":"ea25c25b357cec0fe45fd03facd0e5b2a09ed58d364f66d1b663ee89fa9bf842"},{"id":"func/givenNewInteractiveSession","name":"givenNewInteractiveSession","line":767,"end_line":770,"hash":"477c6c36946723ba0cf7266b3a1f5d389b23c50149a8f51494cadc5faacd253b"},{"id":"func/whenNextTurnDisplayed","name":"whenNextTurnDisplayed","line":772,"end_line":775,"hash":"a6e40b81690976d1c2975b99a8c648996799650a71d4799a9aab1a9487d3d630"},{"id":"func/whenPlayerEntersCommand","name":"whenPlayerEntersCommand","line":777,"end_line":784,"hash":"5ffdcbbeba15175da9f5502259e931cb6f9ac29c268650005a45e65a633414d1"},{"id":"func/thenDisplayedLinesAre","name":"thenDisplayedLinesAre","line":786,"end_line":788,"hash":"25b4303db623d2eacacdc8841aee4354781e6b3ab6249092aa424c1818870930"},{"id":"func/assertStringState","name":"assertStringState","line":790,"end_line":792,"hash":"240180cc31ad79b78d459adb966dd95661a043a78f235013969a4be8c6e20782"},{"id":"func/assertStringList","name":"assertStringList","line":794,"end_line":800,"hash":"a1cf0e2e84b2cbbf92d8d60428d271f9bf87bd6c3fa5718a6d9e2bc9d76604de"},{"id":"func/thenDisplayedLinesInclude","name":"thenDisplayedLinesInclude","line":802,"end_line":815,"hash":"e01ee0d23eb08b265db40f692e282f78cafbc79f42875b9638875d8e512f284c"},{"id":"func/givenPlayerHasLost","name":"givenPlayerHasLost","line":817,"end_line":828,"hash":"9ab8662f712aa674956a999443943da00c8d875891ec21897dbe3f8780a92620"},{"id":"func/whenPlayerAnswersSameSetupPrompt","name":"whenPlayerAnswersSameSetupPrompt","line":830,"end_line":833,"hash":"b40056ab91813524c1a007e541cc3a4a1fbc437a73d8339a5d0c0717951a503a"},{"id":"func/thenNextGameSetupRelation","name":"thenNextGameSetupRelation","line":835,"end_line":846,"hash":"f373ed02fe4161887f0c5b0fd8d23f9336afcd5f0041ff4b58a1ceffdf6d73b2"},{"id":"func/whenPlayerAnswersInstructionsPrompt","name":"whenPlayerAnswersInstructionsPrompt","line":848,"end_line":851,"hash":"688365f36d502ade4504f527a8aea38cac3a37f2e3ce52a07ed88663d9528655"},{"id":"func/givenNextSleepyWumpusObservation","name":"givenNextSleepyWumpusObservation","line":853,"end_line":863,"hash":"e8bbfc658c59bb801862c9c6996fdffd98ab5c3bf3fedfe5d32b679aeea794d9"},{"id":"func/givenWumpusAsleep","name":"givenWumpusAsleep","line":865,"end_line":868,"hash":"775614c16c935e2f1becc245c106fe93f94f5599755ed886e3be81c0be438e83"},{"id":"func/givenWumpusAwake","name":"givenWumpusAwake","line":870,"end_line":873,"hash":"fecad09c04ee15b17e8d6c5b1ed4e94443be1a4c565a8b52a76e58b42cc47ba4"},{"id":"func/thenWumpusSleepState","name":"thenWumpusSleepState","line":875,"end_line":884,"hash":"14061f226838ca345f6f9d9243573625034516897b021e91c65e0608a66803d5"},{"id":"func/thenWumpusSleepStateAsleep","name":"thenWumpusSleepStateAsleep","line":886,"end_line":888,"hash":"fefafe1cb1767eb3e546c843fc8bb62a7c4b32f4ce4c19843fa35c35d9973027"},{"id":"func/thenWumpusSleepStateAwake","name":"thenWumpusSleepStateAwake","line":890,"end_line":892,"hash":"99181bed1443f4e804f57860c5e77e44649debe26ead5b1bf481823bf6249064"},{"id":"func/requireWumpusSleepState","name":"requireWumpusSleepState","line":894,"end_line":899,"hash":"8ce74cdf65372de647d1209a5f11f40e1c77c53844d0a7976ae1703f12ab2728"},{"id":"func/thenTurnWarningsAre","name":"thenTurnWarningsAre","line":901,"end_line":908,"hash":"a712e9f57634f1acbcda42327dc4adba3800ba7406dad54f9e12dd68b7c833dc"},{"id":"func/givenNextSleepingWumpusEntryOutcome","name":"givenNextSleepingWumpusEntryOutcome","line":910,"end_line":917,"hash":"d002770f3368b47adfcbed35f5f35ef005465658c60d0fdc29c683f4c48e952a"},{"id":"func/givenPlayerHasSeenSleepingWumpus","name":"givenPlayerHasSeenSleepingWumpus","line":919,"end_line":922,"hash":"23b5522d50edfc470e8abcb111fd28e7a83e5ab191ce909fed3d419ededb2e6e"},{"id":"func/whenBothGamesObserveSleepyWumpus","name":"whenBothGamesObserveSleepyWumpus","line":924,"end_line":928,"hash":"95a946dfe99a4bf8fef832a0c8991f99502328fb3b205c21240a2300c2ef4382"},{"id":"func/thenBothSleepyObservationsIdentical","name":"thenBothSleepyObservationsIdentical","line":930,"end_line":932,"hash":"21b286afdf10af00f73d37a09876474f2db52dcc0d2d5abb565b8efe5a7bf5f6"},{"id":"func/recordPairedStringObservations","name":"recordPairedStringObservations","line":934,"end_line":942,"hash":"4e85a6a31ab1c29f69acbe505435ab8239be6982c47543a07381ccdaaa38f2a3"},{"id":"func/assertStringObservationsIdentical","name":"assertStringObservationsIdentical","line":944,"end_line":951,"hash":"e69a1ec150047412f96f4293974b72a8c99318fe8e8c781e754df1a57fad4a46"},{"id":"func/setupSnapshot","name":"setupSnapshot","line":953,"end_line":959,"hash":"281567869ec92ee702b7087da698bcee785d98859c843b54ef13a421f1250f89"},{"id":"func/inspectSetup","name":"inspectSetup","line":961,"end_line":963,"hash":"a4cedce9df2e91f5b1f9455c6c992ba2feec3a9c1cec9d1493458c9732224500"},{"id":"func/storeSession","name":"storeSession","line":965,"end_line":968,"hash":"eb2ce186b12c5f32d9fd26d0c29fb554aee73d98d8e3fdfac583a9fe5c7355a1"},{"id":"func/gameFrom","name":"gameFrom","line":970,"end_line":972,"hash":"aaeee81e6f19c53a6b9197619d41ed998cf9f207a5ebdb67f6232bf0536d0312"},{"id":"func/sessionFrom","name":"sessionFrom","line":974,"end_line":976,"hash":"4b0f7f5d7d798d7125b08cf5243664b4cd8761b5ae7acec8c3a1efe79286791c"},{"id":"func/firstPresent","name":"firstPresent","line":978,"end_line":985,"hash":"50894efd3c5e00889ba877aa508c09685d2436d1b67501237727b5b4be2b4e2d"},{"id":"func/intAnyExample","name":"intAnyExample","line":987,"end_line":994,"hash":"e84e50dc42a925c038232f6c875275921708de97c348bcd6f57e63d7434c92bd"},{"id":"func/intExample","name":"intExample","line":996,"end_line":1006,"hash":"7bd03b5c2a2d06befbebf041ca10e7f319b032af0503ce8f3988f9fe147d9d33"},{"id":"func/parseWakeChoice","name":"parseWakeChoice","line":1008,"end_line":1021,"hash":"c2c170daf77137167aa459747296236d9c088f2c326af553c0508d598a96bada"},{"id":"func/int64Example","name":"int64Example","line":1023,"end_line":1033,"hash":"a4d6fde31b1bd879696ab39055fe116a932bac87d06b2b655c3aea3c1bcaf926"},{"id":"func/roomList","name":"roomList","line":1035,"end_line":1045,"hash":"12f84457daa82e5cf93a778fa64236ad21ac0a65aa731438bc4997abc57541bf"},{"id":"func/optionalRoomList","name":"optionalRoomList","line":1047,"end_line":1052,"hash":"52c2dcc631c7ecfd3c653fea0adeb4042e43a64bb640c697748a8eafdd61d9ad"},{"id":"func/stringList","name":"stringList","line":1054,"end_line":1065,"hash":"4b2429b2a6b552fc95cbfa4568e072889a4b70699ec224c66cdd503f422e66fa"},{"id":"func/hazardList","name":"hazardList","line":1067,"end_line":1069,"hash":"d4abe62e21fc6925976a4e5c439ece7273c34a6bf059caec5dac61281d875e7b"},{"id":"func/commaSeparatedStrings","name":"commaSeparatedStrings","line":1071,"end_line":1080,"hash":"9489e01c576628d5a219d1ff7829f3e5e8331959f9db1fa5b9d37e83b4fc72e7"},{"id":"func/firstUnoccupiedRoom","name":"firstUnoccupiedRoom","line":1082,"end_line":1091,"hash":"24be287e7954551866a994516d19d6f36601dff10ab00bb9c998a3cf29acf476"}]}
// mutate4go-manifest-end
