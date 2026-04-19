package ports

import "sync"

// SequenceTracker assigns monotonically increasing sequence numbers to
// port-change events so consumers can detect gaps or reordering.
type SequenceTracker struct {
	mu      sync.Mutex
	counter uint64
	last    map[int]uint64 // port -> last sequence number assigned
}

// NewSequenceTracker returns an initialised SequenceTracker.
func NewSequenceTracker() *SequenceTracker {
	return &SequenceTracker{last: make(map[int]uint64)}
}

// Next returns the next global sequence number and records it against the port.
func (s *SequenceTracker) Next(port int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	s.last[port] = s.counter
	return s.counter
}

// LastFor returns the most recent sequence number assigned to port, or 0 if
// the port has never been seen.
func (s *SequenceTracker) LastFor(port int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last[port]
}

// Reset clears all state and resets the counter to zero.
func (s *SequenceTracker) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = 0
	s.last = make(map[int]uint64)
}

// Len returns the number of distinct ports that have been assigned a sequence.
func (s *SequenceTracker) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.last)
}
