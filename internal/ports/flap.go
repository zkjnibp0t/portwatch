package ports

import (
	"sync"
	"time"
)

// FlapEvent records a single open/close transition for a port.
type FlapEvent struct {
	Port      int
	OpenedAt  time.Time
	ClosedAt  time.Time
}

// FlapDetector tracks ports that open and close repeatedly within a window,
// flagging them as "flapping" when transitions exceed a threshold.
type FlapDetector struct {
	mu        sync.Mutex
	events    map[int][]FlapEvent
	window    time.Duration
	threshold int
	clock     func() time.Time
}

// NewFlapDetector creates a FlapDetector. threshold is the minimum number of
// open→close cycles within window to be considered flapping.
func NewFlapDetector(window time.Duration, threshold int, clock func() time.Time) *FlapDetector {
	if clock == nil {
		clock = time.Now
	}
	return &FlapDetector{
		events:    make(map[int][]FlapEvent),
		window:    window,
		threshold: threshold,
		clock:     clock,
	}
}

// RecordOpen notes that port was opened at the current time.
func (f *FlapDetector) RecordOpen(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events[port] = append(f.events[port], FlapEvent{Port: port, OpenedAt: f.clock()})
}

// RecordClose closes the most recent open event for port and prunes old events.
func (f *FlapDetector) RecordClose(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	evs := f.events[port]
	for i := len(evs) - 1; i >= 0; i-- {
		if evs[i].ClosedAt.IsZero() {
			evs[i].ClosedAt = now
			break
		}
	}
	f.events[port] = f.prune(evs, now)
}

// IsFlapping returns true if port has completed at least threshold open/close
// cycles within the configured window.
func (f *FlapDetector) IsFlapping(port int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	cutoff := now.Add(-f.window)
	count := 0
	for _, ev := range f.events[port] {
		if !ev.ClosedAt.IsZero() && ev.OpenedAt.After(cutoff) {
			count++
		}
	}
	return count >= f.threshold
}

// FlappingPorts returns the set of ports currently considered flapping.
func (f *FlapDetector) FlappingPorts() []int {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	cutoff := now.Add(-f.window)
	var result []int
	for port, evs := range f.events {
		count := 0
		for _, ev := range evs {
			if !ev.ClosedAt.IsZero() && ev.OpenedAt.After(cutoff) {
				count++
			}
		}
		if count >= f.threshold {
			result = append(result, port)
		}
	}
	return result
}

func (f *FlapDetector) prune(evs []FlapEvent, now time.Time) []FlapEvent {
	cutoff := now.Add(-f.window)
	out := evs[:0]
	for _, ev := range evs {
		if ev.OpenedAt.After(cutoff) {
			out = append(out, ev)
		}
	}
	return out
}
