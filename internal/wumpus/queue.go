package wumpus

func dequeueOr[T any](values *[]T, fallback T) T {
	if len(*values) == 0 {
		return fallback
	}
	value := (*values)[0]
	*values = (*values)[1:]
	return value
}

func periodicBehavior(player, wumpus, turnCount, divisor int, match, miss string) []string {
	events := make([]string, 0, turnCount)
	for i := 0; i < turnCount; i++ {
		if (player+wumpus+i)%divisor == 0 {
			events = append(events, match)
		} else {
			events = append(events, miss)
		}
	}
	return events
}
