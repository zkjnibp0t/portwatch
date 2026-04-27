package ports

import (
	"testing"
	"time"
)

var fixedCascadeClock = func() func() int64 {
	now := time.Now().UnixNano()
	return func() int64 { return now }
}()

func TestCascadeNotDetectedBelowMinGroup(t *testing.T) {
	cd := NewCascadeDetector(int64(time.Second), 3, fixedCascadeClock)
	events := cd.Record([]int{8080})
	if len(events) != 0 {
		t.Fatalf("expected no cascade, got %d", len(events))
	}
}

func TestCascadeDetectedAtMinGroup(t *testing.T) {
	cd := NewCascadeDetector(int64(time.Second), 2, fixedCascadeClock)
	events := cd.Record([]int{80, 443})
	if len(events) != 1 {
		t.Fatalf("expected 1 cascade, got %d", len(events))
	}
	if events[0].Trigger != 80 {
		t.Errorf("expected trigger 80, got %d", events[0].Trigger)
	}
	if len(events[0].FollowOn) != 1 || events[0].FollowOn[0] != 443 {
		t.Errorf("unexpected follow-on: %v", events[0].FollowOn)
	}
}

func TestCascadeAccumulatesAcrossCalls(t *testing.T) {
	var ts int64 = 1_000_000_000
	clock := func() int64 { return ts }
	window := int64(5 * time.Second)
	cd := NewCascadeDetector(window, 3, clock)

	// First call — only 1 port, below minGroup.
	events := cd.Record([]int{8080})
	if len(events) != 0 {
		t.Fatal("expected no cascade yet")
	}

	// Second call within window — total 3, should trigger.
	ts += int64(2 * time.Second)
	events = cd.Record([]int{8081, 8082})
	if len(events) != 1 {
		t.Fatalf("expected cascade, got %d", len(events))
	}
	if events[0].Trigger != 8080 {
		t.Errorf("expected trigger 8080, got %d", events[0].Trigger)
	}
}

func TestCascadeEvictsOldEvents(t *testing.T) {
	var ts int64 = 1_000_000_000
	clock := func() int64 { return ts }
	window := int64(2 * time.Second)
	cd := NewCascadeDetector(window, 2, clock)

	// Record first port.
	cd.Record([]int{9000})

	// Advance time beyond window.
	ts += int64(3 * time.Second)

	// New port alone should not trigger cascade.
	events := cd.Record([]int{9001})
	if len(events) != 0 {
		t.Fatalf("expected no cascade after eviction, got %d", len(events))
	}
}

func TestCascadeResetClearsState(t *testing.T) {
	cd := NewCascadeDetector(int64(time.Second), 2, fixedCascadeClock)
	cd.Record([]int{1234})
	cd.Reset()
	events := cd.Record([]int{5678})
	if len(events) != 0 {
		t.Fatal("expected no cascade after reset")
	}
}

func TestCascadeMinGroupClampedToTwo(t *testing.T) {
	// minGroup < 2 should be clamped to 2.
	cd := NewCascadeDetector(int64(time.Second), 0, fixedCascadeClock)
	events := cd.Record([]int{80, 443})
	if len(events) != 1 {
		t.Fatalf("expected cascade with clamped minGroup=2, got %d", len(events))
	}
}
