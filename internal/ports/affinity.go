package ports

import (
	"fmt"
	"sort"
	"sync"
)

// AffinityPair represents two ports that frequently appear together.
type AffinityPair struct {
	PortA int
	PortB int
	Count int
}

func (p AffinityPair) String() string {
	return fmt.Sprintf("%d<->%d (count=%d)", p.PortA, p.PortB, p.Count)
}

// AffinityTracker records how often pairs of ports are seen open together.
type AffinityTracker struct {
	mu      sync.Mutex
	counts  map[string]int
	pairKey func(a, b int) string
	minCount int
}

// NewAffinityTracker creates an AffinityTracker that reports pairs seen
// together at least minCount times.
func NewAffinityTracker(minCount int) *AffinityTracker {
	if minCount < 1 {
		minCount = 1
	}
	return &AffinityTracker{
		counts:   make(map[string]int),
		minCount: minCount,
		pairKey: func(a, b int) string {
			if a > b {
				a, b = b, a
			}
			return fmt.Sprintf("%d:%d", a, b)
		},
	}
}

// Observe records all co-occurring port pairs from the given open set.
func (t *AffinityTracker) Observe(ports []int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			t.counts[t.pairKey(sorted[i], sorted[j])]++
		}
	}
}

// Pairs returns all pairs whose co-occurrence count meets the minimum threshold,
// sorted by count descending.
func (t *AffinityTracker) Pairs() []AffinityPair {
	t.mu.Lock()
	defer t.mu.Unlock()
	var result []AffinityPair
	for key, count := range t.counts {
		if count < t.minCount {
			continue
		}
		var a, b int
		fmt.Sscanf(key, "%d:%d", &a, &b)
		result = append(result, AffinityPair{PortA: a, PortB: b, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].PortA < result[j].PortA
	})
	return result
}

// Reset clears all recorded affinity data.
func (t *AffinityTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]int)
}
