package ports

import (
	"sync"
	"time"
)

// HeartbeatTracker records the last time a port was seen open and reports
// ports that have gone silent longer than a configurable deadline.
type HeartbeatTracker struct {
	mu       sync.Mutex
	hearts   map[int]time.Time
	deadline time.Duration
	clock    func() time.Time
}

// NewHeartbeatTracker returns a tracker that flags ports unseen for longer
// than deadline.
func NewHeartbeatTracker(deadline time.Duration) *HeartbeatTracker {
	return &HeartbeatTracker{
		hearts:   make(map[int]time.Time),
		deadline: deadline,
		clock:    time.Now,
	}
}

// Beat records that port was observed open at the current time.
func (h *HeartbeatTracker) Beat(port int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hearts[port] = h.clock()
}

// Remove stops tracking port entirely (e.g. after it closes cleanly).
func (h *HeartbeatTracker) Remove(port int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.hearts, port)
}

// Silent returns all ports whose last heartbeat is older than the deadline.
func (h *HeartbeatTracker) Silent() []int {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := h.clock()
	var out []int
	for port, last := range h.hearts {
		if now.Sub(last) > h.deadline {
			out = append(out, port)
		}
	}
	return out
}

// LastSeen returns the last heartbeat time for port and whether it is tracked.
func (h *HeartbeatTracker) LastSeen(port int) (time.Time, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	t, ok := h.hearts[port]
	return t, ok
}
