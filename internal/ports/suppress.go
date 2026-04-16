package ports

import (
	"sync"
	"time"
)

// SuppressRule defines a port+process pair to suppress for a duration.
type SuppressRule struct {
	Port    int
	Process string // empty means any process
	Until   time.Time
}

// Suppressor temporarily silences alerts for specific ports.
type Suppressor struct {
	mu    sync.Mutex
	rules []SuppressRule
	now   func() time.Time
}

// NewSuppressor creates a Suppressor with real-time clock.
func NewSuppressor() *Suppressor {
	return &Suppressor{now: time.Now}
}

// Suppress adds a rule silencing the given port (and optionally process) for d.
func (s *Suppressor) Suppress(port int, process string, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = append(s.rules, SuppressRule{
		Port:    port,
		Process: process,
		Until:   s.now().Add(d),
	})
}

// IsSuppressed returns true if the port/process combo is currently suppressed.
func (s *Suppressor) IsSuppressed(port int, process string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	for _, r := range s.rules {
		if r.Port != port {
			continue
		}
		if now.After(r.Until) {
			continue
		}
		if r.Process == "" || r.Process == process {
			return true
		}
	}
	return false
}

// Prune removes expired rules.
func (s *Suppressor) Prune() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	active := s.rules[:0]
	for _, r := range s.rules {
		if now.Before(r.Until) {
			active = append(active, r)
		}
	}
	s.rules = active
}

// ActiveCount returns the number of non-expired rules.
func (s *Suppressor) ActiveCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	count := 0
	for _, r := range s.rules {
		if now.Before(r.Until) {
			count++
		}
	}
	return count
}
