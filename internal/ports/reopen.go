package ports

import (
	"sync"
	"time"
)

// ReopenEvent records a port that was closed and then reopened.
type ReopenEvent struct {
	Port      int
	ReopenedAt time.Time
	ClosedAt   time.Time
	Gap        time.Duration
}

// ReopenDetector tracks ports that close and reopen within a time window,
// which may indicate a flapping or restarting service.
type ReopenDetector struct {
	mu       sync.Mutex
	window   time.Duration
	closedAt map[int]time.Time
	clock    func() time.Time
}

// NewReopenDetector creates a ReopenDetector with the given observation window.
func NewReopenDetector(window time.Duration) *ReopenDetector {
	return &ReopenDetector{
		window:   window,
		closedAt: make(map[int]time.Time),
		clock:    time.Now,
	}
}

// RecordClosed notes that a port was closed at the current time.
func (r *ReopenDetector) RecordClosed(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closedAt[port] = r.clock()
}

// RecordOpened checks if a port was recently closed and returns a ReopenEvent
// if the port reopened within the observation window. Returns nil otherwise.
func (r *ReopenDetector) RecordOpened(port int) *ReopenEvent {
	r.mu.Lock()
	defer r.mu.Unlock()

	closed, ok := r.closedAt[port]
	if !ok {
		return nil
	}

	now := r.clock()
	gap := now.Sub(closed)
	if gap > r.window {
		delete(r.closedAt, port)
		return nil
	}

	delete(r.closedAt, port)
	return &ReopenEvent{
		Port:       port,
		ReopenedAt: now,
		ClosedAt:   closed,
		Gap:        gap,
	}
}

// Purge removes stale closed-port entries that are outside the window.
func (r *ReopenDetector) Purge() {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := r.clock().Add(-r.window)
	for port, t := range r.closedAt {
		if t.Before(cutoff) {
			delete(r.closedAt, port)
		}
	}
}
