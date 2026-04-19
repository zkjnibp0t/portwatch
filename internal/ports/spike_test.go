package ports

import (
	"testing"
	"time"
)

var fixedSpikeBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newSpikeDetector(window time.Duration) (*SpikeDetector, *time.Time) {
	now := fixedSpikeBase
	sd := NewSpikeDetector(window)
	sd.clock = func() time.Time { return now }
	return sd, &now
}

func TestSpikeNotDetectedWithoutOpen(t *testing.T) {
	sd, _ := newSpikeDetector(5 * time.Second)
	_, ok := sd.RecordClosed(8080)
	if ok {
		t.Fatal("expected no spike for port that was never opened")
	}
}

func TestSpikeDetectedWithinWindow(t *testing.T) {
	sd, now := newSpikeDetector(5 * time.Second)
	sd.RecordOpened(8080)
	*now = now.Add(2 * time.Second)
	ev, ok := sd.RecordClosed(8080)
	if !ok {
		t.Fatal("expected spike to be detected")
	}
	if ev.Port != 8080 {
		t.Errorf("expected port 8080, got %d", ev.Port)
	}
	if ev.Duration != 2*time.Second {
		t.Errorf("unexpected duration %v", ev.Duration)
	}
}

func TestSpikeNotDetectedOutsideWindow(t *testing.T) {
	sd, now := newSpikeDetector(5 * time.Second)
	sd.RecordOpened(9090)
	*now = now.Add(10 * time.Second)
	_, ok := sd.RecordClosed(9090)
	if ok {
		t.Fatal("expected no spike outside window")
	}
}

func TestSpikeAccumulates(t *testing.T) {
	sd, now := newSpikeDetector(5 * time.Second)
	for _, p := range []int{80, 443, 22} {
		sd.RecordOpened(p)
		*now = now.Add(1 * time.Second)
		sd.RecordClosed(p)
	}
	if len(sd.Spikes()) != 3 {
		t.Errorf("expected 3 spikes, got %d", len(sd.Spikes()))
	}
}

func TestSpikeResetClearsAll(t *testing.T) {
	sd, now := newSpikeDetector(5 * time.Second)
	sd.RecordOpened(1234)
	*now = now.Add(1 * time.Second)
	sd.RecordClosed(1234)
	sd.Reset()
	if len(sd.Spikes()) != 0 {
		t.Fatal("expected empty spikes after reset")
	}
	_, ok := sd.RecordClosed(1234)
	if ok {
		t.Fatal("expected no spike after reset")
	}
}
