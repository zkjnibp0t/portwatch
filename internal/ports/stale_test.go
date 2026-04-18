package ports

import (
	"testing"
	"time"
)

func fixedStaleClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestStaleNotStaleBeforeThreshold(t *testing.T) {
	now := time.Now()
	st := NewStaleTracker(10 * time.Minute)
	st.clock = fixedStaleClock(now)
	st.Observe(8080)
	st.clock = fixedStaleClock(now.Add(5 * time.Minute))
	if st.IsStale(8080) {
		t.Fatal("expected not stale before threshold")
	}
}

func TestStaleDetectedAfterThreshold(t *testing.T) {
	now := time.Now()
	st := NewStaleTracker(10 * time.Minute)
	st.clock = fixedStaleClock(now)
	st.Observe(8080)
	st.clock = fixedStaleClock(now.Add(10 * time.Minute))
	if !st.IsStale(8080) {
		t.Fatal("expected stale at threshold")
	}
}

func TestStaleRemoveClearsPort(t *testing.T) {
	now := time.Now()
	st := NewStaleTracker(1 * time.Minute)
	st.clock = fixedStaleClock(now)
	st.Observe(443)
	st.Remove(443)
	st.clock = fixedStaleClock(now.Add(5 * time.Minute))
	if st.IsStale(443) {
		t.Fatal("expected not stale after removal")
	}
}

func TestStalePortsReturnsAll(t *testing.T) {
	now := time.Now()
	st := NewStaleTracker(5 * time.Minute)
	st.clock = fixedStaleClock(now)
	st.Observe(80)
	st.Observe(443)
	st.clock = fixedStaleClock(now.Add(6 * time.Minute))
	ports := st.StalePorts()
	if len(ports) != 2 {
		t.Fatalf("expected 2 stale ports, got %d", len(ports))
	}
}

func TestStaleAgeUnknownPortIsZero(t *testing.T) {
	st := NewStaleTracker(time.Minute)
	if st.Age(9999) != 0 {
		t.Fatal("expected zero age for untracked port")
	}
}

func TestStaleObserveIdempotent(t *testing.T) {
	now := time.Now()
	st := NewStaleTracker(time.Minute)
	st.clock = fixedStaleClock(now)
	st.Observe(22)
	st.clock = fixedStaleClock(now.Add(30 * time.Second))
	st.Observe(22) // should not reset first seen
	st.clock = fixedStaleClock(now.Add(90 * time.Second))
	if !st.IsStale(22) {
		t.Fatal("expected stale; second Observe should not reset timer")
	}
}
