package ports

import (
	"testing"
	"time"
)

func fixedReopenClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestReopenNotDetectedWithoutClose(t *testing.T) {
	d := NewReopenDetector(5 * time.Second)
	event := d.RecordOpened(8080)
	if event != nil {
		t.Fatalf("expected nil event, got %+v", event)
	}
}

func TestReopenDetectedWithinWindow(t *testing.T) {
	now := time.Now()
	d := NewReopenDetector(10 * time.Second)
	d.clock = fixedReopenClock(now)

	d.RecordClosed(8080)
	d.clock = fixedReopenClock(now.Add(3 * time.Second))

	event := d.RecordOpened(8080)
	if event == nil {
		t.Fatal("expected reopen event, got nil")
	}
	if event.Port != 8080 {
		t.Errorf("expected port 8080, got %d", event.Port)
	}
	if event.Gap != 3*time.Second {
		t.Errorf("expected gap 3s, got %v", event.Gap)
	}
}

func TestReopenNotDetectedOutsideWindow(t *testing.T) {
	now := time.Now()
	d := NewReopenDetector(5 * time.Second)
	d.clock = fixedReopenClock(now)

	d.RecordClosed(9000)
	d.clock = fixedReopenClock(now.Add(10 * time.Second))

	event := d.RecordOpened(9000)
	if event != nil {
		t.Fatalf("expected nil event outside window, got %+v", event)
	}
}

func TestReopenDifferentPortsAreIndependent(t *testing.T) {
	now := time.Now()
	d := NewReopenDetector(10 * time.Second)
	d.clock = fixedReopenClock(now)

	d.RecordClosed(80)
	d.clock = fixedReopenClock(now.Add(2 * time.Second))

	if d.RecordOpened(443) != nil {
		t.Error("port 443 should not trigger reopen event")
	}
	if d.RecordOpened(80) == nil {
		t.Error("port 80 should trigger reopen event")
	}
}

func TestReopenPurgeRemovesStaleEntries(t *testing.T) {
	now := time.Now()
	d := NewReopenDetector(5 * time.Second)
	d.clock = fixedReopenClock(now)

	d.RecordClosed(3000)
	d.clock = fixedReopenClock(now.Add(10 * time.Second))
	d.Purge()

	// After purge, reopening should not produce an event
	if d.RecordOpened(3000) != nil {
		t.Error("expected nil after purge")
	}
}
