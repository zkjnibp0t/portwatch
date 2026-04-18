package ports

import "sort"

// RollupEntry summarises activity for a single port across a window.
type RollupEntry struct {
	Port    int
	Opened  int
	Closed  int
	Net     int // Opened - Closed
}

// Rollup aggregates diff events into per-port summaries.
type Rollup struct {
	entries map[int]*RollupEntry
}

// NewRollup creates an empty Rollup.
func NewRollup() *Rollup {
	return &Rollup{entries: make(map[int]*RollupEntry)}
}

// Record folds a diff into the rollup.
func (r *Rollup) Record(diff Diff) {
	for _, p := range diff.Opened {
		e := r.entry(p)
		e.Opened++
		e.Net++
	}
	for _, p := range diff.Closed {
		e := r.entry(p)
		e.Closed++
		e.Net--
	}
}

// Entries returns all rollup entries sorted by port number.
func (r *Rollup) Entries() []RollupEntry {
	out := make([]RollupEntry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })
	return out
}

// Reset clears all accumulated data.
func (r *Rollup) Reset() {
	r.entries = make(map[int]*RollupEntry)
}

func (r *Rollup) entry(port int) *RollupEntry {
	if e, ok := r.entries[port]; ok {
		return e
	}
	e := &RollupEntry{Port: port}
	r.entries[port] = e
	return e
}
