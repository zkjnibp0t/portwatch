package ports

import (
	"sync"
	"time"
)

// SpikeDetector flags ports that open and close within a short spike window.
type SpikeDetector struct {
	mu      sync.Mutex
	window  time.Duration
	opened  map[int]time.Time
	spikes  []SpikeEvent
	clock   func() time.Time
}

// SpikeEvent records a port that exhibited spike behaviour.
type SpikeEvent struct {
	Port     int
	OpenedAt time.Time
	ClosedAt time.Time
	Duration time.Duration
}

// NewSpikeDetector creates a SpikeDetector with the given spike window.
func NewSpikeDetector(window time.Duration) *SpikeDetector {
	return &SpikeDetector{
		window: window,
		opened: make(map[int]time.Time),
		clock:  time.Now,
	}
}

// RecordOpened marks a port as opened at the current time.
func (s *SpikeDetector) RecordOpened(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.opened[port] = s.clock()
}

// RecordClosed checks whether the port closed within the spike window.
// If so, a SpikeEvent is stored and true is returned.
func (s *SpikeDetector) RecordClosed(port int) (SpikeEvent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	openedAt, ok := s.opened[port]
	if !ok {
		return SpikeEvent{}, false
	}
	delete(s.opened, port)

	now := s.clock()
	dur := now.Sub(openedAt)
	if dur > s.window {
		return SpikeEvent{}, false
	}

	ev := SpikeEvent{Port: port, OpenedAt: openedAt, ClosedAt: now, Duration: dur}
	s.spikes = append(s.spikes, ev)
	return ev, true
}

// Spikes returns all recorded spike events.
func (s *SpikeDetector) Spikes() []SpikeEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SpikeEvent, len(s.spikes))
	copy(out, s.spikes)
	return out
}

// Reset clears all state.
func (s *SpikeDetector) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.opened = make(map[int]time.Time)
	s.spikes = nil
}
