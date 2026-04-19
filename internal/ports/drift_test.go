package ports

import (
	"testing"
	"time"
)

var fixedDriftBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedDriftClock(offset time.Duration) func() time.Time {
	return func() time.Time { return fixedDriftBase.Add(offset) }
}

func TestDriftNotDetectedBelowThreshold(t *testing.T) {
	var tick time.Duration
	clock := func() time.Time { return fixedDriftBase.Add(tick) }
	d := NewDriftDetector(10*time.Minute, clock)
	d.Observe([]int{8080}, nil)
	tick = 5 * time.Minute
	events := d.Detect()
	if len(events) != 0 {
		t.Fatalf("expected no drift, got %d", len(events))
	}
}

func TestDriftDetectedAfterThreshold(t *testing.T) {
	var tick time.Duration
	clock := func() time.Time { return fixedDriftBase.Add(tick) }
	d := NewDriftDetector(10*time.Minute, clock)
	d.Observe([]int{8080}, nil)
	tick = 15 * time.Minute
	events := d.Detect()
	if len(events) != 1 {
		t.Fatalf("expected 1 drift event, got %d", len(events))
	}
	if events[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Port)
	}
	if events[0].Drift < 10*time.Minute {
		t.Errorf("drift too small: %v", events[0].Drift)
	}
}

func TestDriftClosedPortRemoved(t *testing.T) {
	var tick time.Duration
	clock := func() time.Time { return fixedDriftBase.Add(tick) }
	d := NewDriftDetector(5*time.Minute, clock)
	d.Observe([]int{9090}, nil)
	tick = 10 * time.Minute
	d.Observe(nil, []int{9090})
	events := d.Detect()
	if len(events) != 0 {
		t.Fatalf("expected no events after close, got %d", len(events))
	}
}

func TestDriftMultiplePorts(t *testing.T) {
	var tick time.Duration
	clock := func() time.Time { return fixedDriftBase.Add(tick) }
	d := NewDriftDetector(10*time.Minute, clock)
	d.Observe([]int{80, 443, 8080}, nil)
	tick = 20 * time.Minute
	events := d.Detect()
	if len(events) != 3 {
		t.Fatalf("expected 3 drift events, got %d", len(events))
	}
}

func TestDriftResetClearsAll(t *testing.T) {
	var tick time.Duration
	clock := func() time.Time { return fixedDriftBase.Add(tick) }
	d := NewDriftDetector(5*time.Minute, clock)
	d.Observe([]int{22}, nil)
	tick = 10 * time.Minute
	d.Reset()
	events := d.Detect()
	if len(events) != 0 {
		t.Fatalf("expected no events after reset, got %d", len(events))
	}
}
