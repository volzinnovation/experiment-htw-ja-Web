package interactive

import (
	"reflect"
	"testing"

	"htwgo/internal/wumpus"
)

func TestNextTurnCanBeginWithJumpingWumpusEvent(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{11, 12})
	session := NewSessionWithGame(game)

	lines := session.BeginTurn()

	if !containsLine(lines, "YOU HEAR WHUMP, WHUMP.") {
		t.Fatalf("lines = %v", lines)
	}
	if session.Game().Setup().Wumpus != 12 {
		t.Fatalf("Wumpus room = %d, want 12", session.Game().Setup().Wumpus)
	}
}

func TestMoveCommandRunsJumpingWumpusBeforeMove(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{11, 12})
	session := NewSessionWithGame(game)

	lines := session.EnterCommand("m 5")

	if !reflect.DeepEqual(lines, []string{"YOU HEAR WHUMP, WHUMP."}) {
		t.Fatalf("lines = %v", lines)
	}
	if session.Game().Setup().Player != 5 {
		t.Fatalf("player room = %d, want 5", session.Game().Setup().Player)
	}
	if session.Game().Setup().Wumpus != 12 {
		t.Fatalf("Wumpus room = %d, want 12", session.Game().Setup().Wumpus)
	}
}

func TestPendingGrenadeDetonatesAfterJumpAndCommand(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 10, Pits: []int{14, 15}, Bats: []int{16, 17}}, 5)
	game.SetPendingGrenade(14)
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{11, 12})
	session := NewSessionWithGame(game)

	lines := session.EnterCommand("m 5")

	want := []string{"YOU HEAR WHUMP, WHUMP.", "YOU HEAR A HORRENDOUS EXPLOSION!"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %v, want %v", lines, want)
	}
	if _, ok := session.Game().PendingGrenadeRoom(); ok {
		t.Fatal("grenade should no longer be pending")
	}
}
