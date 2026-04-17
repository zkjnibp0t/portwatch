package ports

import (
	"fmt"
	"sync"
	"time"
)

// QuotaExceededError is returned when the scan quota is exceeded.
type QuotaExceededError struct {
	Limit    int
	Window   time.Duration
	Consumed int
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("scan quota exceeded: %d/%d scans in %s", e.Consumed, e.Limit, e.Window)
}

// QuotaTracker limits the number of scans within a rolling time window.
type QuotaTracker struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clock   func() time.Time
	buckets []time.Time
}

// NewQuotaTracker creates a tracker allowing at most limit scans per window.
func NewQuotaTracker(limit int, window time.Duration) *QuotaTracker {
	return &QuotaTracker{
		limit:  limit,
		window: window,
		clock:  time.Now,
	}
}

// Allow returns nil if a scan is permitted, or a QuotaExceededError.
func (q *QuotaTracker) Allow() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.clock()
	cutoff := now.Add(-q.window)

	// prune expired entries
	valid := q.buckets[:0]
	for _, t := range q.buckets {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	q.buckets = valid

	if len(q.buckets) >= q.limit {
		return &QuotaExceededError{Limit: q.limit, Window: q.window, Consumed: len(q.buckets)}
	}
	q.buckets = append(q.buckets, now)
	return nil
}

// Remaining returns how many scans are still allowed in the current window.
func (q *QuotaTracker) Remaining() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.clock()
	cutoff := now.Add(-q.window)
	count := 0
	for _, t := range q.buckets {
		if t.After(cutoff) {
			count++
		}
	}
	r := q.limit - count
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded scan timestamps.
func (q *QuotaTracker) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.buckets = nil
}
