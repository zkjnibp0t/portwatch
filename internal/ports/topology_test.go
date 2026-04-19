package ports

import (
	"testing"
)

func TestTopologyEdgeString(t *testing.T) {
	e := TopologyEdge{From: 80, To: 443}
	if e.String() != "80->443" {
		t.Fatalf("unexpected string: %s", e.String())
	}
}

func TestTopologyRecordPair(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80, 443})
	edges := tr.Edges()
	if edges[TopologyEdge{From: 80, To: 443}] != 1 {
		t.Fatal("expected edge 80->443 with count 1")
	}
}

func TestTopologyRecordSinglePortNoEdge(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80})
	if len(tr.Edges()) != 0 {
		t.Fatal("expected no edges for single port")
	}
}

func TestTopologyRecordAccumulatesCounts(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80, 443})
	tr.Record([]int{443, 80})
	edges := tr.Edges()
	if edges[TopologyEdge{From: 80, To: 443}] != 2 {
		t.Fatalf("expected count 2, got %d", edges[TopologyEdge{From: 80, To: 443}])
	}
}

func TestTopologyRecordThreePorts(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{22, 80, 443})
	edges := tr.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
}

func TestTopologyNeighbors(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80, 443, 8080})
	neighbors := tr.Neighbors(80)
	if len(neighbors) != 2 {
		t.Fatalf("expected 2 neighbors for port 80, got %d", len(neighbors))
	}
	if neighbors[0] != 443 || neighbors[1] != 8080 {
		t.Fatalf("unexpected neighbors: %v", neighbors)
	}
}

func TestTopologyNeighborsUnknownPort(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80, 443})
	if len(tr.Neighbors(9999)) != 0 {
		t.Fatal("expected no neighbors for unknown port")
	}
}

func TestTopologyReset(t *testing.T) {
	tr := NewTopologyTracker()
	tr.Record([]int{80, 443})
	tr.Reset()
	if len(tr.Edges()) != 0 {
		t.Fatal("expected empty edges after reset")
	}
}
