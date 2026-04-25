package ports

import (
	"testing"
	"time"
)

func fixedFreqClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestFrequencyCountZeroInitially(t *testing.T) {
	ft := NewFrequencyTracker(time.Minute)
	if got := ft.Count(8080); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestFrequencyCountIncrementsOnRecord(t *testing.T) {
	base := time.Now()
	ft := NewFrequencyTracker(time.Minute)
	ft.clock = fixedFreqClock(base)
	ft.Record(8080)
	ft.Record(8080)
	if got := ft.Count(8080); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestFrequencyEvictsOldEvents(t *testing.T) {
	base := time.Now()
	ft := NewFrequencyTracker(time.Minute)
	ft.clock = fixedFreqClock(base)
	ft.Record(9090)
	ft.Record(9090)
	// advance past window
	ft.clock = fixedFreqClock(base.Add(2 * time.Minute))
	ft.Record(9090)
	if got := ft.Count(9090); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestFrequencyDifferentPortsAreIndependent(t *testing.T) {
	base := time.Now()
	ft := NewFrequencyTracker(time.Minute)
	ft.clock = fixedFreqClock(base)
	ft.Record(80)
	ft.Record(80)
	ft.Record(443)
	if got := ft.Count(80); got != 2 {
		t.Fatalf("port 80: expected 2, got %d", got)
	}
	if got := ft.Count(443); got != 1 {
		t.Fatalf("port 443: expected 1, got %d", got)
	}
}

func TestFrequencyTopN(t *testing.T) {
	base := time.Now()
	ft := NewFrequencyTracker(time.Minute)
	ft.clock = fixedFreqClock(base)
	for i := 0; i < 3; i++ {
		ft.Record(80)
	}
	for i := 0; i < 5; i++ {
		ft.Record(443)
	}
	ft.Record(8080)
	top := ft.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0].Port != 443 || top[0].Count != 5 {
		t.Errorf("expected 443 with count 5, got port=%d count=%d", top[0].Port, top[0].Count)
	}
	if top[1].Port != 80 || top[1].Count != 3 {
		t.Errorf("expected 80 with count 3, got port=%d count=%d", top[1].Port, top[1].Count)
	}
}

func TestFrequencyResetClearsAll(t *testing.T) {
	base := time.Now()
	ft := NewFrequencyTracker(time.Minute)
	ft.clock = fixedFreqClock(base)
	ft.Record(8080)
	ft.Record(8080)
	ft.Reset()
	if got := ft.Count(8080); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
