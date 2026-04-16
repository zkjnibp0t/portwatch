package ports

import "fmt"

// Whitelist holds a set of known-safe (port, process) pairs that should
// never trigger an alert even when they appear as newly opened ports.
type Whitelist struct {
	entries map[string]struct{}
}

// WhitelistEntry describes a trusted port/process combination.
type WhitelistEntry struct {
	Port    int    // 0 means "any port"
	Process string // empty means "any process"
}

// NewWhitelist creates a Whitelist from the provided entries.
func NewWhitelist(entries []WhitelistEntry) *Whitelist {
	wl := &Whitelist{entries: make(map[string]struct{}, len(entries))}
	for _, e := range entries {
		wl.entries[entryKey(e.Port, e.Process)] = struct{}{}
	}
	return wl
}

// Allow returns true when the port+process combination is whitelisted.
// A wildcard match (port=0 or process="") is checked in addition to
// the exact match so partial rules work as expected.
func (w *Whitelist) Allow(port int, process string) bool {
	if _, ok := w.entries[entryKey(port, process)]; ok {
		return true
	}
	// port-only rule
	if _, ok := w.entries[entryKey(port, "")]; ok {
		return true
	}
	// process-only rule
	if _, ok := w.entries[entryKey(0, process)]; ok {
		return true
	}
	return false
}

// FilterDiff removes opened ports from diff that are covered by the whitelist.
// It enriches each port with process information via the provided lookup
// function before checking the whitelist.
func (w *Whitelist) FilterDiff(diff Diff, lookup func(port int) string) Diff {
	if len(w.entries) == 0 {
		return diff
	}
	filtered := make([]int, 0, len(diff.Opened))
	for _, p := range diff.Opened {
		proc := ""
		if lookup != nil {
			proc = lookup(p)
		}
		if !w.Allow(p, proc) {
			filtered = append(filtered, p)
		}
	}
	return Diff{Opened: filtered, Closed: diff.Closed}
}

// Len returns the number of entries in the whitelist.
func (w *Whitelist) Len() int {
	return len(w.entries)
}

// Add inserts a new entry into the whitelist at runtime.
func (w *Whitelist) Add(e WhitelistEntry) {
	w.entries[entryKey(e.Port, e.Process)] = struct{}{}
}

func entryKey(port int, process string) string {
	return fmt.Sprintf("%d|%s", port, process)
}
