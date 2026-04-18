package ports

import (
	"sync"
	"time"
)

// WindowAggregator accumulates port sets over a sliding time window
// and returns the union of all sets seen within that window.
type WindowAggregator struct {
	mu      sync.Mutex
	window  time.Duration
	buckets []windowBucket
	clock   func() time.Time
}

type windowBucket struct {
	at  time.Time
	set map[int]struct{}
}

// NewWindowAggregator creates a WindowAggregator with the given window duration.
func NewWindowAggregator(window time.Duration) *WindowAggregator {
	return &WindowAggregator{
		window: window,
		clock:  time.Now,
	}
}

// Add records a port set at the current time.
func (w *WindowAggregator) Add(ports map[int]struct{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	copy := make(map[int]struct{}, len(ports))
	for p := range ports {
		copy[p] = struct{}{}
	}
	w.buckets = append(w.buckets, windowBucket{at: w.clock(), set: copy})
	w.evict()
}

// Union returns the union of all port sets within the current window.
func (w *WindowAggregator) Union() map[int]struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	result := make(map[int]struct{})
	for _, b := range w.buckets {
		for p := range b.set {
			result[p] = struct{}{}
		}
	}
	return result
}

// Len returns the number of buckets currently in the window.
func (w *WindowAggregator) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	return len(w.buckets)
}

func (w *WindowAggregator) evict() {
	cutoff := w.clock().Add(-w.window)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
