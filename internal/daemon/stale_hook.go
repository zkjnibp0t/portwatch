package daemon

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// StaleHook observes open ports each cycle and logs any that exceed the
// staleness threshold.
type StaleHook struct {
	tracker   *ports.StaleTracker
	threshold time.Duration
	logger    io.Writer
}

// NewStaleHook creates a StaleHook with the given threshold.
func NewStaleHook(threshold time.Duration, logger io.Writer) *StaleHook {
	if logger == nil {
		logger = os.Stdout
	}
	return &StaleHook{
		tracker:   ports.NewStaleTracker(threshold),
		threshold: threshold,
		logger:    logger,
	}
}

// BeforeScan is a no-op for this hook.
func (h *StaleHook) BeforeScan() {}

// AfterScan updates the tracker with the current open ports and logs stale ones.
func (h *StaleHook) AfterScan(current ports.PortSet, diff ports.Diff) {
	// Remove ports that closed.
	for _, p := range diff.Closed {
		h.tracker.Remove(p)
	}
	// Observe all currently open ports.
	for p := range current {
		h.tracker.Observe(p)
	}
	// Report stale ports.
	stale := h.tracker.StalePorts()
	if len(stale) == 0 {
		return
	}
	for _, p := range stale {
		age := h.tracker.Age(p)
		fmt.Fprintf(h.logger, "[stale] port %d open for %s (threshold %s)\n",
			p, age.Round(time.Second), h.threshold)
	}
}
