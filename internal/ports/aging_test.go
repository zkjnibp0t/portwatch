package ports

import (
	"testing"
	"time"
)

func fixedAgingClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAgingTrackerNewPortTracked(t *testing.T) {
	now := time.Now()
	at := NewAgingTracker()
	at.clock = fixedAgingClock(now)
	at.Observe([]int{8080})

	age, ok := at.Age(8080)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age != 0 {
		t.Errorf("expected age 0, got %v", age)
	}
}

func TestAgingTrackerAgeGrowsOverTime(t *testing.T) {
	base := time.Now()
	at := NewAgingTracker()
	at.clock = fixedAgingClock(base)
	at.Observe([]int{443})

	at.clock = fixedAgingClock(base.Add(5 * time.Minute))
	at.Observe([]int{443})

	age, ok := at.Age(443)
	if !ok {
		t.Fatal("expected port tracked")
	}
	if age != 5*time.Minute {
		t.Errorf("expected 5m, got %v", age)
	}
}

func TestAgingTrackerPortRemovedWhenClosed(t *testing.T) {
	at := NewAgingTracker()
	at.Observe([]int{9000})
	at.Observe([]int{}) // port closed

	_, ok := at.Age(9000)
	if ok {
		t.Error("expected port to be removed after close")
	}
}

func TestAgingTrackerOlderThan(t *testing.T) {
	base := time.Now()
	at := NewAgingTracker()
	at.clock = fixedAgingClock(base)
	at.Observe([]int{80, 443})

	at.clock = fixedAgingClock(base.Add(10 * time.Minute))
	at.Observe([]int{80, 443})

	old := at.OlderThan(5 * time.Minute)
	if len(old) != 2 {
		t.Errorf("expected 2 old ports, got %d", len(old))
	}
}

func TestAgingTrackerOlderThanFiltersRecent(t *testing.T) {
	base := time.Now()
	at := NewAgingTracker()
	at.clock = fixedAgingClock(base)
	at.Observe([]int{80})

	at.clock = fixedAgingClock(base.Add(2 * time.Minute))
	at.Observe([]int{80})

	old := at.OlderThan(5 * time.Minute)
	if len(old) != 0 {
		t.Errorf("expected 0 old ports, got %d", len(old))
	}
}
