package daemon

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// HealthHook integrates a HealthTracker into the daemon cycle.
type HealthHook struct {
	tracker       *ports.HealthTracker
	errorThreshold int
	w              io.Writer
}

// NewHealthHook creates a HealthHook that warns when consecutive errors exceed threshold.
func NewHealthHook(tracker *ports.HealthTracker, errorThreshold int, w io.Writer) *HealthHook {
	if w == nil {
		w = os.Stderr
	}
	if errorThreshold <= 0 {
		errorThreshold = 3
	}
	return &HealthHook{tracker: tracker, errorThreshold: errorThreshold, w: w}
}

// OnSuccess records a successful scan.
func (h *HealthHook) OnSuccess() {
	h.tracker.RecordSuccess()
}

// OnError records a scan error and emits a warning when threshold is exceeded.
func (h *HealthHook) OnError(err error) {
	h.tracker.RecordError(err)
	s := h.tracker.Status()
	if s.ConsecErrors >= h.errorThreshold {
		fmt.Fprintf(h.w, "[portwatch] WARNING: %d consecutive scan errors (last: %s)\n",
			s.ConsecErrors, s.LastError)
	}
}

// Summary prints a brief health summary to the writer.
func (h *HealthHook) Summary() {
	s := h.tracker.Status()
	status := "OK"
	if !s.LastScanOK {
		status = "ERROR"
	}
	fmt.Fprintf(h.w, "[portwatch] health: status=%s scans=%d errors=%d consec=%d\n",
		status, s.TotalScans, s.TotalErrors, s.ConsecErrors)
}
