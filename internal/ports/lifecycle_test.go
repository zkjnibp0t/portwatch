package ports

import (
	"testing"
	"time"
)

func fixedLifecycleClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestLifecycleRecordOpened(t *testing.T) {
	now := time.Unix(1000, 0)
	tracker := NewLifecycleTracker(fixedLifecycleClock(now))
	tracker.RecordOpened(8080)

	at, ok := tracker.OpenSince(8080)
	if !ok {
		t.Fatal("expected port 8080 to be open")
	}
	if !at.Equal(now) {
		t.Errorf("expected open time %v, got %v", now, at)
	}
}

func TestLifecycleRecordClosed(t *testing.T) {
	base := time.Unix(1000, 0)
	var tick int64
	tracker := NewLifecycleTracker(func() time.Time {
		t := base.Add(time.Duration(tick) * time.Second)
		tick++
		return t
	})
	tracker.RecordOpened(9090) // t=1000
	tracker.RecordClosed(9090) // t=1001

	_, ok := tracker.OpenSince(9090)
	if ok {
		t.Error("expected port 9090 to be closed")
	}

	events := tracker.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	closed := events[1]
	if closed.Duration != time.Second {
		t.Errorf("expected duration 1s, got %v", closed.Duration)
	}
	if closed.PrevState != "open" || closed.NextState != "closed" {
		t.Errorf("unexpected states: %s -> %s", closed.PrevState, closed.NextState)
	}
}

func TestLifecycleOpenedEventState(t *testing.T) {
	now := time.Unix(500, 0)
	tracker := NewLifecycleTracker(fixedLifecycleClock(now))
	tracker.RecordOpened(443)

	events := tracker.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].PrevState != "closed" || events[0].NextState != "open" {
		t.Errorf("unexpected states: %s -> %s", events[0].PrevState, events[0].NextState)
	}
}

func TestLifecycleDuplicateOpenIgnored(t *testing.T) {
	now := time.Unix(1000, 0)
	tracker := NewLifecycleTracker(fixedLifecycleClock(now))
	tracker.RecordOpened(80)
	tracker.RecordOpened(80) // duplicate

	events := tracker.Events()
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestLifecycleClosedWithoutOpenHasDurationZero(t *testing.T) {
	now := time.Unix(1000, 0)
	tracker := NewLifecycleTracker(fixedLifecycleClock(now))
	tracker.RecordClosed(22) // never opened

	events := tracker.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Duration != 0 {
		t.Errorf("expected zero duration, got %v", events[0].Duration)
	}
}
