package ports

import (
	"testing"
	"time"
)

func fixedVelocityClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestVelocityEmptySnapshot(t *testing.T) {
	v := NewVelocityTracker(time.Minute)
	snap := v.Snapshot()
	if snap.Opened != 0 || snap.Closed != 0 {
		t.Fatalf("expected zero snapshot, got %+v", snap)
	}
}

func TestVelocityAggregatesWithinWindow(t *testing.T) {
	now := time.Now()
	v := NewVelocityTracker(time.Minute)
	v.clock = fixedVelocityClock(now)
	v.Record(3, 1)
	v.Record(2, 4)
	snap := v.Snapshot()
	if snap.Opened != 5 {
		t.Errorf("expected Opened=5, got %d", snap.Opened)
	}
	if snap.Closed != 5 {
		t.Errorf("expected Closed=5, got %d", snap.Closed)
	}
}

func TestVelocityEvictsOldEvents(t *testing.T) {
	now := time.Now()
	v := NewVelocityTracker(time.Minute)
	v.clock = fixedVelocityClock(now.Add(-2 * time.Minute))
	v.Record(10, 5)
	v.clock = fixedVelocityClock(now)
	v.Record(1, 1)
	snap := v.Snapshot()
	if snap.Opened != 1 {
		t.Errorf("expected Opened=1 after eviction, got %d", snap.Opened)
	}
	if snap.Closed != 1 {
		t.Errorf("expected Closed=1 after eviction, got %d", snap.Closed)
	}
}

func TestVelocityResetClearsAll(t *testing.T) {
	v := NewVelocityTracker(time.Minute)
	v.Record(5, 3)
	v.Reset()
	snap := v.Snapshot()
	if snap.Opened != 0 || snap.Closed != 0 {
		t.Fatalf("expected zero after reset, got %+v", snap)
	}
}

func TestVelocityWindowReturned(t *testing.T) {
	win := 30 * time.Second
	v := NewVelocityTracker(win)
	snap := v.Snapshot()
	if snap.Window != win {
		t.Errorf("expected window %v, got %v", win, snap.Window)
	}
}
