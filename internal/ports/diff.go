package ports

import "sort"

// Diff holds the result of comparing two port sets.
type Diff struct {
	Opened []int
	Closed []int
}

// IsEmpty returns true when no ports changed.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare returns the ports opened and closed between prev and curr.
func Compare(prev, curr map[int]struct{}) Diff {
	var opened, closed []int

	for p := range curr {
		if _, ok := prev[p]; !ok {
			opened = append(opened, p)
		}
	}
	for p := range prev {
		if _, ok := curr[p]; !ok {
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
