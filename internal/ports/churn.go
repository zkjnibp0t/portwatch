package ports

import (
	"sync"
	"time"
)

// ChurnRecord holds open/close counts for a port within the observation window.
type ChurnRecord struct {
	Port    int
	Opens   int
	Closes  int
	LastSeen time.Time
}

// ChurnRate is the total number of state transitions for a port.
func (r ChurnRecord) ChurnRate() int {
	return r.Opens + r.Closes
}

// ChurnTracker tracks how frequently individual ports open and close over a
// sliding time window. Ports that exceed the churn threshold are considered
// unstable and can be surfaced for operator review.
type ChurnTracker struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    map[int][]churnEvent
	clock     func() time.Time
}

type churnEvent struct {
	at    time.Time
	state string // "open" or "close"
}

// NewChurnTracker creates a ChurnTracker with the given sliding window and
// churn threshold. Ports with ChurnRate >= threshold are flagged as unstable.
func NewChurnTracker(window time.Duration, threshold int, clock func() time.Time) *ChurnTracker {
	if clock == nil {
		clock = time.Now
	}
	return &ChurnTracker{
		window:    window,
		threshold: threshold,
		events:    make(map[int][]churnEvent),
		clock:     clock,
	}
}

// RecordOpen records a port-open event.
func (c *ChurnTracker) RecordOpen(port int) {
	c.record(port, "open")
}

// RecordClose records a port-close event.
func (c *ChurnTracker) RecordClose(port int) {
	c.record(port, "close")
}

func (c *ChurnTracker) record(port int, state string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.evict(port, now)
	c.events[port] = append(c.events[port], churnEvent{at: now, state: state})
}

// evict removes events outside the sliding window. Caller must hold mu.
func (c *ChurnTracker) evict(port int, now time.Time) {
	cutoff := now.Add(-c.window)
	evs := c.events[port]
	var keep []churnEvent
	for _, e := range evs {
		if !e.at.Before(cutoff) {
			keep = append(keep, e)
		}
	}
	if len(keep) == 0 {
		delete(c.events, port)
	} else {
		c.events[port] = keep
	}
}

// Unstable returns ChurnRecords for all ports whose churn rate meets or
// exceeds the configured threshold within the current window.
func (c *ChurnTracker) Unstable() []ChurnRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	var out []ChurnRecord
	for port, evs := range c.events {
		c.evict(port, now)
		evs = c.events[port]
		if len(evs) == 0 {
			continue
		}
		rec := ChurnRecord{Port: port, LastSeen: evs[len(evs)-1].at}
		for _, e := range evs {
			if e.state == "open" {
				rec.Opens++
			} else {
				rec.Closes++
			}
		}
		if rec.ChurnRate() >= c.threshold {
			out = append(out, rec)
		}
	}
	return out
}

// Reset clears all tracked events for a specific port.
func (c *ChurnTracker) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.events, port)
}
