package wumpus

import "testing"

func TestDequeueOrConsumesQueuedValues(t *testing.T) {
	values := []int{3, 5}

	if got := dequeueOr(&values, 1); got != 3 {
		t.Fatalf("first dequeue = %d, want 3", got)
	}
	if got := dequeueOr(&values, 1); got != 5 {
		t.Fatalf("second dequeue = %d, want 5", got)
	}
	if got := dequeueOr(&values, 1); got != 1 {
		t.Fatalf("empty dequeue = %d, want fallback 1", got)
	}
}
