package steps

import (
	"fmt"

	"htwgo/acceptance/runtime"
)

func givenOrThenTurnCount(world *runtime.World, example map[string]string) error {
	return setOrAssertParsedInt(world, example, "turn_count", func(turnCount int) {
		sessionFrom(world).SetTurnCount(turnCount)
	}, func(turnCount int) error {
		return requireTurnCount(world, turnCount)
	})
}

func givenTurnCount(world *runtime.World, example map[string]string) error {
	turnCount, err := intExample(example, "turn_count")
	if err != nil {
		return err
	}
	sessionFrom(world).SetTurnCount(turnCount)
	return nil
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
