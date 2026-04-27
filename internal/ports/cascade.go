package ports

import "sync"

// CascadeEvent represents a port that triggered a cascade — a rapid succession
// of related port openings within a short time window.
type CascadeEvent struct {
	Trigger  int
	FollowOn []int
}

// CascadeDetector detects when the opening of one port is quickly followed by
// openings of other ports, suggesting a coordinated or automated process.
type CascadeDetector struct {
	mu       sync.Mutex
	window   int64 // nanoseconds
	events   []cascadeEntry
	minGroup int
	clock    func() int64
}

type cascadeEntry struct {
	port int
	at   int64
}

// NewCascadeDetector creates a CascadeDetector. window is the time window in
// nanoseconds within which ports are considered part of the same cascade.
// minGroup is the minimum number of ports (including trigger) to qualify.
func NewCascadeDetector(window int64, minGroup int, clock func() int64) *CascadeDetector {
	if minGroup < 2 {
		minGroup = 2
	}
	return &CascadeDetector{
		window:   window,
		minGroup: minGroup,
		clock:    clock,
	}
}

// Record registers newly opened ports and returns any detected cascades.
func (c *CascadeDetector) Record(opened []int) []CascadeEvent {
	if len(opened) == 0 {
		return nil
	}
	now := c.clock()
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict old entries outside the window.
	cutoff := now - c.window
	filtered := c.events[:0]
	for _, e := range c.events {
		if e.at >= cutoff {
			filtered = append(filtered, e)
		}
	}
	c.events = filtered

	// Collect ports already in window before this batch.
	prior := make([]int, len(c.events))
	for i, e := range c.events {
		prior[i] = e.port
	}

	// Add newly opened ports.
	for _, p := range opened {
		c.events = append(c.events, cascadeEntry{port: p, at: now})
	}

	// A cascade exists when prior + opened together meet minGroup.
	total := append(prior, opened...)
	if len(total) < c.minGroup {
		return nil
	}

	// The first port in the window is the trigger.
	trigger := total[0]
	followOn := make([]int, 0, len(total)-1)
	for _, p := range total[1:] {
		followOn = append(followOn, p)
	}
	return []CascadeEvent{{Trigger: trigger, FollowOn: followOn}}
}

// Reset clears all recorded events.
func (c *CascadeDetector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}
