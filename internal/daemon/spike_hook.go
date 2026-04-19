package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SpikeHook watches for ports that open and close within a short window.
type SpikeHook struct {
	detector *ports.SpikeDetector
	logger   *log.Logger
}

// NewSpikeHook creates a SpikeHook with the given spike window duration.
func NewSpikeHook(window time.Duration, w io.Writer) *SpikeHook {
	if w == nil {
		w = os.Stdout
	}
	return &SpikeHook{
		detector: ports.NewSpikeDetector(window),
		logger:   log.New(w, "[spike] ", 0),
	}
}

// BeforeScan is a no-op.
func (h *SpikeHook) BeforeScan() {}

// AfterScan processes diff events to detect spike ports.
func (h *SpikeHook) AfterScan(diff ports.Diff) {
	for _, p := range diff.Opened {
		h.detector.RecordOpened(p)
	}
	for _, p := range diff.Closed {
		if ev, ok := h.detector.RecordClosed(p); ok {
			h.logger.Println(fmt.Sprintf(
				"spike detected on port %d — open for %v",
				ev.Port, ev.Duration.Round(time.Millisecond),
			))
		}
	}
}
