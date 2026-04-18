package ports

import "fmt"

// PortScore holds a composite risk score for a port.
type PortScore struct {
	Port     int
	Score    float64
	Reasons  []string
}

func (s PortScore) String() string {
	return fmt.Sprintf("port=%d score=%.2f reasons=%v", s.Port, s.Score, s.Reasons)
}

// Scorer computes a risk score for a port based on pluggable contributors.
type Scorer struct {
	contributors []ScoreContributor
}

// ScoreContributor adds a partial score and optional reason for a port.
type ScoreContributor interface {
	Contribute(port int) (delta float64, reason string)
}

// NewScorer creates a Scorer with the given contributors.
func NewScorer(contributors ...ScoreContributor) *Scorer {
	return &Scorer{contributors: contributors}
}

// Score computes the aggregate risk score for a port.
func (s *Scorer) Score(port int) PortScore {
	ps := PortScore{Port: port}
	for _, c := range s.contributors {
		delta, reason := c.Contribute(port)
		if delta != 0 {
			ps.Score += delta
			if reason != "" {
				ps.Reasons = append(ps.Reasons, reason)
			}
		}
	}
	return ps
}

// ScoreSet scores all ports in a set and returns them sorted by descending score.
func (s *Scorer) ScoreSet(ports map[int]struct{}) []PortScore {
	results := make([]PortScore, 0, len(ports))
	for p := range ports {
		results = append(results, s.Score(p))
	}
	// simple insertion sort – sets are small
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results
}
