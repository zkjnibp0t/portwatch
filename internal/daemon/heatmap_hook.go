package daemon

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// HeatmapHook records per-port activity and periodically logs the top-N
// most active ports after every scan cycle.
type HeatmapHook struct {
	tracker  *ports.HeatmapTracker
	topN     int
	cycles   int
	every    int
	logger   *log.Logger
}

// NewHeatmapHook creates a HeatmapHook that logs the top-N ports every
// `every` cycles. If every <= 0 it defaults to 1.
func NewHeatmapHook(tracker *ports.HeatmapTracker, topN, every int, w io.Writer) *HeatmapHook {
	if every <= 0 {
		every = 1
	}
	if w == nil {
		w = os.Stdout
	}
	return &HeatmapHook{
		tracker: tracker,
		topN:    topN,
		every:   every,
		logger:  log.New(w, "[heatmap] ", 0),
	}
}

// BeforeScan is a no-op.
func (h *HeatmapHook) BeforeScan() {}

// AfterScan records the diff and logs top-N ports on the configured cadence.
func (h *HeatmapHook) AfterScan(diff ports.Diff) {
	if len(diff.Opened) > 0 || len(diff.Closed) > 0 {
		h.tracker.Record(diff)
	}
	h.cycles++
	if h.cycles%h.every != 0 {
		return
	}
	top := h.tracker.TopN(h.topN)
	if len(top) == 0 {
		return
	}
	h.logger.Printf("top %d active ports:", len(top))
	for i, e := range top {
		h.logger.Printf("  #%d port=%d hits=%d last=%s",
			i+1, e.Port, e.Hits,
			fmt.Sprintf("%s", e.LastSeen.Format("15:04:05")))
	}
}
