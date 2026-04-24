package ports

import (
	"testing"
	"time"
)

var baseFlap = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedFlapClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestFlapNotDetectedBelowThreshold(t *testing.T) {
	clock := baseFlap
	f := NewFlapDetector(time.Minute, 3, fixedFlapClock(clock))

	f.RecordOpen(8080)
	clock = clock.Add(5 * time.Second)
	f.RecordClose(8080)

	if f.IsFlapping(8080) {
		t.Fatal("expected not flapping below threshold")
	}
}

func TestFlapDetectedAtThreshold(t *testing.T) {
	var now time.Time = baseFlap
	f := NewFlapDetector(time.Minute, 3, func() time.Time { return now })

	for i := 0; i < 3; i++ {
		f.RecordOpen(9000)
		now = now.Add(2 * time.Second)
		f.RecordClose(9000)
		now = now.Add(2 * time.Second)
	}

	if !f.IsFlapping(9000) {
		t.Fatal("expected port 9000 to be flapping")
	}
}

func TestFlapEventsExpireOutsideWindow(t *testing.T) {
	var now time.Time = baseFlap
	f := NewFlapDetector(30*time.Second, 2, func() time.Time { return now })

	// Two cycles far in the past
	for i := 0; i < 2; i++ {
		f.RecordOpen(443)
		now = now.Add(2 * time.Second)
		f.RecordClose(443)
		now = now.Add(2 * time.Second)
	}

	// Advance past window
	now = now.Add(60 * time.Second)

	if f.IsFlapping(443) {
		t.Fatal("expected flap events to have expired")
	}
}

func TestFlapDifferentPortsAreIndependent(t *testing.T) {
	var now time.Time = baseFlap
	f := NewFlapDetector(time.Minute, 2, func() time.Time { return now })

	for i := 0; i < 2; i++ {
		f.RecordOpen(80)
		now = now.Add(1 * time.Second)
		f.RecordClose(80)
		now = now.Add(1 * time.Second)
	}

	if f.IsFlapping(443) {
		t.Fatal("port 443 should not be affected by port 80 flaps")
	}
	if !f.IsFlapping(80) {
		t.Fatal("port 80 should be flapping")
	}
}

func TestFlappingPortsReturnsAll(t *testing.T) {
	var now time.Time = baseFlap
	f := NewFlapDetector(time.Minute, 2, func() time.Time { return now })

	for _, port := range []int{3000, 4000} {
		for i := 0; i < 2; i++ {
			f.RecordOpen(port)
			now = now.Add(1 * time.Second)
			f.RecordClose(port)
			now = now.Add(1 * time.Second)
		}
	}

	result := f.FlappingPorts()
	if len(result) != 2 {
		t.Fatalf("expected 2 flapping ports, got %d", len(result))
	}
}
