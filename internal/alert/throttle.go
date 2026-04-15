package alert

import (
	"sync"
	"time"
)

// Throttler prevents duplicate alerts from firing too frequently.
// It tracks the last alert time per key and suppresses alerts
// that occur within the cooldown window.
type Throttler struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
}

// NewThrottler creates a Throttler with the given cooldown duration.
// Alerts with the same key will be suppressed if the previous alert
// was sent within the cooldown window.
func NewThrottler(cooldown time.Duration) *Throttler {
	return &Throttler{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

// Allow returns true if an alert for the given key should be sent.
// It updates the last-sent timestamp when returning true.
func (t *Throttler) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastSent[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSent[key] = now
	return true
}

// Reset clears the throttle state for a given key, allowing the
// next alert through immediately regardless of cooldown.
func (t *Throttler) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, key)
}

// ResetAll clears all throttle state.
func (t *Throttler) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSent = make(map[string]time.Time)
}
