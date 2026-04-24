package ports

import (
	"sort"
	"sync"
	"time"
)

// HeatmapEntry holds the observed activity count and last-seen time for a port.
type HeatmapEntry struct {
	Port     int
	Hits     int
	LastSeen time.Time
}

// HeatmapTracker accumulates per-port open/close event counts over time,
// providing a ranked view of the most active ports.
type HeatmapTracker struct {
	mu      sync.Mutex
	entries map[int]*HeatmapEntry
	clock   func() time.Time
}

// NewHeatmapTracker returns a HeatmapTracker using the real wall clock.
func NewHeatmapTracker() *HeatmapTracker {
	return &HeatmapTracker{
		entries: make(map[int]*HeatmapEntry),
		clock:   time.Now,
	}
}

// Record increments the hit counter for each port in the diff.
func (h *HeatmapTracker) Record(diff Diff) {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := h.clock()
	for _, p := range diff.Opened {
		h.touch(p, now)
	}
	for _, p := range diff.Closed {
		h.touch(p, now)
	}
}

func (h *HeatmapTracker) touch(port int, t time.Time) {
	e, ok := h.entries[port]
	if !ok {
		e = &HeatmapEntry{Port: port}
		h.entries[port] = e
	}
	e.Hits++
	e.LastSeen = t
}

// TopN returns the n most active ports sorted by hit count descending.
// If n <= 0 all entries are returned.
func (h *HeatmapTracker) TopN(n int) []HeatmapEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]HeatmapEntry, 0, len(h.entries))
	for _, e := range h.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Hits != out[j].Hits {
			return out[i].Hits > out[j].Hits
		}
		return out[i].Port < out[j].Port
	})
	if n > 0 && n < len(out) {
		return out[:n]
	}
	return out
}

// Reset clears all heatmap data.
func (h *HeatmapTracker) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = make(map[int]*HeatmapEntry)
}
