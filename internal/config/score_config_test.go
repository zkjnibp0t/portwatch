package config

import "testing"

func TestBuildScorerPrivilegedPort(t *testing.T) {
	s := BuildScorer(3.0, 5.0)
	ps := s.Score(80)
	if ps.Score != 3.0 {
		t.Fatalf("expected 3.0 for port 80, got %f", ps.Score)
	}
}

func TestBuildScorerRiskyPort(t *testing.T) {
	s := BuildScorer(3.0, 5.0)
	// port 23 is both privileged and risky
	ps := s.Score(23)
	if ps.Score != 8.0 {
		t.Fatalf("expected 8.0 for telnet, got %f", ps.Score)
	}
	if len(ps.Reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(ps.Reasons))
	}
}

func TestBuildScorerHighPort(t *testing.T) {
	s := BuildScorer(3.0, 5.0)
	ps := s.Score(8080)
	if ps.Score != 0 {
		t.Fatalf("expected 0 for high port, got %f", ps.Score)
	}
}

func TestBuildScorerRDPIsRisky(t *testing.T) {
	s := BuildScorer(1.0, 4.0)
	ps := s.Score(3389)
	found := false
	for _, r := range ps.Reasons {
		if r == "risky service rdp" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected rdp reason, got %v", ps.Reasons)
	}
}

func TestPrivilegedContributorBoundary(t *testing.T) {
	c := PrivilegedContributor{Delta: 2.0}
	d, _ := c.Contribute(1023)
	if d != 2.0 {
		t.Fatalf("expected 2.0 for port 1023")
	}
	d2, _ := c.Contribute(1024)
	if d2 != 0 {
		t.Fatalf("expected 0 for port 1024")
	}
}
