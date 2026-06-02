package steps

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"htwgo/acceptance/runtime"
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
		"the tunnel from room <from_room> to room <to_room> is queried": whenTunnelQueried,
		"the reverse tunnel also exists":                                thenReverseTunnelExists,
		"room <room> is not one of the exits":                           thenRoomIsNotExit,
		"a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>": givenConfiguredSetup,
		"a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>":   givenConfiguredSetup,
		"adjacent hazards are queried from the player room": whenAdjacentHazardsQueried,
		"the adjacent hazard types are <hazards>":           thenAdjacentHazardsAre,
		"a new game created with seed 1973":                 givenNewGameSeed1973,
		"a new game created with seed <seed>":               givenNewGameSeed,
		"the setup is inspected":                            whenSetupInspected,
		"there is 1 player":                                 thenOnePlayer,
		"there is 1 Wumpus":                                 thenOneWumpus,
		"there are 2 pits":                                  thenTwoPits,
		"there are 2 bats":                                  thenTwoBats,
		"the occupied rooms are inspected":                  whenOccupiedRoomsInspected,
		"every occupied room number is from 1 through 20":   thenOccupiedRoomsValid,
		"exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats": thenDistinctOccupiedCount,
		"another new game created with seed <seed>":                                                  givenAnotherNewGameSeed,
		"both setups are inspected":                                                                  whenBothSetupsInspected,
		"both setups have identical player, Wumpus, pit, and bat rooms":                              thenBothSetupsIdentical,
		"a completed game created with seed <seed>":                                                  givenCompletedGameSeed,
		"a same setup replay is started":                                                             whenSameSetupReplayStarted,
		"the replay setup has identical player, Wumpus, pit, and bat rooms":                          thenReplaySetupIdentical,
		"the player moves to room <to_room>":                                                         whenPlayerMoves,
		"the player moves to room <pit_room>":                                                        whenPlayerMoves,
		"the player moves to room <bat_room>":                                                        whenPlayerMoves,
		"the player moves to room <wumpus_room>":                                                     whenPlayerMoves,
		"the player is in room <to_room>":                                                            thenPlayerInToRoom,
		"the player is in room <from_room>":                                                          thenPlayerInFromRoom,
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
	game, err := wumpus.NewGameWithSetup(setup)
	if err != nil {
		return err
	}
	world.State["setup"] = setup
	world.State["game"] = &game
	return nil
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

func whenAdjacentHazardsQueried(world *runtime.World, _ map[string]string) error {
	setup := world.State["setup"].(wumpus.Setup)
	world.State["hazards"] = wumpus.NewCave().AdjacentHazards(setup.Player, setup)
	return nil
}

func thenAdjacentHazardsAre(world *runtime.World, example map[string]string) error {
	want := hazardList(example["hazards"])
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

func givenCompletedGameSeed(world *runtime.World, example map[string]string) error {
	return setGameFromExample(world, example, "game")
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
	setup := world.State["inspected_setup"].(inspectedSetup)
	return requirePlaced(setup.Player, "player not placed")
}

func thenOneWumpus(world *runtime.World, _ map[string]string) error {
	setup := world.State["inspected_setup"].(inspectedSetup)
	return requirePlaced(setup.Wumpus, "Wumpus not placed")
}

func thenTwoPits(world *runtime.World, _ map[string]string) error {
	setup := world.State["inspected_setup"].(inspectedSetup)
	return requireCount(len(setup.Pits), 2, "pits")
}

func thenTwoBats(world *runtime.World, _ map[string]string) error {
	setup := world.State["inspected_setup"].(inspectedSetup)
	return requireCount(len(setup.Bats), 2, "bats")
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
	return commaSeparatedStrings(value)
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

func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
