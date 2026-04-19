package ports

import (
	"sync"
	"time"
)

// PressureTracker measures how many ports opened within a sliding window
// and reports whether the system is under "port pressure".
type PressureTracker struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    []time.Time
	clock     func() time.Time
}

func NewPressureTracker(window time.Duration, threshold int) *PressureTracker {
	return &PressureTracker{
		window:    window,
		threshold: threshold,
		clock:     time.Now,
	}
}

// Record adds n open-port events at the current time.
func (p *PressureTracker) Record(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.clock()
	for i := 0; i < n; i++ {
		p.events = append(p.events, now)
	}
	p.evict(now)
}

// UnderPressure returns true when the number of events within the window
// meets or exceeds the threshold.
func (p *PressureTracker) UnderPressure() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.evict(p.clock())
	return len(p.events) >= p.threshold
}

// Count returns the number of events currently within the window.
func (p *PressureTracker) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.evict(p.clock())
	return len(p.events)
}

// Reset clears all recorded events.
func (p *PressureTracker) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = p.events[:0]
}

func (p *PressureTracker) evict(now time.Time) {
	cutoff := now.Add(-p.window)
	i := 0
	for i < len(p.events) && p.events[i].Before(cutoff) {
		i++
	}
	p.events = p.events[i:]
}
