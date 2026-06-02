package steps

import (
	"fmt"

	"htwgo/acceptance/runtime"
)

func givenOrThenTurnCount(world *runtime.World, example map[string]string) error {
	turnCount, err := intExample(example, "turn_count")
	if err != nil {
		return err
	}
	if _, actionTaken := world.State["action_taken"]; !actionTaken {
		sessionFrom(world).SetTurnCount(turnCount)
		return nil
	}
	return requireTurnCount(world, turnCount)
}

func thenExpectedTurnCount(world *runtime.World, example map[string]string) error {
	turnCount, err := intExample(example, "expected_turn_count")
	if err != nil {
		return err
	}
	return requireTurnCount(world, turnCount)
}

func requireTurnCount(world *runtime.World, want int) error {
	got := sessionFrom(world).TurnCount()
	if got != want {
		return fmt.Errorf("turn count = %d, want %d", got, want)
	}
	return nil
}
