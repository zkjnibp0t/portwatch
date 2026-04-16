package ports

import (
	"fmt"
	"sync"
	"time"
)

// TrendEntry records a port change event with a timestamp.
type TrendEntry struct {
	Port      int
	Event     string // "opened" or "closed"
	Timestamp time.Time
}

// TrendTracker accumulates port change events over time,
// allowing callers to query how frequently a port flaps.
type TrendTracker struct {
	mu      sync.Mutex
	entries []TrendEntry
	maxAge  time.Duration
}

// NewTrendTracker creates a TrendTracker that discards entries
// older than maxAge on each Record call.
func NewTrendTracker(maxAge time.Duration) *TrendTracker {
	if maxAge <= 0 {
		maxAge = 24 * time.Hour
	}
	return &TrendTracker{maxAge: maxAge}
}

// Record adds a new event for the given port.
func (t *TrendTracker) Record(port int, event string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	t.entries = append(t.entries, TrendEntry{Port: port, Event: event, Timestamp: now})
	t.prune(now)
}

// FlapCount returns how many times a port has changed state within maxAge.
func (t *TrendTracker) FlapCount(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	count := 0
	for _, e := range t.entries {
		if e.Port == port {
			count++
		}
	}
	return count
}

// Summary returns a human-readable summary of the top flapping ports.
func (t *TrendTracker) Summary() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	counts := make(map[int]int)
	for _, e := range t.entries {
		counts[e.Port]++
	}
	var lines []string
	for port, n := range counts {
		lines = append(lines, fmt.Sprintf("port %d: %d event(s)", port, n))
	}
	return lines
}

// RecentEvents returns all recorded events for the given port within maxAge,
// ordered from oldest to newest.
func (t *TrendTracker) RecentEvents(port int) []TrendEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	var result []TrendEntry
	for _, e := range t.entries {
		if e.Port == port {
			result = append(result, e)
		}
	}
	return result
}

// prune removes entries older than maxAge. Caller must hold mu.
func (t *TrendTracker) prune(now time.Time) {
	cutoff := now.Add(-t.maxAge)
	filtered := t.entries[:0]
	for _, e := range t.entries {
		if e.Timestamp.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	t.entries = filtered
}
