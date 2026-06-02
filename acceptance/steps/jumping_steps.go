package steps

import (
	"fmt"
	"slices"

	"htwgo/acceptance/runtime"
	"htwgo/internal/wumpus"
)

func givenNextJumpingWumpusTurnEvent(world *runtime.World, example map[string]string) error {
	return setNextJumpingWumpusTurnEvent(world, example["jump_event"])
}

func givenNextJumpingWumpusTurnEventJumps(world *runtime.World, _ map[string]string) error {
	return setNextJumpingWumpusTurnEvent(world, "jumps")
}

func givenNextJumpingWumpusTurnEventNoJump(world *runtime.World, _ map[string]string) error {
	return setNextJumpingWumpusTurnEvent(world, "no jump")
}

func setNextJumpingWumpusTurnEvent(world *runtime.World, value string) error {
	switch value {
	case "jumps":
		gameFrom(world, "game").SetNextJumpingWumpusTurnEvent(true)
	case "no jump":
		gameFrom(world, "game").SetNextJumpingWumpusTurnEvent(false)
	default:
		return fmt.Errorf("unsupported jumping Wumpus turn event %q", value)
	}
	return nil
}

func givenNextWumpusJumpPath(world *runtime.World, example map[string]string) error {
	path, err := roomList(example["jump_path"])
	if err != nil {
		return err
	}
	world.State["jump_start_room"] = gameFrom(world, "game").Setup().Wumpus
	world.State["jump_path"] = path
	gameFrom(world, "game").SetNextWumpusJumpPath(path)
	return nil
}

func givenNextFirstJumpLandingOutcome(world *runtime.World, example map[string]string) error {
	return setStringChoice(example["landing_outcome"], "first jump landing outcome", []string{
		string(wumpus.FirstJumpTramples),
		string(wumpus.FirstJumpSlams),
	}, func(value string) {
		gameFrom(world, "game").SetNextFirstJumpPlayerLandingOutcome(wumpus.FirstJumpLandingOutcome(value))
	})
}

func whenNextTurnBegins(world *runtime.World, _ map[string]string) error {
	lines := sessionFrom(world).BeginTurn()
	world.State["displayed_lines"] = lines
	world.State["turn_messages"] = lines
	world.State["action_taken"] = true
	return nil
}

func thenDisplayedLinesDoNotInclude(world *runtime.World, example map[string]string) error {
	lines := world.State["displayed_lines"].([]string)
	for _, unwanted := range stringList(example["message"]) {
		if slices.Contains(lines, unwanted) {
			return fmt.Errorf("displayed lines %v include %q", lines, unwanted)
		}
	}
	return nil
}

func thenPlayerMayTakeNextCommand(world *runtime.World, _ map[string]string) error {
	if sessionFrom(world).Game().Status() != wumpus.StatusInProgress {
		return fmt.Errorf("game status = %s, want in progress", sessionFrom(world).Game().Status())
	}
	return nil
}

func thenEveryWumpusJumpSegmentLegal(world *runtime.World, _ map[string]string) error {
	from := world.State["jump_start_room"].(int)
	for _, to := range world.State["jump_path"].([]int) {
		if !wumpus.NewCave().HasTunnel(from, to) {
			return fmt.Errorf("jump segment %d to %d is not legal", from, to)
		}
		from = to
	}
	return nil
}

func whenBothGamesEvaluateJumpingWumpus(world *runtime.World, example map[string]string) error {
	return recordPairedStringObservations(world, example, "jump_events", "another_jump_events", func(game *wumpus.Game, turnCount int) []string {
		return game.ObserveJumpingWumpusBehavior(turnCount)
	})
}

func thenBothJumpEventsIdentical(world *runtime.World, _ map[string]string) error {
	return assertStringObservationsIdentical(world, "jump events", "jump_events", "another_jump_events")
}
