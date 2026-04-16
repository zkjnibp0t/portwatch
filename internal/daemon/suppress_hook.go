package daemon

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SuppressHook integrates a Suppressor into the daemon cycle,
// filtering diff events and logging suppressed ports.
type SuppressHook struct {
	suppressor *ports.Suppressor
	writer     io.Writer
}

// NewSuppressHook creates a SuppressHook with the given suppressor.
func NewSuppressHook(sup *ports.Suppressor, w io.Writer) *SuppressHook {
	if w == nil {
		w = os.Stdout
	}
	return &SuppressHook{suppressor: sup, writer: w}
}

// FilterOpened removes suppressed ports from the opened list and logs them.
func (h *SuppressHook) FilterOpened(opened []int, process func(int) string) []int {
	return h.filter(opened, process, "opened")
}

// FilterClosed removes suppressed ports from the closed list and logs them.
func (h *SuppressHook) FilterClosed(closed []int, process func(int) string) []int {
	return h.filter(closed, process, "closed")
}

func (h *SuppressHook) filter(ports_ []int, process func(int) string, dir string) []int {
	var kept []int
	for _, p := range ports_ {
		proc := ""
		if process != nil {
			proc = process(p)
		}
		if h.suppressor.IsSuppressed(p, proc) {
			fmt.Fprintf(h.writer, "[suppress] port %d (%s) suppressed for %s\n", p, proc, dir)
			continue
		}
		kept = append(kept, p)
	}
	return kept
}

// SuppressFor is a convenience wrapper to add a suppression rule.
func (h *SuppressHook) SuppressFor(port int, process string, d time.Duration) {
	h.suppressor.Suppress(port, process, d)
	fmt.Fprintf(h.writer, "[suppress] port %d suppressed for %s\n", port, d)
}
