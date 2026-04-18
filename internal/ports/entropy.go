package ports

import (
	"math"
	"sort"
)

// EntropyCalculator measures the Shannon entropy of a port set,
// which can indicate unusual diversity in open ports.
type EntropyCalculator struct {
	windowSize int
	history    []map[int]struct{}
}

// NewEntropyCalculator creates a calculator that tracks up to windowSize snapshots.
func NewEntropyCalculator(windowSize int) *EntropyCalculator {
	if windowSize < 1 {
		windowSize = 1
	}
	return &EntropyCalculator{windowSize: windowSize}
}

// Record adds a port snapshot to the rolling window.
func (e *EntropyCalculator) Record(ports map[int]struct{}) {
	copy := make(map[int]struct{}, len(ports))
	for p := range ports {
		copy[p] = struct{}{}
	}
	e.history = append(e.history, copy)
	if len(e.history) > e.windowSize {
		e.history = e.history[len(e.history)-e.windowSize:]
	}
}

// Entropy computes Shannon entropy over port frequencies across the window.
// Returns 0 if no history is recorded.
func (e *EntropyCalculator) Entropy() float64 {
	if len(e.history) == 0 {
		return 0
	}
	freq := make(map[int]int)
	total := 0
	for _, snap := range e.history {
		for p := range snap {
			freq[p]++
			total++
		}
	}
	if total == 0 {
		return 0
	}
	entropy := 0.0
	for _, count := range freq {
		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// TopPorts returns the n most frequently seen ports across the window.
func (e *EntropyCalculator) TopPorts(n int) []int {
	freq := make(map[int]int)
	for _, snap := range e.history {
		for p := range snap {
			freq[p]++
		}
	}
	type kv struct {
		port, count int
	}
	var sorted []kv
	for p, c := range freq {
		sorted = append(sorted, kv{p, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].count != sorted[j].count {
			return sorted[i].count > sorted[j].count
		}
		return sorted[i].port < sorted[j].port
	})
	result := make([]int, 0, n)
	for i := 0; i < n && i < len(sorted); i++ {
		result = append(result, sorted[i].port)
	}
	return result
}
