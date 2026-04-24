package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// ReopenHook integrates ReopenDetector into the daemon scan cycle.
// It logs a warning whenever a port is reopened within the observation window.
type ReopenHook struct {
	detector *ports.ReopenDetector
	logger   *log.Logger
}

// NewReopenHook creates a ReopenHook with the given observation window.
func NewReopenHook(window time.Duration, w io.Writer) *ReopenHook {
	if w == nil {
		w = os.Stdout
	}
	return &ReopenHook{
		detector: ports.NewReopenDetector(window),
		logger:   log.New(w, "", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *ReopenHook) BeforeScan() error { return nil }

// AfterScan processes the diff to record closed ports and detect reopens.
func (h *ReopenHook) AfterScan(diff ports.Diff) error {
	// First record all newly closed ports.
	for _, port := range diff.Closed {
		h.detector.RecordClosed(port)
	}

	// Then check whether any opened ports were recently closed.
	for _, port := range diff.Opened {
		if event := h.detector.RecordOpened(port); event != nil {
			h.logger.Printf(
				"[reopen] port %d reopened after %v (closed at %s)",
				event.Port,
				formatGap(event.Gap),
				event.ClosedAt.Format(time.RFC3339),
			)
		}
	}

	h.detector.Purge()
	return nil
}

func formatGap(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return d.Round(time.Millisecond).String()
}
