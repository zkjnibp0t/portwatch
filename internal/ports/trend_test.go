package ports

import (
	"testing"
	"time"
)

func TestTrendTrackerRecordAndFlap(t *testing.T) {
	tr := NewTrendTracker(time.Hour)
	tr.Record(8080, "opened")
	tr.Record(8080, "closed")
	tr.Record(8080, "opened")

	if got := tr.FlapCount(8080); got != 3 {
		t.Fatalf("expected 3 flaps, got %d", got)
	}
}

func TestTrendTrackerDifferentPortsAreIndependent(t *testing.T) {
	tr := NewTrendTracker(time.Hour)
	tr.Record(80, "opened")
	tr.Record(443, "opened")
	tr.Record(443, "closed")

	if got := tr.FlapCount(80); got != 1 {
		t.Fatalf("port 80: expected 1, got %d", got)
	}
	if got := tr.FlapCount(443); got != 2 {
		t.Fatalf("port 443: expected 2, got %d", got)
	}
}

func TestTrendTrackerPrunesOldEntries(t *testing.T) {
	tr := NewTrendTracker(50 * time.Millisecond)
	tr.Record(9000, "opened")

	time.Sleep(80 * time.Millisecond)

	// Trigger prune via FlapCount
	if got := tr.FlapCount(9000); got != 0 {
		t.Fatalf("expected 0 after pruning, got %d", got)
	}
}

func TestTrendTrackerSummaryEmpty(t *testing.T) {
	tr := NewTrendTracker(time.Hour)
	if lines := tr.Summary(); len(lines) != 0 {
		t.Fatalf("expected empty summary, got %v", lines)
	}
}

func TestTrendTrackerSummaryContainsPort(t *testing.T) {
	tr := NewTrendTracker(time.Hour)
	tr.Record(3000, "opened")
	tr.Record(3000, "closed")

	lines := tr.Summary()
	if len(lines) != 1 {
		t.Fatalf("expected 1 summary line, got %d", len(lines))
	}
}

func TestTrendTrackerZeroMaxAgeDefaultsToDay(t *testing.T) {
	tr := NewTrendTracker(0)
	if tr.maxAge != 24*time.Hour {
		t.Fatalf("expected 24h default, got %v", tr.maxAge)
	}
}
