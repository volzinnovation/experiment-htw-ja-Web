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
		"the cave invariants are inspected":                             whenCaveInvariantsInspected,
		"every room has exactly three exits":                            thenEveryRoomHasThreeExits,
		"no room is one of its own exits":                               thenNoRoomIsItsOwnExit,
		"the tunnel from room <from_room> to room <to_room> is queried": whenTunnelQueried,
		"the reverse tunnel also exists":                                thenReverseTunnelExists,
		"room <room> is not one of the exits":                           thenRoomIsNotExit,
		"a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>": givenConfiguredSetup,
		"a game setup with the player in room 1, the Wumpus in room 2, pits in rooms 3 and 4, and bats in rooms 5 and 6":                                 givenSetupRoom1WumpusBats,
		"a game setup with the player in room 10, the Wumpus in room 2, pits in rooms 9 and 18, and bats in rooms 6 and 7":                               givenSetupRoom10WumpusPit,
		"a game setup with the player in room 6, the Wumpus in room 20, pits in rooms 1 and 2, and bats in rooms 3 and 4":                                givenSetupRoom6NoHazards,
		"adjacent hazards are queried from the player room":                                                                                              whenAdjacentHazardsQueried,
		"the adjacent hazard types are <hazards>":                                                                                                        thenAdjacentHazardsAre,
		"the adjacent hazard types are Wumpus and Bats":                                                                                                  thenAdjacentHazardsWumpusBats,
		"the adjacent hazard types are Wumpus and Pit":                                                                                                   thenAdjacentHazardsWumpusPit,
		"there are no adjacent hazard types":                                                                                                             thenNoAdjacentHazards,
		"a new game created with seed 1973":                                                                                                              givenNewGameSeed1973,
		"a new game created with seed <seed>":                                                                                                            givenNewGameSeed,
		"the setup is inspected":                                                                                                                         whenSetupInspected,
		"there is 1 player":                                                                                                                              thenOnePlayer,
		"there is 1 Wumpus":                                                                                                                              thenOneWumpus,
		"there are 2 pits":                                                                                                                               thenTwoPits,
		"there are 2 bats":                                                                                                                               thenTwoBats,
		"the occupied rooms are inspected":                                                                                                               whenOccupiedRoomsInspected,
		"every occupied room number is from 1 through 20":                                                                                                thenOccupiedRoomsValid,
		"exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats":                                                     thenDistinctOccupiedCount,
		"exactly 6 distinct rooms are occupied by the player, Wumpus, pits, and bats":                                                                    thenSixDistinctOccupiedRooms,
		"another new game created with seed 1973":                                                                                                        givenAnotherNewGameSeed1973,
		"another new game created with seed <seed>":                                                                                                      givenAnotherNewGameSeed,
		"both setups are inspected":                                                                                                                      whenBothSetupsInspected,
		"both setups have identical player, Wumpus, pit, and bat rooms":                                                                                  thenBothSetupsIdentical,
		"a completed game created with seed 1973":                                                                                                        givenCompletedGameSeed1973,
		"a completed game created with seed <seed>":                                                                                                      givenCompletedGameSeed,
		"a same setup replay is started":                                                                                                                 whenSameSetupReplayStarted,
		"the replay setup has identical player, Wumpus, pit, and bat rooms":                                                                              thenReplaySetupIdentical,
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
	player, err := intExample(example, "player_room")
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
	if _, err := wumpus.NewGameWithSetup(setup); err != nil {
		return err
	}
	world.State["setup"] = setup
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
	world.State[key] = game
	return nil
}

func whenSetupInspected(world *runtime.World, _ map[string]string) error {
	world.State["inspected_setup"] = setupSnapshot(world.State["game"].(wumpus.Game).Setup())
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
	setup := world.State["game"].(wumpus.Game).Setup()
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
	world.State["inspected_setup"] = setupSnapshot(world.State["game"].(wumpus.Game).Setup())
	world.State["another_inspected_setup"] = setupSnapshot(world.State["another_game"].(wumpus.Game).Setup())
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
	replay := world.State["game"].(wumpus.Game).ReplaySameSetup()
	world.State["replay_setup"] = setupSnapshot(replay.Setup())
	return nil
}

func thenReplaySetupIdentical(world *runtime.World, _ map[string]string) error {
	original := setupSnapshot(world.State["game"].(wumpus.Game).Setup())
	replay := world.State["replay_setup"].(inspectedSetup)
	if !reflect.DeepEqual(original, replay) {
		return fmt.Errorf("replay setup %v differs from original %v", replay, original)
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
