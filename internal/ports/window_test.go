package ports

import (
	"testing"
	"time"
)

func TestWindowAggregatorEmptyUnion(t *testing.T) {
	w := NewWindowAggregator(5 * time.Second)
	if len(w.Union()) != 0 {
		t.Fatal("expected empty union")
	}
}

func TestWindowAggregatorUnionMergesSets(t *testing.T) {
	w := NewWindowAggregator(5 * time.Second)
	w.Add(map[int]struct{}{80: {}, 443: {}})
	w.Add(map[int]struct{}{8080: {}})
	u := w.Union()
	for _, p := range []int{80, 443, 8080} {
		if _, ok := u[p]; !ok {
			t.Errorf("expected port %d in union", p)
		}
	}
}

func TestWindowAggregatorEvictsOldBuckets(t *testing.T) {
	now := time.Unix(1000, 0)
	w := NewWindowAggregator(5 * time.Second)
	w.clock = func() time.Time { return now }

	w.Add(map[int]struct{}{22: {}})
	now = now.Add(6 * time.Second)
	w.Add(map[int]struct{}{80: {}})

	u := w.Union()
	if _, ok := u[22]; ok {
		t.Error("port 22 should have been evicted")
	}
	if _, ok := u[80]; !ok {
		t.Error("port 80 should be present")
	}
}

func TestWindowAggregatorLen(t *testing.T) {
	now := time.Unix(1000, 0)
	w := NewWindowAggregator(5 * time.Second)
	w.clock = func() time.Time { return now }

	w.Add(map[int]struct{}{1: {}})
	w.Add(map[int]struct{}{2: {}})
	if w.Len() != 2 {
		t.Fatalf("expected 2 buckets, got %d", w.Len())
	}

	now = now.Add(10 * time.Second)
	if w.Len() != 0 {
		t.Fatalf("expected 0 buckets after eviction, got %d", w.Len())
	}
}

func TestWindowAggregatorDoesNotMutateInput(t *testing.T) {
	w := NewWindowAggregator(5 * time.Second)
	input := map[int]struct{}{9090: {}}
	w.Add(input)
	delete(input, 9090)
	if _, ok := w.Union()[9090]; !ok {
		t.Error("window should hold its own copy of the input set")
	}
}
