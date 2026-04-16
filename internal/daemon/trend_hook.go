package daemon

import (
	"log"

	"github.com/user/portwatch/internal/ports"
)

// TrendHook wires a TrendTracker into the daemon scan cycle.
// After each diff it records opened/closed events and logs
// any port that has flapped more than the configured threshold.
type TrendHook struct {
	tracker   *ports.TrendTracker
	threshold int
	logger    *log.Logger
}

// NewTrendHook creates a TrendHook with the supplied tracker and flap threshold.
func NewTrendHook(tracker *ports.TrendTracker, threshold int, logger *log.Logger) *TrendHook {
	if threshold <= 0 {
		threshold = 5
	}
	return &TrendHook{tracker: tracker, threshold: threshold, logger: logger}
}

// Apply records diff events and warns on high-flap ports.
func (h *TrendHook) Apply(diff ports.Diff) {
	for _, p := range diff.Opened {
		h.tracker.Record(p, "opened")
		h.checkFlap(p)
	}
	for _, p := range diff.Closed {
		h.tracker.Record(p, "closed")
		h.checkFlap(p)
	}
}

func (h *TrendHook) checkFlap(port int) {
	count := h.tracker.FlapCount(port)
	if count >= h.threshold {
		h.logger.Printf("[portwatch] WARNING: port %d has flapped %d time(s) recently", port, count)
	}
}
