package interactive

import (
	"reflect"
	"testing"

	"htwgo/internal/wumpus"
)

func TestDisplayTurnShowsWarningsRoomTunnelsArrowsAndPrompt(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{5, 17}}, 5))

	got := session.DisplayTurn()
	want := []string{
		"I SMELL A WUMPUS",
		"BATS NEARBY",
		"YOU ARE IN ROOM 1",
		"TUNNELS LEAD TO 2 5 8",
		"ARROWS LEFT: 5",
		"SHOOT OR MOVE (S-M)?",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("displayed lines = %v, want %v", got, want)
	}
}

func TestMoveCommandIsCaseInsensitive(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("M 2")

	if len(lines) != 0 {
		t.Fatalf("lines = %v, want none", lines)
	}
	if session.Game().Setup().Player != 2 {
		t.Fatalf("player room = %d, want 2", session.Game().Setup().Player)
	}
}

func TestShootCommandDisplaysVictory(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 2, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("s 2")

	want := []string{"AHA! YOU GOT THE WUMPUS! HEE HEE HEE - THE WUMPUS'LL GETCHA NEXT TIME!!"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %v, want %v", lines, want)
	}
	if session.Game().Status() != wumpus.StatusWon {
		t.Fatalf("status = %s, want %s", session.Game().Status(), wumpus.StatusWon)
	}
}

func TestFiveRoomShootCommandIsAccepted(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))
	session.Game().SetNextWumpusWakeChoice(wumpus.WumpusWakeChoice{Stay: true})

	lines := session.EnterCommand("s 2 3 4 5 6")

	if !reflect.DeepEqual(lines, []string{"MISSED"}) {
		t.Fatalf("lines = %v, want [MISSED]", lines)
	}
	if session.Game().Arrows() != 4 {
		t.Fatalf("arrows = %d, want 4", session.Game().Arrows())
	}
}

func TestInvalidCommandDoesNotAdvanceState(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("jump 2")

	if !reflect.DeepEqual(lines, []string{"JUMP IS NOT A COMMAND"}) {
		t.Fatalf("lines = %v", lines)
	}
	if session.Game().Setup().Player != 1 {
		t.Fatalf("player room = %d, want 1", session.Game().Setup().Player)
	}
	if session.Game().Arrows() != 5 {
		t.Fatalf("arrows = %d, want 5", session.Game().Arrows())
	}
}

func TestEmptyCommandReportsInvalidCommand(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("")

	if !reflect.DeepEqual(lines, []string{" IS NOT A COMMAND"}) {
		t.Fatalf("lines = %v", lines)
	}
}

func TestMalformedMoveCommandsReportMoveError(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	for _, command := range []string{"m", "m two", "m 20"} {
		lines := session.EnterCommand(command)
		if !reflect.DeepEqual(lines, []string{"CAN'T MOVE THERE"}) {
			t.Fatalf("%q lines = %v", command, lines)
		}
	}
}

func TestMalformedShootCommandsReportShootError(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	for _, command := range []string{"s", "s two", "s 21", "s 2 3 4 5 6 7"} {
		lines := session.EnterCommand(command)
		if !reflect.DeepEqual(lines, []string{"CAN'T SHOOT THERE"}) {
			t.Fatalf("%q lines = %v", command, lines)
		}
	}
}

func TestLosingMovePromptsForSameSetup(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{2, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("m 2")

	want := []string{"YYYIIIIEEEE . . . FELL IN PIT", "HA HA HA - YOU LOSE!", "SAME SET UP (Y-N)?"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %v, want %v", lines, want)
	}
}

func TestDifferentSetupAnswerWithoutSeedLeavesGameAlone(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{2, 14}, Bats: []int{16, 17}}, 5))
	original := session.Game().Setup()

	session.AnswerSameSetup("n")

	if !reflect.DeepEqual(session.Game().Setup(), original) {
		t.Fatalf("setup = %v, want unchanged %v", session.Game().Setup(), original)
	}
}

func TestSameSetupReplayCanPreserveOrReplaceSetup(t *testing.T) {
	session := NewSessionWithSeed(1973)
	original := session.Game().Setup()
	session.MarkLostForTest()

	session.AnswerSameSetup("Y")
	if !reflect.DeepEqual(session.Game().Setup(), original) {
		t.Fatalf("same setup replay = %v, want %v", session.Game().Setup(), original)
	}

	session.MarkLostForTest()
	session.AnswerSameSetup("n")
	expected, err := wumpus.NewGame(1974)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(session.Game().Setup(), expected.Setup()) {
		t.Fatalf("new setup replay = %v, want seed 1974 setup %v", session.Game().Setup(), expected.Setup())
	}
}

func TestParseRoomBoundaries(t *testing.T) {
	tests := []struct {
		value string
		room  int
		ok    bool
	}{
		{value: "0", room: 0, ok: false},
		{value: "1", room: 1, ok: true},
		{value: "20", room: 20, ok: true},
		{value: "21", room: 0, ok: false},
		{value: "two", room: 0, ok: false},
	}

	for _, test := range tests {
		room, ok := parseRoom(test.value)
		if room != test.room || ok != test.ok {
			t.Fatalf("parseRoom(%q) = (%d, %t), want (%d, %t)", test.value, room, ok, test.room, test.ok)
		}
	}
}

func TestInstructionPrompt(t *testing.T) {
	session := NewSession()

	if got := session.AnswerInstructions("y"); !reflect.DeepEqual(got, []string{"WELCOME TO 'HUNT THE WUMPUS'"}) {
		t.Fatalf("yes instructions = %v", got)
	}
	if got := session.AnswerInstructions("n"); len(got) != 0 {
		t.Fatalf("no instructions = %v, want none", got)
	}
}

func mustInteractiveGame(t *testing.T, setup wumpus.Setup, arrows int) wumpus.Game {
	t.Helper()
	game, err := wumpus.NewGameWithSetup(setup)
	if err != nil {
		t.Fatal(err)
	}
	game.SetArrows(arrows)
	return game
}
