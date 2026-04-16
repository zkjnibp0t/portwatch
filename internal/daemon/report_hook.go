package daemon

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/ports"
)

// ReportHook writes a formatted change report to a writer after each scan cycle.
type ReportHook struct {
	store  *history.Store
	logger *log.Logger
	out    io.Writer
}

// NewReportHook creates a ReportHook backed by the given history store.
func NewReportHook(store *history.Store, out io.Writer) *ReportHook {
	if out == nil {
		out = os.Stdout
	}
	return &ReportHook{
		store:  store,
		logger: log.New(out, "[report] ", log.LstdFlags),
		out:    out,
	}
}

// OnCycle is called after each scan cycle with the diff result.
func (r *ReportHook) OnCycle(diff ports.Diff) {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return
	}
	if len(diff.Opened) > 0 {
		r.logger.Printf("opened ports: %s", formatPorts(diff.Opened))
	}
	if len(diff.Closed) > 0 {
		r.logger.Printf("closed ports: %s", formatPorts(diff.Closed))
	}
	records := r.store.Recent(5)
	if len(records) > 0 {
		r.logger.Printf("recent changes: %d recorded events", len(records))
	}
}

func formatPorts(ps []int) string {
	if len(ps) == 0 {
		return "-"
	}
	out := ""
	for i, p := range ps {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%d", p)
	}
	return out
}
