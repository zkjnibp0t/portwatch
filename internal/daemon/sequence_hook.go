package daemon

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// SequenceHook assigns and logs a sequence number for every port-change event
// so that downstream consumers can detect missed or out-of-order alerts.
type SequenceHook struct {
	tracker *ports.SequenceTracker
	logger  *log.Logger
}

// NewSequenceHook returns a SequenceHook backed by tracker, writing to w.
// If w is nil, os.Stdout is used.
func NewSequenceHook(tracker *ports.SequenceTracker, w io.Writer) *SequenceHook {
	if w == nil {
		w = os.Stdout
	}
	return &SequenceHook{
		tracker: tracker,
		logger:  log.New(w, "[sequence] ", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *SequenceHook) BeforeScan() {}

// AfterScan assigns sequence numbers to every opened and closed port and logs
// each assignment.
func (h *SequenceHook) AfterScan(diff ports.Diff) {
	for _, p := range diff.Opened {
		seq := h.tracker.Next(p)
		h.logger.Println(fmt.Sprintf("opened port=%d seq=%d", p, seq))
	}
	for _, p := range diff.Closed {
		seq := h.tracker.Next(p)
		h.logger.Println(fmt.Sprintf("closed port=%d seq=%d", p, seq))
	}
}
