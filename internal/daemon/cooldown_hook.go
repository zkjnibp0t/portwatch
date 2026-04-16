package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/ports/diff"
)

// CooldownHook filters diff events suppressing ports still within cooldown.
type CooldownHook struct {
	tracker *ports.CooldownTracker
	logger  *log.Logger
}

// NewCooldownHook creates a CooldownHook with the given window duration.
func NewCooldownHook(window time.Duration, logger *log.Logger) *CooldownHook {
	if logger == nil {
		logger = log.Default()
	}
	return &CooldownHook{
		tracker: ports.NewCooldownTracker(window),
		logger:  logger,
	}
}

// Filter removes ports from the diff that are within their cooldown window
// and records newly seen ports.
func (h *CooldownHook) Filter(d diff.Diff) diff.Diff {
	opened := filterCooldown(h.tracker, h.logger, d.Opened, "opened")
	closed := filterCooldown(h.tracker, h.logger, d.Closed, "closed")
	return diff.Diff{Opened: opened, Closed: closed}
}

func filterCooldown(tracker *ports.CooldownTracker, logger *log.Logger, ports []int, direction string) []int {
	var result []int
	for _, p := range ports {
		if tracker.IsActive(p) {
			logger.Printf("cooldown_hook: suppressing %s port %d (still in cooldown)", direction, p)
			continue
		}
		tracker.Record(p)
		result = append(result, p)
	}
	return result
}
