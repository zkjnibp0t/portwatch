package ports

import "sync"

// Deduplicator suppresses repeated identical diffs within a short window,
// preventing the same change from being dispatched more than once per cycle.
type Deduplicator struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewDeduplicator returns a fresh Deduplicator.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]struct{}),
	}
}

// IsDuplicate returns true if the given key was already seen since the last
// Reset call. If not seen, it records the key and returns false.
func (d *Deduplicator) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.seen[key]; ok {
		return true
	}
	d.seen[key] = struct{}{}
	return false
}

// Reset clears all recorded keys, typically called at the start of each scan cycle.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
}

// Len returns the number of unique keys seen since the last reset.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
