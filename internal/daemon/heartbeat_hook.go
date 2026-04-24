package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/user/portwatch/internal/ports"
)

// HeartbeatHook records a heartbeat for every open port after each scan and
// logs ports that have gone silent (open but not seen recently).
type HeartbeatHook struct {
	tracker *ports.HeartbeatTracker
	log     *log.Logger
}

// NewHeartbeatHook returns a hook backed by tracker.
func NewHeartbeatHook(tracker *ports.HeartbeatTracker, w io.Writer) *HeartbeatHook {
	if w == nil {
		w = os.Stdout
	}
	return &HeartbeatHook{
		tracker: tracker,
		log:     log.New(w, "[heartbeat] ", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *HeartbeatHook) BeforeScan() {}

// AfterScan beats every port in the current open set and logs any that have
// gone silent since the last scan cycle.
func (h *HeartbeatHook) AfterScan(current ports.Set, diff ports.Diff) {
	// Update heartbeats for all currently open ports.
	for port := range current {
		h.tracker.Beat(port)
	}

	// Remove ports that closed cleanly so they don't linger as silent.
	for _, port := range diff.Closed {
		h.tracker.Remove(port)
	}

	// Report any ports that are tracked but haven't been seen within the deadline.
	silent := h.tracker.Silent()
	if len(silent) == 0 {
		return
	}
	sort.Ints(silent)
	for _, p := range silent {
		h.log.Println(fmt.Sprintf("port %d silent — no heartbeat within deadline", p))
	}
}
