package ports

import (
	"testing"
	"time"
)

func fixedDecayClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestDecayScoreZeroInitially(t *testing.T) {
	d := NewDecayTracker(0.1)
	if s := d.Score(8080); s != 0 {
		t.Fatalf("expected 0, got %f", s)
	}
}

func TestDecayObserveIncrementsScore(t *testing.T) {
	now := time.Now()
	d := NewDecayTracker(0.1)
	d.clock = fixedDecayClock(now)
	d.Observe(8080)
	d.Observe(8080)
	if s := d.Score(8080); s != 2.0 {
		t.Fatalf("expected 2.0, got %f", s)
	}
}

func TestDecayScoreDecreasesOverTime(t *testing.T) {
	now := time.Now()
	d := NewDecayTracker(0.5) // 50% per second
	d.clock = fixedDecayClock(now)
	d.Observe(9000)
	// advance 1 second
	d.clock = fixedDecayClock(now.Add(time.Second))
	s := d.Score(9000)
	// score should be ~0.5
	if s >= 1.0 || s <= 0.0 {
		t.Fatalf("expected score between 0 and 1, got %f", s)
	}
}

func TestDecayPruneRemovesBelowThreshold(t *testing.T) {
	now := time.Now()
	d := NewDecayTracker(1.0) // full decay per second
	d.clock = fixedDecayClock(now)
	d.Observe(1234)
	d.clock = fixedDecayClock(now.Add(2 * time.Second))
	d.Prune(0.01)
	if s := d.Score(1234); s != 0 {
		t.Fatalf("expected pruned port score 0, got %f", s)
	}
}

func TestDecayDifferentPortsAreIndependent(t *testing.T) {
	now := time.Now()
	d := NewDecayTracker(0.1)
	d.clock = fixedDecayClock(now)
	d.Observe(80)
	d.Observe(443)
	d.Observe(443)
	if d.Score(80) >= d.Score(443) {
		t.Fatalf("expected port 443 score higher than 80")
	}
}

func TestDecayDefaultRateOnZero(t *testing.T) {
	d := NewDecayTracker(0)
	if d.decayRate != 0.1 {
		t.Fatalf("expected default decay rate 0.1, got %f", d.decayRate)
	}
}
