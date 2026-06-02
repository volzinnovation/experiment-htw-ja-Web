package interactive

import (
	"reflect"
	"testing"

	"htwgo/internal/wumpus"
)

func TestArmedTurnPromptIncludesThrow(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.GiveGrenade()
	session := NewSessionWithGame(game)

	lines := session.DisplayTurn()

	if !containsLine(lines, "SHOOT, MOVE OR THROW (S-M-T)?") {
		t.Fatalf("lines = %v", lines)
	}
}

func TestThrowCommandStartsFuseAndNextTurnDetonates(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.GiveGrenade()
	session := NewSessionWithGame(game)

	throwLines := session.EnterCommand("T 2")
	moveLines := session.EnterCommand("m 5")

	if !reflect.DeepEqual(throwLines, []string{"YOU HEAR TIC...TIC..."}) {
		t.Fatalf("throw lines = %v", throwLines)
	}
	if !containsLine(moveLines, "YOU HEAR A HORRENDOUS EXPLOSION!") {
		t.Fatalf("move lines = %v", moveLines)
	}
}

func TestMalformedThrowCommandsReportThrowError(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.GiveGrenade()
	session := NewSessionWithGame(game)

	for _, command := range []string{"t", "t two", "t 20"} {
		lines := session.EnterCommand(command)
		if !reflect.DeepEqual(lines, []string{"CAN'T THROW THERE"}) {
			t.Fatalf("%q lines = %v", command, lines)
		}
	}
}

func containsLine(lines []string, target string) bool {
	for _, line := range lines {
		if line == target {
			return true
		}
	}
	return false
}
