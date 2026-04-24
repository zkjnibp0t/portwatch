package daemon

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// ChurnHook integrates the churn tracker into the daemon cycle.
// After each scan it records open/close events and logs any ports
// that have been flagged as unstable (churning) within the configured
// observation window.
type ChurnHook struct {
	tracker *ports.ChurnTracker
	logger  *log.Logger
}

// NewChurnHook returns a ChurnHook backed by the given tracker.
// If w is nil, output is written to os.Stdout.
func NewChurnHook(tracker *ports.ChurnTracker, w io.Writer) *ChurnHook {
	if w == nil {
		w = os.Stdout
	}
	return &ChurnHook{
		tracker: tracker,
		logger:  log.New(w, "[churn] ", 0),
	}
}

// BeforeScan is a no-op for the churn hook.
func (h *ChurnHook) BeforeScan() error { return nil }

// AfterScan records opened and closed ports into the churn tracker,
// then logs any ports that are currently considered unstable.
func (h *ChurnHook) AfterScan(event ScanEvent) error {
	for _, p := range event.Diff.Opened {
		h.tracker.RecordOpen(p)
	}
	for _, p := range event.Diff.Closed {
		h.tracker.RecordClose(p)
	}

	unstable := h.tracker.Unstable()
	if len(unstable) == 0 {
		return nil
	}

	for _, p := range unstable {
		summary := h.tracker.Summary(p)
		h.logger.Println(fmt.Sprintf(
			"port %d is churning: %d open / %d close events in window",
			p, summary.Opens, summary.Closes,
		))
	}
	return nil
}
