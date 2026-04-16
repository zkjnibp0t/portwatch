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

// Compare returns ports that appeared in next but not prev (Opened)
// and ports that were in prev but not next (Closed).
func Compare(prev, next Set) Diff {
	var opened, closed []int

	for p := range next {
		if !prev[p] {
			opened = append(opened, p)
		}
	}
	for p := range prev {
		if !next[p] {
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
