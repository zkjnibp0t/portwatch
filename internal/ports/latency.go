package ports

import (
	"sync"
	"time"
)

// LatencySample holds the round-trip probe duration for a single port.
type LatencySample struct {
	Port      int
	Latency   time.Duration
	RecordedAt time.Time
}

// LatencyTracker records per-port probe latencies and exposes summary stats.
type LatencyTracker struct {
	mu      sync.Mutex
	samples map[int][]LatencySample
	window  time.Duration
	clock   func() time.Time
}

// NewLatencyTracker creates a LatencyTracker that retains samples within window.
func NewLatencyTracker(window time.Duration) *LatencyTracker {
	return &LatencyTracker{
		samples: make(map[int][]LatencySample),
		window:  window,
		clock:   time.Now,
	}
}

// Record adds a latency observation for port.
func (t *LatencyTracker) Record(port int, latency time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.samples[port] = append(t.samples[port], LatencySample{
		Port:       port,
		Latency:    latency,
		RecordedAt: now,
	})
	t.evict(port, now)
}

// Average returns the mean latency for port over the retention window.
// Returns 0 and false if no samples exist.
func (t *LatencyTracker) Average(port int) (time.Duration, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.evict(port, now)
	samples := t.samples[port]
	if len(samples) == 0 {
		return 0, false
	}
	var total time.Duration
	for _, s := range samples {
		total += s.Latency
	}
	return total / time.Duration(len(samples)), true
}

// Samples returns a copy of all retained samples for port.
func (t *LatencyTracker) Samples(port int) []LatencySample {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.evict(port, now)
	out := make([]LatencySample, len(t.samples[port]))
	copy(out, t.samples[port])
	return out
}

// evict removes samples outside the retention window; must be called with lock held.
func (t *LatencyTracker) evict(port int, now time.Time) {
	cutoff := now.Add(-t.window)
	ss := t.samples[port]
	start := 0
	for start < len(ss) && ss[start].RecordedAt.Before(cutoff) {
		start++
	}
	t.samples[port] = ss[start:]
}
