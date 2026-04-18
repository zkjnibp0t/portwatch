package ports

import (
	"sync"
	"time"
)

// DecayTracker tracks port activity scores that decay over time.
// Each observation increments a score; idle ports decay toward zero.
type DecayTracker struct {
	mu      sync.Mutex
	scores  map[int]float64
	lastSeen map[int]time.Time
	decayRate float64 // fraction lost per second
	clock    func() time.Time
}

func NewDecayTracker(decayRate float64) *DecayTracker {
	if decayRate <= 0 {
		decayRate = 0.1
	}
	return &DecayTracker{
		scores:    make(map[int]float64),
		lastSeen:  make(map[int]time.Time),
		decayRate: decayRate,
		clock:     time.Now,
	}
}

// Observe records an activity event for a port, incrementing its score.
func (d *DecayTracker) Observe(port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	d.applyDecay(port, now)
	d.scores[port] += 1.0
	d.lastSeen[port] = now
}

// Score returns the current decayed score for a port.
func (d *DecayTracker) Score(port int) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	d.applyDecay(port, now)
	return d.scores[port]
}

// Prune removes ports whose score has fallen below the threshold.
func (d *DecayTracker) Prune(threshold float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	for port := range d.scores {
		d.applyDecay(port, now)
		if d.scores[port] < threshold {
			delete(d.scores, port)
			delete(d.lastSeen, port)
		}
	}
}

func (d *DecayTracker) applyDecay(port int, now time.Time) {
	last, ok := d.lastSeen[port]
	if !ok {
		return
	}
	elapsed := now.Sub(last).Seconds()
	if elapsed <= 0 {
		return
	}
	factor := 1.0 - d.decayRate*elapsed
	if factor < 0 {
		factor = 0
	}
	d.scores[port] *= factor
	d.lastSeen[port] = now
}
