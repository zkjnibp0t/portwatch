package ports

import (
	"testing"
	"time"
)

func fixedBurstClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestBurstNotDetectedBelowThreshold(t *testing.T) {
	b := NewBurstDetector(10*time.Second, 3)
	if b.Record(80) {
		t.Fatal("expected no burst on first event")
	}
	if b.Record(80) {
		t.Fatal("expected no burst on second event")
	}
}

func TestBurstDetectedAtThreshold(t *testing.T) {
	b := NewBurstDetector(10*time.Second, 3)
	b.Record(80)
	b.Record(80)
	if !b.Record(80) {
		t.Fatal("expected burst detected on third event")
	}
}

func TestBurstEventsExpireOutsideWindow(t *testing.T) {
	now := time.Now()
	b := NewBurstDetector(5*time.Second, 3)
	b.clock = fixedBurstClock(now)
	b.Record(443)
	b.Record(443)
	// advance past window
	b.clock = fixedBurstClock(now.Add(6 * time.Second))
	// old events should be pruned; this is the first in the new window
	if b.Record(443) {
		t.Fatal("expected no burst after window expiry")
	}
}

func TestBurstDifferentPortsAreIndependent(t *testing.T) {
	b := NewBurstDetector(10*time.Second, 2)
	b.Record(80)
	b.Record(443)
	if b.Record(80) {
		t.Fatal("port 80 should not trigger burst independently")
	}
}

func TestBurstResetClearsHistory(t *testing.T) {
	b := NewBurstDetector(10*time.Second, 2)
	b.Record(22)
	b.Reset(22)
	if b.Count(22) != 0 {
		t.Fatalf("expected count 0 after reset, got %d", b.Count(22))
	}
	if b.Record(22) {
		t.Fatal("expected no burst after reset")
	}
}

func TestBurstCountReturnsRecent(t *testing.T) {
	b := NewBurstDetector(10*time.Second, 5)
	b.Record(8080)
	b.Record(8080)
	if got := b.Count(8080); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}
