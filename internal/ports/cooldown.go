package ports

import (
	"sync"
	"time"
)

// CooldownTracker prevents repeated alerts for the same port within a window.
type CooldownTracker struct {
	mu       sync.Mutex
	window   time.Duration
	entries  map[int]time.Time
	clock    func() time.Time
}

// NewCooldownTracker creates a CooldownTracker with the given cooldown window.
func NewCooldownTracker(window time.Duration) *CooldownTracker {
	return &CooldownTracker{
		window:  window,
		entries: make(map[int]time.Time),
		clock:   time.Now,
	}
}

// IsActive returns true if the port is still within its cooldown window.
func (c *CooldownTracker) IsActive(port int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if t, ok := c.entries[port]; ok {
		if c.clock().Before(t.Add(c.window)) {
			return true
		}
	}
	return false
}

// Record marks the port as having just triggered, starting its cooldown.
func (c *CooldownTracker) Record(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[port] = c.clock()
}

// Reset clears the cooldown for a port immediately.
func (c *CooldownTracker) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// Prune removes expired entries to keep memory bounded.
func (c *CooldownTracker) Prune() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	for port, t := range c.entries {
		if now.After(t.Add(c.window)) {
			delete(c.entries, port)
		}
	}
}
