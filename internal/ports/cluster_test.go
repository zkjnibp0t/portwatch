package ports

import (
	"testing"
)

func TestClusterNotDetectedBelowSupport(t *testing.T) {
	c := NewClusterDetector(3)
	c.Record([]int{80, 443})
	c.Record([]int{80, 443})
	if got := c.Clusters(); len(got) != 0 {
		t.Fatalf("expected no clusters, got %d", len(got))
	}
}

func TestClusterDetectedAtSupport(t *testing.T) {
	c := NewClusterDetector(2)
	c.Record([]int{80, 443})
	c.Record([]int{80, 443})
	clusters := c.Clusters()
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	if clusters[0].Count != 2 {
		t.Errorf("expected count 2, got %d", clusters[0].Count)
	}
	if len(clusters[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(clusters[0].Ports))
	}
}

func TestClusterSinglePortIgnored(t *testing.T) {
	c := NewClusterDetector(1)
	c.Record([]int{80})
	c.Record([]int{80})
	if got := c.Clusters(); len(got) != 0 {
		t.Fatalf("expected no clusters for single port, got %d", len(got))
	}
}

func TestClusterDifferentGroupsAreIndependent(t *testing.T) {
	c := NewClusterDetector(2)
	c.Record([]int{80, 443})
	c.Record([]int{80, 443})
	c.Record([]int{8080, 9090})
	c.Record([]int{8080, 9090})
	clusters := c.Clusters()
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}

func TestClusterSortedByCountDesc(t *testing.T) {
	c := NewClusterDetector(1)
	c.Record([]int{8080, 9090})
	c.Record([]int{80, 443})
	c.Record([]int{80, 443})
	c.Record([]int{80, 443})
	clusters := c.Clusters()
	if len(clusters) < 2 {
		t.Fatal("expected at least 2 clusters")
	}
	if clusters[0].Count < clusters[1].Count {
		t.Errorf("clusters not sorted by count desc: %d < %d", clusters[0].Count, clusters[1].Count)
	}
}

func TestClusterResetClearsState(t *testing.T) {
	c := NewClusterDetector(1)
	c.Record([]int{80, 443})
	c.Reset()
	if got := c.Clusters(); len(got) != 0 {
		t.Fatalf("expected empty after reset, got %d", len(got))
	}
}

func TestClusterOrderIndependent(t *testing.T) {
	c := NewClusterDetector(2)
	c.Record([]int{443, 80})
	c.Record([]int{80, 443})
	clusters := c.Clusters()
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster regardless of order, got %d", len(clusters))
	}
}
