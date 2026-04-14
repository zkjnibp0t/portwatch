package ports

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Opened []PortState
	Closed []PortState
}

// HasChanges returns true when at least one port was opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare calculates which ports were opened or closed between two snapshots.
// previous and current are slices returned by Scanner.Scan().
func Compare(previous, current []PortState) Diff {
	prevSet := ToSet(previous)
	currSet := ToSet(current)

	var diff Diff

	// Ports present in current but not in previous → newly opened.
	for port, state := range currSet {
		if _, exists := prevSet[port]; !exists {
			diff.Opened = append(diff.Opened, state)
		}
	}

	// Ports present in previous but not in current → closed.
	for port, state := range prevSet {
		if _, exists := currSet[port]; !exists {
			diff.Closed = append(diff.Closed, state)
		}
	}

	return diff
}
