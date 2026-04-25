package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// LifecycleHook records port open/close transitions and logs notable durations.
type LifecycleHook struct {
	tracker  *ports.LifecycleTracker
	minLog   time.Duration // log closed ports open longer than this
	logger   *log.Logger
	lastDiff ports.Diff
}

// NewLifecycleHook creates a LifecycleHook with the given minimum log duration.
func NewLifecycleHook(tracker *ports.LifecycleTracker, minLog time.Duration, w io.Writer) *LifecycleHook {
	if w == nil {
		w = os.Stdout
	}
	return &LifecycleHook{
		tracker: tracker,
		minLog:  minLog,
		logger:  log.New(w, "[lifecycle] ", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *LifecycleHook) BeforeScan() error { return nil }

// AfterScan records opened and closed ports, logging long-lived closures.
func (h *LifecycleHook) AfterScan(diff ports.Diff) error {
	h.lastDiff = diff

	for _, p := range diff.Opened {
		h.tracker.RecordOpened(p)
	}

	for _, p := range diff.Closed {
		h.tracker.RecordClosed(p)
	}

	if h.minLog > 0 {
		for _, ev := range h.tracker.Events() {
			if ev.NextState == "closed" && ev.Duration >= h.minLog {
				h.logger.Println(fmt.Sprintf(
					"port %d was open for %s",
					ev.Port, formatLifecycleDur(ev.Duration),
				))
			}
		}
	}
	return nil
}

func formatLifecycleDur(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}
