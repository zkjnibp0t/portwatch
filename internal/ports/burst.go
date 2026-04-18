package ports

import (
	"sync"
	"time"
)

// BurstDetector flags ports that open and close repeatedly in a short window.
type BurstDetector struct {
	mu       sync.Mutex
	events   map[int][]time.Time
	window   time.Duration
	threshold int
	clock    func() time.Time
}

// NewBurstDetector creates a BurstDetector that fires when a port crosses
// threshold events within window.
func NewBurstDetector(window time.Duration, threshold int) *BurstDetector {
	return &BurstDetector{
		events:    make(map[int][]time.Time),
		window:    window,
		threshold: threshold,
		clock:     time.Now,
	}
}

// Record registers an open or close event for port and returns true if the
// burst threshold has been exceeded within the configured window.
func (b *BurstDetector) Record(port int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	cutoff := now.Add(-b.window)

	prev := b.events[port]
	filtered := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	b.events[port] = filtered

	return len(filtered) >= b.threshold
}

// Reset clears the event history for a specific port.
func (b *BurstDetector) Reset(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.events, port)
}

// Count returns the number of recent events recorded for port within the window.
func (b *BurstDetector) Count(port int) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	cutoff := now.Add(-b.window)
	count := 0
	for _, t := range b.events[port] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
