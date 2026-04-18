package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/user/portwatch/internal/ports"
)

// LabelHook logs human-readable labels for opened/closed ports after each scan cycle.
type LabelHook struct {
	labeler *ports.Labeler
	logger  *log.Logger
	out     io.Writer
}

// NewLabelHook creates a LabelHook using the provided Labeler.
// If out is nil, os.Stdout is used.
func NewLabelHook(labeler *ports.Labeler, out io.Writer) *LabelHook {
	if out == nil {
		out = os.Stdout
	}
	return &LabelHook{
		labeler: labeler,
		logger:  log.New(out, "[label] ", 0),
		out:     out,
	}
}

// OnCycle is called after each scan cycle with the diff result.
// It logs opened and closed ports with their human-readable labels.
// If there are no changes, it returns early without logging.
func (h *LabelHook) OnCycle(diff ports.Diff) {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return
	}
	if len(diff.Opened) > 0 {
		h.logger.Printf("opened: %s", h.labelList(diff.Opened))
	}
	if len(diff.Closed) > 0 {
		h.logger.Printf("closed: %s", h.labelList(diff.Closed))
	}
}

// SetPrefix updates the log prefix used by the hook's logger.
func (h *LabelHook) SetPrefix(prefix string) {
	h.logger.SetPrefix(prefix)
}

func (h *LabelHook) labelList(portNums []int) string {
	parts := make([]string, 0, len(portNums))
	for _, p := range portNums {
		label := h.labeler.Label(p)
		parts = append(parts, fmt.Sprintf("%d(%s)", p, label))
	}
	return strings.Join(parts, ", ")
}
