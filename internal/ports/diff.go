package ports

// Diff holds the result of comparing two port sets.
type Diff struct {
	Opened []int
	Closed []int
}

// IsEmpty returns true when no ports were opened or closed.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare calculates which ports were opened and which were closed between
// the previous and current port sets.
func Compare(prev, curr Set) Diff {
	var opened, closed []int

	for p := range curr {
		if _, existed := prev[p]; !existed {
			opened = append(opened, p)
		}
	}
	for p := range prev {
		if _, exists := curr[p]; !exists {
			closed = append(closed, p)
		}
	}

	sortInts(opened)
	sortInts(closed)

	return Diff{Opened: opened, Closed: closed}
}

func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
