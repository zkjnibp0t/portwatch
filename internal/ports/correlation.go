package ports

import (
	"fmt"
	"sync"
	"time"
)

// CorrelationEvent records a port open/close with a timestamp.
type CorrelationEvent struct {
	Port      int
	Opened    bool
	Timestamp time.Time
}

// CorrelationGroup is of ports that changed within a time window.
type CorrelationGroup struct {
	Ports  []int
	Opened bool
	At     time.Time
}

func (g CorrelationGroup) String() string {
	dir := "opened"
	if !g.Opened {
		dir = "closed"
	}
	return fmt.Sprintf("%s %v at %s", dir, g.Ports, g.At.Format(time.RFC3339))
}

// Correlator groups port events that occur within a correlation window.
type Correlator struct {
	mu     sync.Mutex
	window time.Duration
	events []CorrelationEvent
	clock  func() time.Time
}

func NewCorrelator(window time.Duration) *Correlator {
	return &Correlator{window: window, clock: time.Now}
}

// Record adds a port event.
func (c *Correlator) Record(port int, opened bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, CorrelationEvent{Port: port, Opened: opened, Timestamp: c.clock()})
}

// Groups returns correlated port groups within the current window.
func (c *Correlator) Groups() []CorrelationGroup {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := c.clock().Add(-c.window)
	buckets := map[bool][]int{}
	var latest time.Time
	for _, e := range c.events {
		if e.Timestamp.Before(cutoff) {
			continue
		}
		buckets[e.Opened] = append(buckets[e.Opened], e.Port)
		if e.Timestamp.After(latest) {
			latest = e.Timestamp
		}
	}
	var groups []CorrelationGroup
	for opened, ports := range buckets {
		if len(ports) > 0 {
			groups = append(groups, CorrelationGroup{Ports: ports, Opened: opened, At: latest})
		}
	}
	return groups
}

// Flush clears all recorded events.
func (c *Correlator) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = nil
}
