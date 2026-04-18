package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// ShadowHook integrates ShadowTracker into the daemon cycle.
type ShadowHook struct {
	tracker *ports.ShadowTracker
	log     *log.Logger
	out     io.Writer
}

// NewShadowHook creates a ShadowHook with the given flash window.
func NewShadowHook(window time.Duration, logger *log.Logger) *ShadowHook {
	if logger == nil {
		logger = log.New(os.Stdout, "[shadow] ", 0)
	}
	return &ShadowHook{
		tracker: ports.NewShadowTracker(window),
		log:     logger,
		out:     os.Stdout,
	}
}

// BeforeScan is a no-op for this hook.
func (h *ShadowHook) BeforeScan() {}

// AfterScan processes diff events to detect shadow ports.
func (h *ShadowHook) AfterScan(diff ports.Diff) {
	for _, p := range diff.Opened {
		h.tracker.RecordOpened(p)
	}
	for _, p := range diff.Closed {
		if e, ok := h.tracker.RecordClosed(p); ok {
			h.log.Print(fmt.Sprintf(
				"shadow port detected: %d was open for %v",
				e.Port, e.Duration.Round(time.Millisecond),
			))
		}
	}
}

// Shadows exposes accumulated shadow entries for reporting.
func (h *ShadowHook) Shadows() []ports.ShadowEntry {
	return h.tracker.Shadows()
}
