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
func NewLabelHook(labeler *ports.Labeler, out io.Writer) *LabelHook {
	if out == nil {
		out = os.Stdout
	}
	return &LabelHook{
		labeler: labeler,
		logger:  log.New(out, "[label] ", 0),
	}
}

// OnCycle is called after each scan cycle with the diff result.
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

func (h *LabelHook) labelList(portNums []int) string {
	parts := make([]string, 0, len(portNums))
	for _, p := range portNums {
		label := h.labeler.Label(p)
		parts = append(parts, fmt.Sprintf("%d(%s)", p, label))
	}
	return strings.Join(parts, ", ")
}
