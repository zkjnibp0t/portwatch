package ports

import (
	"sync"
	"time"
)

// ShadowEntry records a port that appeared briefly and then closed.
type ShadowEntry struct {
	Port      int
	OpenedAt  time.Time
	ClosedAt  time.Time
	Duration  time.Duration
}

// ShadowTracker detects ports that open and close within a short window.
type ShadowTracker struct {
	mu       sync.Mutex
	window   time.Duration
	opened   map[int]time.Time
	shadows  []ShadowEntry
	clock    func() time.Time
}

// NewShadowTracker creates a ShadowTracker with the given flash window.
func NewShadowTracker(window time.Duration) *ShadowTracker {
	return &ShadowTracker{
		window:  window,
		opened:  make(map[int]time.Time),
		clock:   time.Now,
	}
}

// RecordOpened notes that a port was observed opening.
func (s *ShadowTracker) RecordOpened(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.opened[port] = s.clock()
}

// RecordClosed checks if the port closed within the shadow window.
func (s *ShadowTracker) RecordClosed(port int) (ShadowEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	openedAt, ok := s.opened[port]
	if !ok {
		return ShadowEntry{}, false
	}
	delete(s.opened, port)
	now := s.clock()
	dur := now.Sub(openedAt)
	if dur > s.window {
		return ShadowEntry{}, false
	}
	e := ShadowEntry{Port: port, OpenedAt: openedAt, ClosedAt: now, Duration: dur}
	s.shadows = append(s.shadows, e)
	return e, true
}

// Shadows returns all recorded shadow entries.
func (s *ShadowTracker) Shadows() []ShadowEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ShadowEntry, len(s.shadows))
	copy(out, s.shadows)
	return out
}

// Reset clears all state.
func (s *ShadowTracker) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.opened = make(map[int]time.Time)
	s.shadows = nil
}
