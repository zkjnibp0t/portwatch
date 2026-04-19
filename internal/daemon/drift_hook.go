package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// DriftHook integrates DriftDetector into the daemon scan cycle.
type DriftHook struct {
	detector  *ports.DriftDetector
	threshold time.Duration
	logger    *log.Logger
}

// NewDriftHook creates a DriftHook with the given threshold and optional writer.
func NewDriftHook(threshold time.Duration, w io.Writer) *DriftHook {
	if w == nil {
		w = os.Stdout
	}
	return &DriftHook{
		detector:  ports.NewDriftDetector(threshold, nil),
		threshold: threshold,
		logger:    log.New(w, "[drift] ", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *DriftHook) BeforeScan() {}

// AfterScan observes the diff and logs any ports drifting past the threshold.
func (h *DriftHook) AfterScan(diff ports.Diff) {
	h.detector.Observe(diff.Opened, diff.Closed)
	events := h.detector.Detect()
	for _, e := range events {
		h.logger.Println(fmt.Sprintf(
			"port %d has been open for %s (threshold %s)",
			e.Port,
			e.Drift.Round(time.Second),
			h.threshold,
		))
	}
}
