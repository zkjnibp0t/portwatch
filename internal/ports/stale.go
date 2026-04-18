package ports

import "time"

// StaleTracker flags ports that have been open continuously beyond a threshold.
type StaleTracker struct {
	firstSeen map[int]time.Time
	threshold time.Duration
	clock     func() time.Time
}

// NewStaleTracker creates a StaleTracker with the given staleness threshold.
func NewStaleTracker(threshold time.Duration) *StaleTracker {
	return &StaleTracker{
		firstSeen: make(map[int]time.Time),
		threshold: threshold,
		clock:     time.Now,
	}
}

// Observe records that a port is currently open.
func (s *StaleTracker) Observe(port int) {
	if _, ok := s.firstSeen[port]; !ok {
		s.firstSeen[port] = s.clock()
	}
}

// Remove stops tracking a port (e.g. it closed).
func (s *StaleTracker) Remove(port int) {
	delete(s.firstSeen, port)
}

// IsStale returns true if the port has been open longer than the threshold.
func (s *StaleTracker) IsStale(port int) bool {
	t, ok := s.firstSeen[port]
	if !ok {
		return false
	}
	return s.clock().Sub(t) >= s.threshold
}

// StalePorts returns all ports that exceed the staleness threshold.
func (s *StaleTracker) StalePorts() []int {
	var out []int
	for port := range s.firstSeen {
		if s.IsStale(port) {
			out = append(out, port)
		}
	}
	return out
}

// Age returns how long a port has been tracked. Returns 0 if not tracked.
func (s *StaleTracker) Age(port int) time.Duration {
	t, ok := s.firstSeen[port]
	if !ok {
		return 0
	}
	return s.clock().Sub(t)
}
