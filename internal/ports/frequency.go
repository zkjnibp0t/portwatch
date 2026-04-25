package ports

import (
	"sync"
	"time"
)

// FrequencyTracker tracks how often each port appears across scans
// within a rolling time window, allowing detection of ports that
// open and close with high frequency.
type FrequencyTracker struct {
	mu      sync.Mutex
	window  time.Duration
	events  map[int][]time.Time
	clock   func() time.Time
}

// NewFrequencyTracker creates a FrequencyTracker with the given rolling window.
func NewFrequencyTracker(window time.Duration) *FrequencyTracker {
	return &FrequencyTracker{
		window: window,
		events: make(map[int][]time.Time),
		clock:  time.Now,
	}
}

// Record registers an observation of the given port at the current time.
func (f *FrequencyTracker) Record(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	f.events[port] = append(f.prune(f.events[port], now), now)
}

// Count returns the number of observations for the given port within the window.
func (f *FrequencyTracker) Count(port int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	pruned := f.prune(f.events[port], now)
	f.events[port] = pruned
	return len(pruned)
}

// TopN returns the top n ports by observation count within the window.
// The result is a slice of (port, count) pairs sorted descending by count.
func (f *FrequencyTracker) TopN(n int) []PortCount {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	var results []PortCount
	for port, times := range f.events {
		pruned := f.prune(times, now)
		f.events[port] = pruned
		if len(pruned) > 0 {
			results = append(results, PortCount{Port: port, Count: len(pruned)})
		}
	}
	sortPortCounts(results)
	if n > 0 && len(results) > n {
		return results[:n]
	}
	return results
}

// Reset clears all recorded observations.
func (f *FrequencyTracker) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = make(map[int][]time.Time)
}

func (f *FrequencyTracker) prune(times []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-f.window)
	for i, t := range times {
		if t.After(cutoff) {
			return times[i:]
		}
	}
	return nil
}

// PortCount holds a port and its observation count.
type PortCount struct {
	Port  int
	Count int
}

func sortPortCounts(pcs []PortCount) {
	for i := 1; i < len(pcs); i++ {
		for j := i; j > 0 && pcs[j].Count > pcs[j-1].Count; j-- {
			pcs[j], pcs[j-1] = pcs[j-1], pcs[j]
		}
	}
}
