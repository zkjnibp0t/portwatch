package ports

import "time"

// DriftEvent records a port that has been open longer than expected.
type DriftEvent struct {
	Port      int
	OpenSince time.Time
	Drift     time.Duration
}

// DriftDetector flags ports whose open duration exceeds a configurable threshold.
type DriftDetector struct {
	threshold time.Duration
	clock     func() time.Time
	openSince map[int]time.Time
}

// NewDriftDetector creates a DriftDetector with the given threshold.
func NewDriftDetector(threshold time.Duration, clock func() time.Time) *DriftDetector {
	if clock == nil {
		clock = time.Now
	}
	return &DriftDetector{
		threshold: threshold,
		clock:     clock,
		openSince: make(map[int]time.Time),
	}
}

// Observe records opened and closed ports for drift tracking.
func (d *DriftDetector) Observe(opened, closed []int) {
	now := d.clock()
	for _, p := range opened {
		if _, exists := d.openSince[p]; !exists {
			d.openSince[p] = now
		}
	}
	for _, p := range closed {
		delete(d.openSince, p)
	}
}

// Detect returns all ports that have been open longer than the threshold.
func (d *DriftDetector) Detect() []DriftEvent {
	now := d.clock()
	var events []DriftEvent
	for port, since := range d.openSince {
		drift := now.Sub(since)
		if drift >= d.threshold {
			events = append(events, DriftEvent{
				Port:      port,
				OpenSince: since,
				Drift:     drift,
			})
		}
	}
	return events
}

// Reset clears all tracked ports.
func (d *DriftDetector) Reset() {
	d.openSince = make(map[int]time.Time)
}
