package ports

import (
	"sync"
	"time"
)

// ScanThrottle prevents redundant scans when the port set has not changed
// within a minimum quiet period.
type ScanThrottle struct {
	mu          sync.Mutex
	minInterval time.Duration
	lastScan    time.Time
	lastHash    string
	clock       func() time.Time
}

// NewScanThrottle creates a ScanThrottle with the given minimum interval.
func NewScanThrottle(minInterval time.Duration) *ScanThrottle {
	return &ScanThrottle{
		minInterval: minInterval,
		clock:       time.Now,
	}
}

// ShouldSkip returns true when the fingerprint matches the last scan and the
// minimum interval has not yet elapsed.
func (t *ScanThrottle) ShouldSkip(fingerprint string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	if fingerprint == t.lastHash && now.Sub(t.lastScan) < t.minInterval {
		return true
	}
	return false
}

// Record updates the last scan time and fingerprint.
func (t *ScanThrottle) Record(fingerprint string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastScan = t.clock()
	t.lastHash = fingerprint
}

// Reset clears the stored state, allowing the next scan unconditionally.
func (t *ScanThrottle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastScan = time.Time{}
	t.lastHash = ""
}

// LastScan returns the time of the most recent recorded scan.
func (t *ScanThrottle) LastScan() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScan
}
