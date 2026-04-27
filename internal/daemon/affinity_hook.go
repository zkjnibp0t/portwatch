package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/user/portwatch/internal/ports"
)

// AffinityHook observes each scan's open port set and logs port pairs
// that frequently appear together once they meet the minimum count threshold.
type AffinityHook struct {
	tracker  *ports.AffinityTracker
	minCount int
	logger   *log.Logger
}

// NewAffinityHook creates an AffinityHook. minCount is the co-occurrence
// threshold before a pair is reported.
func NewAffinityHook(minCount int, w io.Writer) *AffinityHook {
	if w == nil {
		w = os.Stdout
	}
	return &AffinityHook{
		tracker:  ports.NewAffinityTracker(minCount),
		minCount: minCount,
		logger:   log.New(w, "[affinity] ", 0),
	}
}

// BeforeScan is a no-op for this hook.
func (h *AffinityHook) BeforeScan() {}

// AfterScan receives the current open port set, observes pairs, and logs
// any pairs that have reached or exceeded the minimum co-occurrence count.
func (h *AffinityHook) AfterScan(current map[int]struct{}, diff ports.Diff) {
	open := sortedKeys(current)
	if len(open) < 2 {
		return
	}
	h.tracker.Observe(open)
	pairs := h.tracker.Pairs()
	if len(pairs) == 0 {
		return
	}
	for _, p := range pairs {
		h.logger.Printf("affinity pair detected: %s", p)
	}
}

// sortedKeys extracts sorted integer keys from a port set.
func sortedKeys(m map[int]struct{}) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// formatAffinityPairs returns a human-readable summary of affinity pairs.
func formatAffinityPairs(pairs []ports.AffinityPair) string {
	if len(pairs) == 0 {
		return "(none)"
	}
	s := ""
	for i, p := range pairs {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d<->%d", p.PortA, p.PortB)
	}
	return s
}
