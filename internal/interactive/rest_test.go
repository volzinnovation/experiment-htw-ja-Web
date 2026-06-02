package interactive

import (
	"reflect"
	"testing"

	"htwgo/internal/wumpus"
)

func TestRestCommandLeavesPlayerAndArrowsAlone(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))

	lines := session.EnterCommand("ReSt")

	if len(lines) != 0 {
		t.Fatalf("lines = %v, want none", lines)
	}
	if session.Game().Setup().Player != 1 {
		t.Fatalf("player room = %d, want 1", session.Game().Setup().Player)
	}
	if session.Game().Arrows() != 5 {
		t.Fatalf("arrows = %d, want 5", session.Game().Arrows())
	}
}

func TestRestCommandConsumesTurn(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))
	session.SetTurnCount(7)

	session.EnterCommand("r")

	if session.TurnCount() != 8 {
		t.Fatalf("turn count = %d, want 8", session.TurnCount())
	}
}

func TestRestRunsJumpingWumpusBeforeCompletingTurn(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 10, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5)
	game.SetNextJumpingWumpusTurnEvent(true)
	game.SetNextWumpusJumpPath([]int{11, 12})
	session := NewSessionWithGame(game)

	lines := session.EnterCommand("r")

	if !reflect.DeepEqual(lines, []string{"YOU HEAR WHUMP, WHUMP."}) {
		t.Fatalf("lines = %v", lines)
	}
	if session.Game().Setup().Wumpus != 12 {
		t.Fatalf("Wumpus room = %d, want 12", session.Game().Setup().Wumpus)
	}
}

func TestRestDetonatesPendingGrenadeAtEndOfTurn(t *testing.T) {
	game := mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{14, 15}, Bats: []int{16, 17}}, 5)
	game.SetPendingGrenade(13)
	session := NewSessionWithGame(game)

	lines := session.EnterCommand("rest")

	if !containsLine(lines, "YOU HEAR A HORRENDOUS EXPLOSION!") {
		t.Fatalf("lines = %v", lines)
	}
	if _, ok := session.Game().PendingGrenadeRoom(); ok {
		t.Fatal("grenade should no longer be pending")
	}
}

func TestInvalidRestSyntaxDoesNotAdvanceTurn(t *testing.T) {
	session := NewSessionWithGame(mustInteractiveGame(t, wumpus.Setup{Player: 1, Wumpus: 20, Pits: []int{13, 14}, Bats: []int{16, 17}}, 5))
	session.SetTurnCount(1)

	lines := session.EnterCommand("rest 12")

	if !reflect.DeepEqual(lines, []string{"REST IS NOT A COMMAND"}) {
		t.Fatalf("lines = %v", lines)
	}
	if session.TurnCount() != 1 {
		t.Fatalf("turn count = %d, want 1", session.TurnCount())
	}
}
