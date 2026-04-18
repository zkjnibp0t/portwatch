package ports

import (
	"testing"
)

type fixedContributor struct {
	delta  float64
	reason string
}

func (f fixedContributor) Contribute(_ int) (float64, string) {
	return f.delta, f.reason
}

type portMatchContributor struct {
	target int
	delta  float64
	reason string
}

func (p portMatchContributor) Contribute(port int) (float64, string) {
	if port == p.target {
		return p.delta, p.reason
	}
	return 0, ""
}

func TestScorerZeroWithNoContributors(t *testing.T) {
	s := NewScorer()
	ps := s.Score(80)
	if ps.Score != 0 {
		t.Fatalf("expected 0, got %f", ps.Score)
	}
}

func TestScorerAggregatesContributors(t *testing.T) {
	s := NewScorer(
		fixedContributor{2.5, "high traffic"},
		fixedContributor{1.0, "known vuln"},
	)
	ps := s.Score(443)
	if ps.Score != 3.5 {
		t.Fatalf("expected 3.5, got %f", ps.Score)
	}
	if len(ps.Reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(ps.Reasons))
	}
}

func TestScorerIgnoresZeroDelta(t *testing.T) {
	s := NewScorer(fixedContributor{0, "irrelevant"})
	ps := s.Score(22)
	if len(ps.Reasons) != 0 {
		t.Fatalf("expected no reasons, got %v", ps.Reasons)
	}
}

func TestScoreSetSortedDescending(t *testing.T) {
	s := NewScorer(
		portMatchContributor{22, 5.0, "ssh"},
		portMatchContributor{80, 1.0, "http"},
	)
	set := map[int]struct{}{80: {}, 22: {}, 443: {}}
	results := s.ScoreSet(set)
	if results[0].Port != 22 {
		t.Fatalf("expected port 22 first, got %d", results[0].Port)
	}
	if results[1].Port != 80 {
		t.Fatalf("expected port 80 second, got %d", results[1].Port)
	}
}

func TestPortScoreString(t *testing.T) {
	ps := PortScore{Port: 8080, Score: 2.1, Reasons: []string{"r1"}}
	got := ps.String()
	if got == "" {
		t.Fatal("expected non-empty string")
	}
}
