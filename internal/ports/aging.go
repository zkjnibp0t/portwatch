package ports

import (
	"sync"
	"time"
)

// AgingEntry tracks when a port was first and last seen.
type AgingEntry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Port      int
}

// AgingTracker records how long ports have been continuously open.
type AgingTracker struct {
	mu      sync.Mutex
	entries map[int]*AgingEntry
	clock   func() time.Time
}

// NewAgingTracker returns an AgingTracker using the real clock.
func NewAgingTracker() *AgingTracker {
	return &AgingTracker{
		entries: make(map[int]*AgingEntry),
		clock:   time.Now,
	}
}

// Observe records that the given ports are currently open.
func (a *AgingTracker) Observe(ports []int) {
	now := a.clock()
	a.mu.Lock()
	defer a.mu.Unlock()

	seen := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		seen[p] = struct{}{}
		if e, ok := a.entries[p]; ok {
			e.LastSeen = now
		} else {
			a.entries[p] = &AgingEntry{Port: p, FirstSeen: now, LastSeen: now}
		}
	}

	// Remove ports no longer open.
	for p := range a.entries {
		if _, ok := seen[p]; !ok {
			delete(a.entries, p)
		}
	}
}

// Age returns how long a port has been continuously open.
// Returns 0 and false if the port is not tracked.
func (a *AgingTracker) Age(port int) (time.Duration, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if e, ok := a.entries[port]; ok {
		return a.clock().Sub(e.FirstSeen), true
	}
	return 0, false
}

// OlderThan returns all ports that have been open longer than the given duration.
func (a *AgingTracker) OlderThan(d time.Duration) []AgingEntry {
	now := a.clock()
	a.mu.Lock()
	defer a.mu.Unlock()
	var result []AgingEntry
	for _, e := range a.entries {
		if now.Sub(e.FirstSeen) > d {
			result = append(result, *e)
		}
	}
	return result
}
