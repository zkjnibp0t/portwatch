package daemon

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/portwatch/internal/ports"
)

// TagHook logs port changes annotated with their tag labels.
type TagHook struct {
	tagger *ports.Tagger
	w      io.Writer
}

// NewTagHook creates a TagHook with the given tagger and writer.
func NewTagHook(tagger *ports.Tagger, w io.Writer) *TagHook {
	if w == nil {
		w = os.Stdout
	}
	return &TagHook{tagger: tagger, w: w}
}

// OnCycle annotates opened and closed ports with tag labels and writes a summary.
func (h *TagHook) OnCycle(diff ports.Diff) {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return
	}
	h.logTagged("opened", diff.Opened)
	h.logTagged("closed", diff.Closed)
}

func (h *TagHook) logTagged(direction string, ps []int) {
	if len(ps) == 0 {
		return
	}
	sorted := make([]int, len(ps))
	copy(sorted, ps)
	sort.Ints(sorted)
	for _, p := range sorted {
		label := h.tagger.Label(p)
		fmt.Fprintf(h.w, "[tag] %s port %d (%s)\n", direction, p, label)
	}
}
