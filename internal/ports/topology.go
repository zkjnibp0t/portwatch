package ports

import (
	"fmt"
	"sort"
	"sync"
)

// TopologyEdge represents a directional relationship between two ports.
type TopologyEdge struct {
	From int
	To   int
}

func (e TopologyEdge) String() string {
	return fmt.Sprintf("%d->%d", e.From, e.To)
}

// TopologyTracker records co-occurrence relationships between ports
// that open or close within the same diff cycle.
type TopologyTracker struct {
	mu    sync.Mutex
	edges map[TopologyEdge]int
}

// NewTopologyTracker returns an initialised TopologyTracker.
func NewTopologyTracker() *TopologyTracker {
	return &TopologyTracker{edges: make(map[TopologyEdge]int)}
}

// Record registers all pairwise edges for a set of ports seen together.
func (t *TopologyTracker) Record(ports []int) {
	if len(ports) < 2 {
		return
	}
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			e := TopologyEdge{From: sorted[i], To: sorted[j]}
			t.edges[e]++
		}
	}
}

// Edges returns all recorded edges with their co-occurrence counts.
func (t *TopologyTracker) Edges() map[TopologyEdge]int {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[TopologyEdge]int, len(t.edges))
	for k, v := range t.edges {
		out[k] = v
	}
	return out
}

// Neighbors returns all ports that have co-occurred with the given port.
func (t *TopologyTracker) Neighbors(port int) []int {
	t.mu.Lock()
	defer t.mu.Unlock()
	seen := map[int]struct{}{}
	for e := range t.edges {
		if e.From == port {
			seen[e.To] = struct{}{}
		} else if e.To == port {
			seen[e.From] = struct{}{}
		}
	}
	out := make([]int, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

// Reset clears all recorded topology data.
func (t *TopologyTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.edges = make(map[TopologyEdge]int)
}
