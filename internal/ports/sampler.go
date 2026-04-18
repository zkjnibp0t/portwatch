package ports

import "time"

// Sample holds a port set snapshot with a timestamp.
type Sample struct {
	Timestamp time.Time
	Ports     map[int]struct{}
}

// Sampler keeps a rolling window of recent port-set samples.
type Sampler struct {
	window  time.Duration
	samples []Sample
	clock   func() time.Time
}

// NewSampler creates a Sampler that retains samples within the given window.
func NewSampler(window time.Duration) *Sampler {
	return &Sampler{window: window, clock: time.Now}
}

// Add records a new sample, pruning entries older than the window.
func (s *Sampler) Add(ports map[int]struct{}) {
	now := s.clock()
	s.samples = append(s.samples, Sample{Timestamp: now, Ports: ports})
	s.prune(now)
}

// Latest returns the most recent sample, or nil if none.
func (s *Sampler) Latest() *Sample {
	if len(s.samples) == 0 {
		return nil
	}
	return &s.samples[len(s.samples)-1]
}

// All returns all retained samples.
func (s *Sampler) All() []Sample {
	return s.samples
}

// Len returns the number of retained samples.
func (s *Sampler) Len() int {
	return len(s.samples)
}

func (s *Sampler) prune(now time.Time) {
	cutoff := now.Add(-s.window)
	i := 0
	for i < len(s.samples) && s.samples[i].Timestamp.Before(cutoff) {
		i++
	}
	s.samples = s.samples[i:]
}
