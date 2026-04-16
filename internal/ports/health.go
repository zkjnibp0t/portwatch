package ports

import "time"

// HealthStatus represents the current health of the port scanner.
type HealthStatus struct {
	LastScanAt   time.Time
	LastScanOK   bool
	LastError    string
	ConsecErrors int
	TotalScans   int
	TotalErrors  int
}

// HealthTracker tracks scanner health across cycles.
type HealthTracker struct {
	status HealthStatus
	clock  func() time.Time
}

// NewHealthTracker creates a HealthTracker with an optional clock override.
func NewHealthTracker(clock func() time.Time) *HealthTracker {
	if clock == nil {
		clock = time.Now
	}
	return &HealthTracker{clock: clock}
}

// RecordSuccess records a successful scan cycle.
func (h *HealthTracker) RecordSuccess() {
	h.status.LastScanAt = h.clock()
	h.status.LastScanOK = true
	h.status.LastError = ""
	h.status.ConsecErrors = 0
	h.status.TotalScans++
}

// RecordError records a failed scan cycle.
func (h *HealthTracker) RecordError(err error) {
	h.status.LastScanAt = h.clock()
	h.status.LastScanOK = false
	h.status.LastError = err.Error()
	h.status.ConsecErrors++
	h.status.TotalScans++
	h.status.TotalErrors++
}

// Status returns a copy of the current health status.
func (h *HealthTracker) Status() HealthStatus {
	return h.status
}

// IsHealthy returns true if the last scan was successful.
func (h *HealthTracker) IsHealthy() bool {
	return h.status.LastScanOK
}
