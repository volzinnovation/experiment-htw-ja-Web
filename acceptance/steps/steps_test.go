package steps

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"htwgo/acceptance/runtime"
)

func TestHandlersSatisfyCurrentAcceptanceFeatures(t *testing.T) {
	for _, feature := range []runtime.Feature{
		caveTopologyFeature(),
		entityPlacementFeature(),
		movementAndHazardsFeature(),
		turnWarningsFeature(),
		crookedArrowFeature(),
		interactiveLoopFeature(),
		holyHandGrenadeFeature(),
		sleepyWumpusFeature(),
		jumpingWumpusFeature(),
		restCommandFeature(),
	} {
		t.Run(feature.Name, func(t *testing.T) {
			path := writeFeature(t, feature)
			runtime.RunFeatureFile(t, path, NewHandlers())
		})
	}
}

func TestExampleParsingReportsMissingAndInvalidIntegers(t *testing.T) {
	if _, err := intExample(map[string]string{}, "room"); err == nil {
		t.Fatal("expected missing integer example error")
	}
	if _, err := intExample(map[string]string{"room": "x"}, "room"); err == nil {
		t.Fatal("expected invalid integer example error")
	}
	if _, err := int64Example(map[string]string{"seed": "x"}, "seed"); err == nil {
		t.Fatal("expected invalid int64 example error")
	}
	if _, err := intAnyExample(map[string]string{}, "room", "from_room"); err == nil {
		t.Fatal("expected missing any-key integer example error")
	}
}

func TestListsAndWakeChoicesParseFeatureValues(t *testing.T) {
	rooms, err := roomList("1, 2, 20")
	if err != nil {
		t.Fatal(err)
	}
	if want := []int{1, 2, 20}; !reflect.DeepEqual(rooms, want) {
		t.Fatalf("rooms = %v, want %v", rooms, want)
	}
	if _, err := roomList("1, x"); err == nil {
		t.Fatal("expected invalid room list error")
	}

	if hazards := hazardList("none"); hazards != nil {
		t.Fatalf("hazards = %v, want nil", hazards)
	}
	if hazards := hazardList("Wumpus, Bats"); !reflect.DeepEqual(hazards, []string{"Wumpus", "Bats"}) {
		t.Fatalf("hazards = %v", hazards)
	}

	if values := stringList("none"); values != nil {
		t.Fatalf("strings = %v, want nil", values)
	}
	if values := stringList("one, two"); !reflect.DeepEqual(values, []string{"one", "two"}) {
		t.Fatalf("strings = %v", values)
	}

	stay, err := parseWakeChoice("stay")
	if err != nil || !stay.Stay {
		t.Fatalf("stay choice = %#v, %v", stay, err)
	}
	move, err := parseWakeChoice("move to 11")
	if err != nil || move.Destination != 11 {
		t.Fatalf("move choice = %#v, %v", move, err)
	}
	if _, err := parseWakeChoice("move to x"); err == nil {
		t.Fatal("expected invalid move wake choice error")
	}
	if _, err := parseWakeChoice("wander"); err == nil {
		t.Fatal("expected unsupported wake choice error")
	}

	rooms, err = optionalRoomList("none")
	if err != nil || rooms != nil {
		t.Fatalf("optional rooms = %v, %v; want nil, nil", rooms, err)
	}
}

func caveTopologyFeature() runtime.Feature {
	return runtime.Feature{
		Name:       "Cave topology",
		Background: []runtime.Step{{Text: "a new cave"}},
		Scenarios: []runtime.Scenario{
			{
				Name:  "canonical room exits",
				Steps: []runtime.Step{{Text: "the exits for room <room> are queried"}, {Text: "the exits are <exits>"}, {Text: "the exit count is <exit_count>"}, {Text: "room <room> is not one of the exits"}},
				Examples: []map[string]string{
					{"room": "1", "exits": "2, 5, 8", "exit_count": "3"},
					{"room": "10", "exits": "2, 9, 11", "exit_count": "3"},
				},
			},
			{
				Name:  "cave is connected",
				Steps: []runtime.Step{{Text: "the cave is traversed from room 1"}, {Text: "every room from 1 through 20 is reachable"}},
			},
			{
				Name:  "cave invariants hold",
				Steps: []runtime.Step{{Text: "the cave invariants are inspected"}, {Text: "every room has exactly three exits"}, {Text: "no room is one of its own exits"}},
			},
			{
				Name:     "tunnel links are bidirectional",
				Steps:    []runtime.Step{{Text: "the tunnel from room <from_room> to room <to_room> is queried"}, {Text: "the reverse tunnel also exists"}},
				Examples: []map[string]string{{"from_room": "1", "to_room": "2"}, {"from_room": "16", "to_room": "20"}},
			},
			{
				Name: "adjacent hazard query reports neighboring hazards",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "adjacent hazards are queried from the player room"},
					{Text: "the adjacent hazard types are <hazards>"},
				},
				Examples: []map[string]string{
					{"player_room": "1", "wumpus_room": "2", "pit_rooms": "3, 4", "bat_rooms": "5, 6", "hazards": "Wumpus, Bats"},
					{"player_room": "6", "wumpus_room": "20", "pit_rooms": "1, 2", "bat_rooms": "3, 4", "hazards": "none"},
				},
			},
		},
	}
}

func entityPlacementFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Entity placement",
		Scenarios: []runtime.Scenario{
			{
				Name:  "random setup creates required occupants",
				Steps: []runtime.Step{{Text: "a new game created with seed 1973"}, {Text: "the setup is inspected"}, {Text: "there is 1 player"}, {Text: "there is 1 Wumpus"}, {Text: "there are 2 pits"}, {Text: "there are 2 bats"}},
			},
			{
				Name:     "occupied rooms are valid",
				Steps:    []runtime.Step{{Text: "a new game created with seed <seed>"}, {Text: "the occupied rooms are inspected"}, {Text: "every occupied room number is from 1 through 20"}, {Text: "exactly <occupied_count> distinct rooms are occupied by the player, Wumpus, pits, and bats"}},
				Examples: []map[string]string{{"seed": "1", "occupied_count": "6"}, {"seed": "2026", "occupied_count": "6"}},
			},
			{
				Name:     "same seed creates same setup",
				Steps:    []runtime.Step{{Text: "a new game created with seed <seed>"}, {Text: "another new game created with seed <seed>"}, {Text: "both setups are inspected"}, {Text: "both setups have identical player, Wumpus, pit, and bat rooms"}},
				Examples: []map[string]string{{"seed": "1973"}},
			},
			{
				Name:     "same setup replay preserves placement",
				Steps:    []runtime.Step{{Text: "a completed game created with seed <seed>"}, {Text: "a same setup replay is started"}, {Text: "the replay setup has identical player, Wumpus, pit, and bat rooms"}},
				Examples: []map[string]string{{"seed": "1973"}},
			},
		},
	}
}

func movementAndHazardsFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Movement and hazard resolution",
		Scenarios: []runtime.Scenario{
			{
				Name: "legal move enters an empty adjacent room",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player moves to room <to_room>"},
					{Text: "the player is in room <to_room>"},
					{Text: "the game is still in progress"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "to_room": "2", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "messages": "none"}},
			},
			{
				Name: "illegal move is rejected",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player moves to room <to_room>"},
					{Text: "the move is rejected with message <message>"},
					{Text: "the player is in room <from_room>"},
					{Text: "the game is still in progress"},
				},
				Examples: []map[string]string{{"from_room": "1", "to_room": "20", "wumpus_room": "19", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "message": "CAN'T MOVE THERE"}},
			},
			{
				Name: "moving into a pit loses",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player moves to room <pit_room>"},
					{Text: "the player loses"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "pit_room": "2", "wumpus_room": "20", "pit_rooms": "2, 14", "bat_rooms": "16, 17", "messages": "YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!"}},
			},
			{
				Name: "moving into bats relocates",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the next bat relocation room is <relocation_room>"},
					{Text: "the player moves to room <bat_room>"},
					{Text: "the player is in room <relocation_room>"},
					{Text: "the game is <game_status>"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "bat_room": "2", "relocation_room": "13", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "2, 17", "game_status": "lost", "messages": "ZAP -- SUPER BAT SNATCH! ELSEWHEREVILLE FOR YOU!, YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!"}},
			},
			{
				Name: "moving into Wumpus room wakes Wumpus",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the next Wumpus wake choice is <wake_choice>"},
					{Text: "the player moves to room <wumpus_room>"},
					{Text: "the Wumpus is in room <expected_wumpus_room>"},
					{Text: "the game is <game_status>"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "wumpus_room": "2", "wake_choice": "move to 3", "expected_wumpus_room": "3", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "game_status": "in progress", "messages": "none"}},
			},
		},
	}
}

func turnWarningsFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Turn warnings",
		Scenarios: []runtime.Scenario{{
			Name: "warnings appear for adjacent hazards",
			Steps: []runtime.Step{
				{Text: "a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
				{Text: "the turn warnings are requested"},
				{Text: "the warning messages are <warnings>"},
			},
			Examples: []map[string]string{
				{"player_room": "6", "wumpus_room": "5", "pit_rooms": "7, 15", "bat_rooms": "1, 2", "warnings": "I SMELL A WUMPUS, BATS NEARBY, I FEEL A DRAFT"},
				{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "warnings": "none"},
			},
		}},
	}
}

func crookedArrowFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Crooked arrow shooting",
		Scenarios: []runtime.Scenario{
			{
				Name: "arrow path that reaches Wumpus wins",
				Steps: []runtime.Step{
					{Text: "a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>"},
					{Text: "the player starts with <initial_arrows> arrows"},
					{Text: "the player shoots the path <path>"},
					{Text: "the requested shot path is <expected_path>"},
					{Text: "the player wins"},
					{Text: "the player has <remaining_arrows> arrows"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "10", "initial_arrows": "5", "path": "2, 10", "expected_path": "2, 10", "remaining_arrows": "4", "messages": "AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"}},
			},
			{
				Name: "invalid arrow segment deviates",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player has <arrows> arrows"},
					{Text: "the next arrow deviation room is <deviation_room>"},
					{Text: "the next Wumpus wake choice is <wake_choice>"},
					{Text: "the player shoots the path <path>"},
					{Text: "the arrow traversed rooms are <traversed_rooms>"},
					{Text: "the game is <game_status>"},
					{Text: "the player has <remaining_arrows> arrows"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "path": "3, 4", "deviation_room": "5", "traversed_rooms": "5, 4", "wake_choice": "stay", "game_status": "in progress", "remaining_arrows": "4"}},
			},
			{
				Name: "arrow can hit player",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player has <arrows> arrows"},
					{Text: "the next arrow deviation room is <deviation_room>"},
					{Text: "the player shoots the path <path>"},
					{Text: "the player loses"},
					{Text: "the player has <remaining_arrows> arrows"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "path": "3", "deviation_room": "1", "remaining_arrows": "4", "messages": "OUCH! ARROW GOT YOU!, HA HA HA - YOU LOSE!"}},
			},
			{
				Name: "missed arrow wakes Wumpus",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the player has <arrows> arrows"},
					{Text: "the next Wumpus wake choice is <wake_choice>"},
					{Text: "the player shoots the path <path>"},
					{Text: "the Wumpus is in room <expected_wumpus_room>"},
					{Text: "the game is <game_status>"},
					{Text: "the player has <remaining_arrows> arrows"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "2", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "path": "5", "wake_choice": "move to 1", "expected_wumpus_room": "1", "game_status": "lost", "remaining_arrows": "4", "messages": "MISSED, TSK TSK TSK - WUMPUS GOT YOU!, HA HA HA - YOU LOSE!"}},
			},
			{
				Name: "shooting path must contain rooms",
				Steps: []runtime.Step{
					{Text: "a shooting game setup with the player in room <player_room> and the Wumpus in room <wumpus_room>"},
					{Text: "the player starts with <initial_arrows> arrows"},
					{Text: "the player shoots the path <path>"},
					{Text: "the requested shot path is <expected_path>"},
					{Text: "the shot is rejected with message <message>"},
					{Text: "the player has <remaining_arrows> arrows"},
					{Text: "the game is still in progress"},
					{Text: "the player is in room <expected_player_room>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "10", "initial_arrows": "5", "path": "none", "expected_path": "none", "remaining_arrows": "5", "expected_player_room": "1", "message": "CAN'T SHOOT THERE"}},
			},
		},
	}
}

func interactiveLoopFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Interactive game loop",
		Scenarios: []runtime.Scenario{
			{
				Name: "each turn displays room, tunnels, warnings, arrows, and prompt",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the next turn is displayed"},
					{Text: "the displayed lines are <lines>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "2", "pit_rooms": "13, 14", "bat_rooms": "5, 17", "arrows": "5", "lines": "I SMELL A WUMPUS, BATS NEARBY, YOU ARE IN ROOM 1, TUNNELS LEAD TO 2 5 8, ARROWS LEFT: 5, SHOOT OR MOVE (S-M)?"}},
			},
			{
				Name: "move command is case insensitive",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player enters command <command>"},
					{Text: "the player is in room <to_room>"},
					{Text: "the game is still in progress"},
				},
				Examples: []map[string]string{{"from_room": "1", "to_room": "2", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "M 2"}},
			},
			{
				Name: "shoot command can win",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player enters command <command>"},
					{Text: "the player wins"},
					{Text: "the displayed lines include <message>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "2", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "s 2", "message": "AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"}},
			},
			{
				Name: "invalid command preserves state",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player enters command <command>"},
					{Text: "the displayed lines include <message>"},
					{Text: "the player is in room <player_room>"},
					{Text: "the player has <arrows> arrows"},
					{Text: "the game is still in progress"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "jump 2", "message": "JUMP IS NOT A COMMAND"}},
			},
			{
				Name: "losing move prompts same setup",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player enters command <command>"},
					{Text: "the player loses"},
					{Text: "the displayed lines include <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "wumpus_room": "20", "pit_rooms": "2, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "m 2", "messages": "YYYIIIIEEEE . . . FELL IN PIT, HA HA HA - YOU LOSE!, SAME SET UP (Y-N)?"}},
			},
			{
				Name: "same setup replay can preserve or replace setup",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with seed <seed>"},
					{Text: "the player has lost"},
					{Text: "the player answers same setup prompt with <answer>"},
					{Text: "the next game setup is <setup_relation> to the lost game setup"},
				},
				Examples: []map[string]string{
					{"seed": "1973", "answer": "y", "setup_relation": "identical"},
					{"seed": "1973", "answer": "n", "setup_relation": "different"},
				},
			},
			{
				Name: "instruction prompt can show or skip instructions",
				Steps: []runtime.Step{
					{Text: "a new interactive session"},
					{Text: "the player answers instructions prompt with <answer>"},
					{Text: "the displayed lines are <lines>"},
				},
				Examples: []map[string]string{
					{"answer": "y", "lines": "WELCOME TO 'HUNT THE WUMPUS'"},
					{"answer": "n", "lines": "none"},
				},
			},
		},
	}
}

func holyHandGrenadeFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Holy Hand Grenade",
		Scenarios: []runtime.Scenario{
			{
				Name: "setup places one grenade in a valid empty room",
				Steps: []runtime.Step{
					{Text: "a new game created with seed <seed>"},
					{Text: "the setup is inspected"},
					{Text: "there is 1 Holy Hand Grenade"},
					{Text: "the Holy Hand Grenade room is from 1 through 20"},
					{Text: "the Holy Hand Grenade room is not occupied by the player, Wumpus, pits, or bats"},
				},
				Examples: []map[string]string{{"seed": "1973"}},
			},
			{
				Name: "same seed creates same grenade placement",
				Steps: []runtime.Step{
					{Text: "a new game created with seed <seed>"},
					{Text: "another new game created with seed <seed>"},
					{Text: "both setups are inspected"},
					{Text: "both setups have identical Holy Hand Grenade rooms"},
				},
				Examples: []map[string]string{{"seed": "1973"}},
			},
			{
				Name: "entering grenade room acquires it",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and the Holy Hand Grenade in room <grenade_room>"},
					{Text: "the player moves to room <grenade_room>"},
					{Text: "the player carries the Holy Hand Grenade"},
					{Text: "the cave has no unclaimed Holy Hand Grenade"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"from_room": "1", "grenade_room": "2", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "messages": "YOU FOUND THE HOLY HAND GRENADE! USE IT WISELY!"}},
			},
			{
				Name: "armed prompt includes throw",
				Steps: []runtime.Step{
					{Text: "an interactive game setup where the player carries the Holy Hand Grenade and has <arrows> arrows"},
					{Text: "the next turn is displayed"},
					{Text: "the displayed lines include <prompt>"},
				},
				Examples: []map[string]string{{"arrows": "5", "prompt": "SHOOT, MOVE OR THROW (S-M-T)?"}},
			},
			{
				Name: "player without grenade cannot throw",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player enters command <command>"},
					{Text: "the displayed lines include <message>"},
					{Text: "the player does not carry the Holy Hand Grenade"},
					{Text: "the game is still in progress"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "t 2", "message": "CAN'T THROW THERE"}},
			},
			{
				Name: "invalid throw does not consume grenade",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player carries the Holy Hand Grenade"},
					{Text: "the player enters command <command>"},
					{Text: "the displayed lines include <message>"},
					{Text: "the player carries the Holy Hand Grenade"},
					{Text: "no Holy Hand Grenade detonation is pending"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "t 20", "message": "CAN'T THROW THERE"}},
			},
			{
				Name: "successful throw starts fuse",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the player carries the Holy Hand Grenade"},
					{Text: "the player enters command <command>"},
					{Text: "the player does not carry the Holy Hand Grenade"},
					{Text: "the Holy Hand Grenade is pending detonation in room <target_room>"},
					{Text: "the displayed lines include <message>"},
					{Text: "the Wumpus is alive"},
				},
				Examples: []map[string]string{{"player_room": "1", "target_room": "2", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "command": "t 2", "message": "YOU HEAR TIC...TIC..."}},
			},
			{
				Name: "pending grenade detonates after legal turn",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the Holy Hand Grenade is pending detonation in room <target_room>"},
					{Text: "the player enters command <command>"},
					{Text: "the remaining bat rooms are <remaining_bat_rooms>"},
					{Text: "the pit rooms are <pit_rooms>"},
					{Text: "the displayed lines include <message>"},
					{Text: "no Holy Hand Grenade detonation is pending"},
				},
				Examples: []map[string]string{{"player_room": "1", "target_room": "10", "wumpus_room": "20", "pit_rooms": "9, 14", "bat_rooms": "2, 16", "arrows": "5", "command": "m 5", "remaining_bat_rooms": "16", "message": "YOU HEAR A HORRENDOUS EXPLOSION!"}},
			},
			{
				Name: "same setup replay preserves grenade state",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with seed <seed>"},
					{Text: "the player carries the Holy Hand Grenade"},
					{Text: "the Holy Hand Grenade is pending detonation in room <target_room>"},
					{Text: "the player has lost"},
					{Text: "the player answers same setup prompt with <answer>"},
					{Text: "the replay setup has identical player, Wumpus, pit, bat, and Holy Hand Grenade rooms"},
					{Text: "the replay Holy Hand Grenade pending detonation room is <target_room>"},
				},
				Examples: []map[string]string{{"seed": "1973", "target_room": "10", "answer": "y"}},
			},
		},
	}
}

func sleepyWumpusFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Sleepy Wumpus",
		Scenarios: []runtime.Scenario{
			{
				Name: "adjacent room can include snoring with normal smell",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the next sleepy Wumpus adjacent observation is <sleepy_observation>"},
					{Text: "the player moves to room <to_room>"},
					{Text: "the Wumpus sleep state is <sleep_state>"},
					{Text: "the turn warnings are <warnings>"},
				},
				Examples: []map[string]string{
					{"from_room": "6", "to_room": "5", "wumpus_room": "1", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "sleepy_observation": "asleep", "sleep_state": "asleep", "warnings": "I SMELL A WUMPUS, YOU HEAR HORRIBLE SNORING"},
					{"from_room": "6", "to_room": "5", "wumpus_room": "1", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "sleepy_observation": "awake", "sleep_state": "awake", "warnings": "I SMELL A WUMPUS"},
				},
			},
			{
				Name: "entering sleeping Wumpus room can wake and kill",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <from_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the Wumpus is asleep"},
					{Text: "the next sleeping Wumpus room entry outcome is <entry_outcome>"},
					{Text: "the player moves to room <wumpus_room>"},
					{Text: "the game is <game_status>"},
					{Text: "the Wumpus sleep state is <sleep_state>"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{
					{"from_room": "1", "wumpus_room": "2", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "entry_outcome": "wakes", "game_status": "lost", "sleep_state": "awake", "messages": "YOU HEAR THE WUMPUS SAY \"YUMMY BREAKFAST!\", HA HA HA - YOU LOSE!"},
					{"from_room": "1", "wumpus_room": "2", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "entry_outcome": "stays asleep", "game_status": "in progress", "sleep_state": "asleep", "messages": "YOU SEE THE HUDDLED HORRIBLE SHAPE OF THE SLEEPING WUMPUS"},
				},
			},
			{
				Name: "leaving after seeing sleeping Wumpus awakens it",
				Steps: []runtime.Step{
					{Text: "a game setup with the player in room <wumpus_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, and bats in rooms <bat_rooms>"},
					{Text: "the Wumpus is asleep"},
					{Text: "the player has seen the sleeping Wumpus shape"},
					{Text: "the player moves to room <to_room>"},
					{Text: "the Wumpus sleep state is awake"},
					{Text: "the player is in room <to_room>"},
					{Text: "the turn messages are <messages>"},
				},
				Examples: []map[string]string{{"wumpus_room": "2", "to_room": "1", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "messages": "YOU HEAR A PETULANT SCREAM!"}},
			},
			{
				Name: "seeded sleepy observations are reproducible",
				Steps: []runtime.Step{
					{Text: "a new game created with seed <seed>"},
					{Text: "another new game created with seed <seed>"},
					{Text: "both games observe sleepy Wumpus behavior for <turn_count> turns"},
					{Text: "both games produce identical sleepy Wumpus observations"},
				},
				Examples: []map[string]string{{"seed": "1973", "turn_count": "10"}},
			},
		},
	}
}

func jumpingWumpusFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Jumping Wumpus",
		Scenarios: []runtime.Scenario{
			{
				Name: "seeded turn can trigger double jump event",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the next jumping Wumpus turn event is <jump_event>"},
					{Text: "the next Wumpus jump path is <jump_path>"},
					{Text: "the next turn begins"},
					{Text: "the Wumpus is in room <expected_wumpus_room>"},
					{Text: "the displayed lines include <message>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "10", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "jump_event": "jumps", "jump_path": "11, 12", "expected_wumpus_room": "12", "message": "YOU HEAR WHUMP, WHUMP."}},
			},
			{
				Name: "no jump event leaves Wumpus in place",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the next jumping Wumpus turn event is <jump_event>"},
					{Text: "the next turn begins"},
					{Text: "the Wumpus is in room <wumpus_room>"},
					{Text: "the displayed lines do not include <message>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "10", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "jump_event": "no jump", "message": "YOU HEAR WHUMP, WHUMP."}},
			},
			{
				Name: "first jump landing on player can trample or slam",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the next jumping Wumpus turn event is jumps"},
					{Text: "the next Wumpus jump path is <jump_path>"},
					{Text: "the next first jump player landing outcome is <landing_outcome>"},
					{Text: "the next turn begins"},
					{Text: "the game is <game_status>"},
					{Text: "the displayed lines include <messages>"},
				},
				Examples: []map[string]string{
					{"player_room": "2", "wumpus_room": "10", "jump_path": "2, 1", "landing_outcome": "trample", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "game_status": "lost", "messages": "YOU HEAR WHUMP, WHUMP., THE WUMPUS TRAMPLES YOU TO DEATH!, HA HA HA - YOU LOSE!"},
					{"player_room": "2", "wumpus_room": "10", "jump_path": "2, 1", "landing_outcome": "slam", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "game_status": "in progress", "messages": "YOU HEAR WHUMP, WHUMP., YOU ARE SLAMMED AGAINST THE CAVE WALL BY THE SNARLING WUMPUS!"},
				},
			},
			{
				Name: "jumps only follow legal tunnels and seeded events are reproducible",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the next jumping Wumpus turn event is jumps"},
					{Text: "the next Wumpus jump path is <jump_path>"},
					{Text: "the next turn begins"},
					{Text: "every Wumpus jump segment is a legal tunnel"},
					{Text: "the Wumpus is in room <expected_wumpus_room>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "jump_path": "19, 18", "expected_wumpus_room": "18", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5"}},
			},
			{
				Name: "seeded jump events are reproducible",
				Steps: []runtime.Step{
					{Text: "a new game created with seed <seed>"},
					{Text: "another new game created with seed <seed>"},
					{Text: "both games evaluate jumping Wumpus behavior for <turn_count> turns"},
					{Text: "both games produce identical jumping Wumpus events"},
				},
				Examples: []map[string]string{{"seed": "1973", "turn_count": "20"}},
			},
		},
	}
}

func restCommandFeature() runtime.Feature {
	return runtime.Feature{
		Name: "Rest command",
		Scenarios: []runtime.Scenario{
			{
				Name: "rest consumes a full turn",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the turn count is <turn_count>"},
					{Text: "the player enters command <command>"},
					{Text: "the turn count is <expected_turn_count>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "turn_count": "1", "command": "r", "expected_turn_count": "2"}},
			},
			{
				Name: "invalid rest syntax does not advance the game",
				Steps: []runtime.Step{
					{Text: "an interactive game setup with the player in room <player_room>, the Wumpus in room <wumpus_room>, pits in rooms <pit_rooms>, bats in rooms <bat_rooms>, and <arrows> arrows"},
					{Text: "the turn count is <turn_count>"},
					{Text: "the player enters command <command>"},
					{Text: "the displayed lines include <message>"},
					{Text: "the turn count is <turn_count>"},
					{Text: "the player is in room <player_room>"},
				},
				Examples: []map[string]string{{"player_room": "1", "wumpus_room": "20", "pit_rooms": "13, 14", "bat_rooms": "16, 17", "arrows": "5", "turn_count": "1", "command": "rest 12", "message": "REST IS NOT A COMMAND"}},
			},
		},
	}
}

func writeFeature(t *testing.T, feature runtime.Feature) string {
	t.Helper()
	dir, err := os.MkdirTemp(".", "steps-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})
	path := filepath.Join(dir, "feature.json")
	data, err := json.Marshal(feature)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
