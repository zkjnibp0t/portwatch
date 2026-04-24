package ports

import (
	"testing"
	"time"
)

func fixedHeartbeatClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestHeartbeatBeatAndLastSeen(t *testing.T) {
	now := time.Unix(1_000, 0)
	h := NewHeartbeatTracker(5 * time.Second)
	h.clock = fixedHeartbeatClock(now)

	h.Beat(8080)
	got, ok := h.LastSeen(8080)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}

func TestHeartbeatNotSilentBeforeDeadline(t *testing.T) {
	now := time.Unix(1_000, 0)
	h := NewHeartbeatTracker(10 * time.Second)
	h.clock = fixedHeartbeatClock(now)
	h.Beat(443)

	// advance 5 s — still within deadline
	h.clock = fixedHeartbeatClock(now.Add(5 * time.Second))
	if s := h.Silent(); len(s) != 0 {
		t.Fatalf("expected no silent ports, got %v", s)
	}
}

func TestHeartbeatSilentAfterDeadline(t *testing.T) {
	now := time.Unix(1_000, 0)
	h := NewHeartbeatTracker(10 * time.Second)
	h.clock = fixedHeartbeatClock(now)
	h.Beat(443)

	// advance past deadline
	h.clock = fixedHeartbeatClock(now.Add(11 * time.Second))
	s := h.Silent()
	if len(s) != 1 || s[0] != 443 {
		t.Fatalf("expected [443], got %v", s)
	}
}

func TestHeartbeatRemoveClearsPort(t *testing.T) {
	now := time.Unix(1_000, 0)
	h := NewHeartbeatTracker(5 * time.Second)
	h.clock = fixedHeartbeatClock(now)
	h.Beat(22)
	h.Remove(22)

	_, ok := h.LastSeen(22)
	if ok {
		t.Fatal("expected port to be removed")
	}
	h.clock = fixedHeartbeatClock(now.Add(10 * time.Second))
	if s := h.Silent(); len(s) != 0 {
		t.Fatalf("expected no silent ports after remove, got %v", s)
	}
}

func TestHeartbeatMultiplePortsIndependent(t *testing.T) {
	now := time.Unix(1_000, 0)
	h := NewHeartbeatTracker(10 * time.Second)
	h.clock = fixedHeartbeatClock(now)
	h.Beat(80)

	// advance 5 s, beat port 443 fresh
	h.clock = fixedHeartbeatClock(now.Add(5 * time.Second))
	h.Beat(443)

	// advance 6 more s — port 80 is 11 s old, port 443 is 6 s old
	h.clock = fixedHeartbeatClock(now.Add(11 * time.Second))
	s := h.Silent()
	if len(s) != 1 || s[0] != 80 {
		t.Fatalf("expected only port 80 silent, got %v", s)
	}
}
