package ports

import (
	"testing"
	"time"
)

func fixedPressureClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestPressureNotUnderPressureInitially(t *testing.T) {
	pt := NewPressureTracker(10*time.Second, 3)
	if pt.UnderPressure() {
		t.Fatal("expected no pressure initially")
	}
}

func TestPressureUnderPressureAtThreshold(t *testing.T) {
	now := time.Now()
	pt := NewPressureTracker(10*time.Second, 3)
	pt.clock = fixedPressureClock(now)
	pt.Record(3)
	if !pt.UnderPressure() {
		t.Fatal("expected pressure at threshold")
	}
}

func TestPressureBelowThreshold(t *testing.T) {
	now := time.Now()
	pt := NewPressureTracker(10*time.Second, 5)
	pt.clock = fixedPressureClock(now)
	pt.Record(4)
	if pt.UnderPressure() {
		t.Fatal("expected no pressure below threshold")
	}
	if pt.Count() != 4 {
		t.Fatalf("expected count 4, got %d", pt.Count())
	}
}

func TestPressureEvictsOldEvents(t *testing.T) {
	now := time.Now()
	pt := NewPressureTracker(5*time.Second, 3)
	pt.clock = fixedPressureClock(now)
	pt.Record(3)
	// advance clock beyond window
	pt.clock = fixedPressureClock(now.Add(6 * time.Second))
	if pt.UnderPressure() {
		t.Fatal("expected old events to be evicted")
	}
	if pt.Count() != 0 {
		t.Fatalf("expected count 0 after eviction, got %d", pt.Count())
	}
}

func TestPressureResetClearsEvents(t *testing.T) {
	now := time.Now()
	pt := NewPressureTracker(10*time.Second, 2)
	pt.clock = fixedPressureClock(now)
	pt.Record(5)
	pt.Reset()
	if pt.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", pt.Count())
	}
}

func TestPressureAccumulatesAcrossRecords(t *testing.T) {
	now := time.Now()
	pt := NewPressureTracker(10*time.Second, 5)
	pt.clock = fixedPressureClock(now)
	pt.Record(2)
	pt.Record(3)
	if pt.Count() != 5 {
		t.Fatalf("expected count 5, got %d", pt.Count())
	}
}
