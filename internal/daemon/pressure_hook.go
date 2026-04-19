package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// PressureHook records opened-port counts into a PressureTracker and logs
// a warning whenever the system crosses the pressure threshold.
type PressureHook struct {
	tracker *ports.PressureTracker
	logger  *log.Logger
}

func NewPressureHook(window time.Duration, threshold int, w io.Writer) *PressureHook {
	if w == nil {
		w = os.Stdout
	}
	return &PressureHook{
		tracker: ports.NewPressureTracker(window, threshold),
		logger:  log.New(w, "", 0),
	}
}

func (h *PressureHook) BeforeScan() error { return nil }

func (h *PressureHook) AfterScan(diff ports.Diff) error {
	opened := len(diff.Opened)
	if opened == 0 {
		return nil
	}
	h.tracker.Record(opened)
	if h.tracker.UnderPressure() {
		h.logger.Println(fmt.Sprintf(
			"[pressure] high port-open rate detected: %d events in window",
			h.tracker.Count(),
		))
	}
	return nil
}
