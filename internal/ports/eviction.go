package ports

import (
	"sync"
	"time"
)

// EvictionPolicy controls how stale port observations are removed.
type EvictionPolicy struct {
	mu      sync.Mutex
	ttl     time.Duration
	seen    map[int]time.Time
	clock   func() time.Time
}

// NewEvictionPolicy creates an EvictionPolicy that evicts ports not seen within ttl.
func NewEvictionPolicy(ttl time.Duration) *EvictionPolicy {
	return &EvictionPolicy{
		ttl:   ttl,
		seen:  make(map[int]time.Time),
		clock: time.Now,
	}
}

// Touch marks a set of ports as observed at the current time.
func (e *EvictionPolicy) Touch(ports map[int]struct{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.clock()
	for p := range ports {
		e.seen[p] = now
	}
}

// Evict returns ports that have not been touched within the TTL and removes them.
func (e *EvictionPolicy) Evict() []int {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.clock()
	var evicted []int
	for p, last := range e.seen {
		if now.Sub(last) > e.ttl {
			evicted = append(evicted, p)
			delete(e.seen, p)
		}
	}
	return evicted
}

// Len returns the number of currently tracked ports.
func (e *EvictionPolicy) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.seen)
}

// Reset clears all tracked ports.
func (e *EvictionPolicy) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.seen = make(map[int]time.Time)
}
