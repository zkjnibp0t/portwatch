package ports

import (
	"testing"
	"time"
)

func fixedLatencyClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestLatencyAverageEmptyReturnsZero(t *testing.T) {
	tr := NewLatencyTracker(time.Minute)
	_, ok := tr.Average(8080)
	if ok {
		t.Fatal("expected false for empty tracker")
	}
}

func TestLatencyAverageAfterRecord(t *testing.T) {
	now := time.Now()
	tr := NewLatencyTracker(time.Minute)
	tr.clock = fixedLatencyClock(now)
	tr.Record(8080, 10*time.Millisecond)
	tr.Record(8080, 20*time.Millisecond)
	avg, ok := tr.Average(8080)
	if !ok {
		t.Fatal("expected ok")
	}
	if avg != 15*time.Millisecond {
		t.Fatalf("expected 15ms, got %v", avg)
	}
}

func TestLatencyEvictsOldSamples(t *testing.T) {
	now := time.Now()
	tr := NewLatencyTracker(30 * time.Second)
	tr.clock = fixedLatencyClock(now)
	tr.Record(9090, 50*time.Millisecond)
	// advance clock beyond window
	tr.clock = fixedLatencyClock(now.Add(31 * time.Second))
	tr.Record(9090, 5*time.Millisecond)
	avg, ok := tr.Average(9090)
	if !ok {
		t.Fatal("expected ok")
	}
	if avg != 5*time.Millisecond {
		t.Fatalf("expected 5ms after eviction, got %v", avg)
	}
}

func TestLatencySamplesReturnsCopy(t *testing.T) {
	now := time.Now()
	tr := NewLatencyTracker(time.Minute)
	tr.clock = fixedLatencyClock(now)
	tr.Record(443, 8*time.Millisecond)
	samples := tr.Samples(443)
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}
	samples[0].Latency = 999 * time.Second
	// original must be unchanged
	orig := tr.Samples(443)
	if orig[0].Latency == 999*time.Second {
		t.Fatal("Samples must return a copy, not a reference")
	}
}

func TestLatencyDifferentPortsAreIndependent(t *testing.T) {
	now := time.Now()
	tr := NewLatencyTracker(time.Minute)
	tr.clock = fixedLatencyClock(now)
	tr.Record(80, 10*time.Millisecond)
	tr.Record(443, 40*time.Millisecond)
	a, _ := tr.Average(80)
	b, _ := tr.Average(443)
	if a == b {
		t.Fatalf("expected different averages, both %v", a)
	}
}
