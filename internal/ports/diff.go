package ports

import "sort"

// Diff holds the ports that appeared or disappeared between two scans.
type Diff struct {
	Opened []int
	Closed []int
}

// IsEmpty returns true when no changes are present.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare returns the Diff between a previous and current port set.
func Compare(prev, curr Set) Diff {
	var opened, closed []int
	for p := range curr {
		if !prev[p] {
			opened = append(opened, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			closed = append(closed, p)
		}
	}
	sortInts(opened)
	sortInts(closed)
	return Diff{Opened: opened, Closed: closed}
}

func sortInts(s []int) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}
