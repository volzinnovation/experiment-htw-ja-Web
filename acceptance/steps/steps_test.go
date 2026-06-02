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
	for _, feature := range []runtime.Feature{caveTopologyFeature(), entityPlacementFeature()} {
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
}

func TestRoomAndHazardListsParseFeatureValues(t *testing.T) {
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
