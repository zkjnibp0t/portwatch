package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// CorrelationHook records port diffs into a Correlator and logs groups
// that exceed a minimum size, indicating coordinated port activity.
type CorrelationHook struct {
	correlator *ports.Correlator
	minGroup   int
	logger     *log.Logger
}

func NewCorrelationHook(window time.Duration, minGroup int, w io.Writer) *CorrelationHook {
	if w == nil {
		w = os.Stdout
	}
	return &CorrelationHook{
		correlator: ports.NewCorrelator(window),
		minGroup:   minGroup,
		logger:     log.New(w, "[correlation] ", 0),
	}
}

func (h *CorrelationHook) BeforeScan() error { return nil }

func (h *CorrelationHook) AfterScan(diff ports.Diff) error {
	for _, p := range diff.Opened {
		h.correlator.Record(p, true)
	}
	for _, p := range diff.Closed {
		h.correlator.Record(p, false)
	}
	for _, g := range h.correlator.Groups() {
		if len(g.Ports) >= h.minGroup {
			h.logger.Println(fmt.Sprintf("correlated group detected: %s", g))
		}
	}
	return nil
}
