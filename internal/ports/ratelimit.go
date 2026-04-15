package ports

import (
	"sync"
	"time"
)

// ScanRateLimiter enforces a minimum interval between successive scans
// to prevent CPU/network saturation when the daemon cycle is tight.
type ScanRateLimiter struct {
	mu       sync.Mutex
	minGap   time.Duration
	lastScan time.Time
	clock    func() time.Time
}

// NewScanRateLimiter creates a ScanRateLimiter that enforces at least minGap
// between scans. If minGap is zero or negative the limiter is effectively
// disabled (every call to Wait returns immediately).
func NewScanRateLimiter(minGap time.Duration) *ScanRateLimiter {
	return &ScanRateLimiter{
		minGap: minGap,
		clock:  time.Now,
	}
}

// Wait blocks until it is safe to start the next scan. It returns the time
// at which the scan was allowed to proceed.
func (r *ScanRateLimiter) Wait() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock()
	if r.minGap > 0 && !r.lastScan.IsZero() {
		next := r.lastScan.Add(r.minGap)
		if now.Before(next) {
			time.Sleep(next.Sub(now))
			now = r.clock()
		}
	}
	r.lastScan = now
	return now
}

// Reset clears the last-scan timestamp so the very next Wait call returns
// immediately regardless of the configured gap.
func (r *ScanRateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastScan = time.Time{}
}

// LastScan returns the timestamp of the most recent allowed scan, or the
// zero time if no scan has been allowed yet.
func (r *ScanRateLimiter) LastScan() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastScan
}
