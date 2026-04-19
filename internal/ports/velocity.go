package ports

import (
	"sync"
	"time"
)

// VelocityTracker measures how quickly ports are opening or closing
// over a sliding time window.
type VelocityTracker struct {
	mu      sync.Mutex
	window  time.Duration
	events  []velocityEvent
	clock   func() time.Time
}

type velocityEvent struct {
	at     time.Time
	opened int
	closed int
}

// VelocitySnapshot holds the computed rate for a window.
type VelocitySnapshot struct {
	Opened int
	Closed int
	Window time.Duration
}

// NewVelocityTracker creates a tracker with the given sliding window.
func NewVelocityTracker(window time.Duration) *VelocityTracker {
	return &VelocityTracker{
		window: window,
		clock:  time.Now,
	}
}

// Record adds a new observation of opened/closed port counts.
func (v *VelocityTracker) Record(opened, closed int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.prune()
	v.events = append(v.events, velocityEvent{
		at:     v.clock(),
		opened: opened,
		closed: closed,
	})
}

// Snapshot returns the aggregated velocity over the current window.
func (v *VelocityTracker) Snapshot() VelocitySnapshot {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.prune()
	snap := VelocitySnapshot{Window: v.window}
	for _, e := range v.events {
		snap.Opened += e.opened
		snap.Closed += e.closed
	}
	return snap
}

// Reset clears all recorded events.
func (v *VelocityTracker) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.events = nil
}

func (v *VelocityTracker) prune() {
	cutoff := v.clock().Add(-v.window)
	i := 0
	for i < len(v.events) && v.events[i].at.Before(cutoff) {
		i++
	}
	v.events = v.events[i:]
}
