package daemon

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// RollupHook accumulates diffs and periodically logs a port-activity summary.
type RollupHook struct {
	rollup  *ports.Rollup
	every   int // log every N cycles
	cycle   int
	logger  *log.Logger
}

// NewRollupHook creates a RollupHook that prints a summary every `every` cycles.
func NewRollupHook(every int, w io.Writer) *RollupHook {
	if w == nil {
		w = os.Stdout
	}
	if every < 1 {
		every = 1
	}
	return &RollupHook{
		rollup:  ports.NewRollup(),
		every:   every,
		logger:  log.New(w, "[rollup] ", 0),
	}
}

// BeforeScan is a no-op.
func (h *RollupHook) BeforeScan() {}

// AfterScan accumulates the diff and logs the summary every N cycles.
func (h *RollupHook) AfterScan(diff ports.Diff) {
	h.rollup.Record(diff)
	h.cycle++
	if h.cycle < h.every {
		return
	}
	h.cycle = 0
	entries := h.rollup.Entries()
	if len(entries) == 0 {
		h.logger.Println("no port activity in window")
	} else {
		for _, e := range entries {
			h.logger.Println(fmt.Sprintf("port=%d opened=%d closed=%d net=%+d",
				e.Port, e.Opened, e.Closed, e.Net))
		}
	}
	h.rollup.Reset()
}
