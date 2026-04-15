package ports

// PortEntry pairs a resolved port with optional process ownership metadata.
type PortEntry struct {
	Port    ResolvedPort
	Process ProcessInfo
}

// Enricher attaches ProcessInfo to a set of ResolvedPorts.
type Enricher struct {
	lookup func(pid int) ProcessInfo
}

// NewEnricher creates an Enricher. Pass nil to use the default LookupProcess.
func NewEnricher(lookup func(pid int) ProcessInfo) *Enricher {
	if lookup == nil {
		lookup = LookupProcess
	}
	return &Enricher{lookup: lookup}
}

// Enrich converts a slice of ResolvedPorts into PortEntries with process data.
func (e *Enricher) Enrich(ports []ResolvedPort) []PortEntry {
	entries := make([]PortEntry, 0, len(ports))
	for _, rp := range ports {
		entry := PortEntry{Port: rp}
		if rp.PID > 0 {
			entry.Process = e.lookup(rp.PID)
		}
		entries = append(entries, entry)
	}
	return entries
}

// EnrichSet converts a PortSet into PortEntries using the provided Resolver.
// Ports whose PID cannot be determined are included with an empty ProcessInfo.
func (e *Enricher) EnrichSet(set PortSet, r *Resolver) []PortEntry {
	resolved := r.ResolveSet(set)
	return e.Enrich(resolved)
}
