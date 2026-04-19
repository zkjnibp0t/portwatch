package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/user/portwatch/internal/ports"
)

// TopologyHook records port co-occurrence topology from each diff and logs
// the top correlated pairs when the opened set has more than one port.
type TopologyHook struct {
	tracker *ports.TopologyTracker
	minPair int
	log     *log.Logger
}

// NewTopologyHook creates a TopologyHook that logs pairs seen at least
// minPair times together.
func NewTopologyHook(minPair int, w io.Writer) *TopologyHook {
	if w == nil {
		w = os.Stdout
	}
	return &TopologyHook{
		tracker: ports.NewTopologyTracker(),
		minPair: minPair,
		log:     log.New(w, "[topology] ", 0),
	}
}

// BeforeScan is a no-op.
func (h *TopologyHook) BeforeScan() {}

// AfterScan records opened ports into the topology tracker and logs
// any edge pairs that meet or exceed the minPair threshold.
func (h *TopologyHook) AfterScan(diff ports.Diff) {
	if len(diff.Opened) < 2 {
		return
	}

	opened := make([]int, 0, len(diff.Opened))
	for p := range diff.Opened {
		opened = append(opened, p)
	}
	sort.Ints(opened)
	h.tracker.Record(opened)

	edges := h.tracker.Edges()
	type kv struct {
		edge  ports.TopologyEdge
		count int
	}
	var notable []kv
	for e, c := range edges {
		if c >= h.minPair {
			notable = append(notable, kv{e, c})
		}
	}
	if len(notable) == 0 {
		return
	}
	sort.Slice(notable, func(i, j int) bool {
		if notable[i].count != notable[j].count {
			return notable[i].count > notable[j].count
		}
		return notable[i].edge.String() < notable[j].edge.String()
	})
	for _, n := range notable {
		h.log.Println(fmt.Sprintf("co-occurrence %s seen %d time(s)", n.edge, n.count))
	}
}
