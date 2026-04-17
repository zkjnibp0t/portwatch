package ports

import (
	"fmt"
	"sort"
	"strings"
)

// Fingerprint is a stable hash-like string representing a set of resolved ports.
type Fingerprint string

// Fingerprintер computes fingerprints from resolved port sets.
type Fingerprinter struct{}

// NewFingerprinter returns a new Fingerprinter.
func NewFingerprinter() *Fingerprinter {
	return &Fingerprinter{}
}

// Compute returns a deterministic fingerprint for the given resolved port set.
// The fingerprint encodes port:pid pairs sorted by port number.
func (f *Fingerprinter) Compute(ports []ResolvedPort) Fingerprint {
	if len(ports) == 0 {
		return Fingerprint("empty")
	}

	type entry struct {
		port int
		pid  int
	}

	entries := make([]entry, 0, len(ports))
	for _, p := range ports {
		entries = append(entries, entry{port: p.Port, pid: p.PID})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].port != entries[j].port {
			return entries[i].port < entries[j].port
		}
		return entries[i].pid < entries[j].pid
	})

	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		parts = append(parts, fmt.Sprintf("%d:%d", e.port, e.pid))
	}
	return Fingerprint(strings.Join(parts, ","))
}

// Equal returns true if two fingerprints match.
func (f *Fingerprinter) Equal(a, b Fingerprint) bool {
	return a == b
}
