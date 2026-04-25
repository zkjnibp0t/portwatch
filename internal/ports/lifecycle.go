package ports

import (
	"sync"
	"time"
)

// LifecycleEvent represents a state transition for a port.
type LifecycleEvent struct {
	Port      int
	PrevState string
	NextState string
	At        time.Time
	Duration  time.Duration // time spent in previous state
}

// LifecycleTracker records open/close transitions and computes time spent in each state.
type LifecycleTracker struct {
	mu      sync.Mutex
	clockFn func() time.Time
	open    map[int]time.Time // port -> time it was opened
	events  []LifecycleEvent
}

// NewLifecycleTracker returns a new LifecycleTracker.
func NewLifecycleTracker(clockFn func() time.Time) *LifecycleTracker {
	if clockFn == nil {
		clockFn = time.Now
	}
	return &LifecycleTracker{
		clockFn: clockFn,
		open:    make(map[int]time.Time),
	}
}

// RecordOpened marks a port as opened at the current clock time.
func (l *LifecycleTracker) RecordOpened(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clockFn()
	if _, exists := l.open[port]; !exists {
		l.open[port] = now
		l.events = append(l.events, LifecycleEvent{
			Port:      port,
			PrevState: "closed",
			NextState: "open",
			At:        now,
		})
	}
}

// RecordClosed marks a port as closed and records how long it was open.
func (l *LifecycleTracker) RecordClosed(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clockFn()
	var dur time.Duration
	if openedAt, exists := l.open[port]; exists {
		dur = now.Sub(openedAt)
		delete(l.open, port)
	}
	l.events = append(l.events, LifecycleEvent{
		Port:      port,
		PrevState: "open",
		NextState: "closed",
		At:        now,
		Duration:  dur,
	})
}

// Events returns a copy of all recorded lifecycle events.
func (l *LifecycleTracker) Events() []LifecycleEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LifecycleEvent, len(l.events))
	copy(out, l.events)
	return out
}

// OpenSince returns the time a port was opened, and whether it is currently open.
func (l *LifecycleTracker) OpenSince(port int) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.open[port]
	return t, ok
}
