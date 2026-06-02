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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:24:53-05:00","module_hash":"e8c47142fb5676e0cb5c25b2d5c0e74ad8781d83cef01ff08fc18aa4b2e7039f","functions":[{"id":"func/dequeueOr","name":"dequeueOr","line":3,"end_line":10,"hash":"caec5f5a06e813508f15c4f91388406aef5f67f0008721cf799d2e8233cf49c9"},{"id":"func/periodicBehavior","name":"periodicBehavior","line":12,"end_line":22,"hash":"84ad3f8a85bf6789bf929aa7e1b80622467ed0ed74dfe7e2611f65b8acbb4f9c"}]}
// mutate4go-manifest-end
