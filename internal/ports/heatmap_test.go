package ports

import (
	"testing"
	"time"
)

func fixedHeatmapClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestHeatmapEmptyTopN(t *testing.T) {
	h := NewHeatmapTracker()
	if got := h.TopN(5); len(got) != 0 {
		t.Fatalf("expected empty, got %d entries", len(got))
	}
}

func TestHeatmapRecordsOpenedAndClosed(t *testing.T) {
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(time.Now())
	h.Record(Diff{Opened: []int{80, 443}, Closed: []int{8080}})
	top := h.TopN(0)
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
}

func TestHeatmapAccumulatesHits(t *testing.T) {
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(time.Now())
	h.Record(Diff{Opened: []int{80}})
	h.Record(Diff{Opened: []int{80}})
	h.Record(Diff{Closed: []int{80}})
	top := h.TopN(1)
	if top[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", top[0].Port)
	}
	if top[0].Hits != 3 {
		t.Fatalf("expected 3 hits, got %d", top[0].Hits)
	}
}

func TestHeatmapTopNRanksCorrectly(t *testing.T) {
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(time.Now())
	h.Record(Diff{Opened: []int{443}})
	h.Record(Diff{Opened: []int{443}})
	h.Record(Diff{Opened: []int{80}})
	top := h.TopN(1)
	if top[0].Port != 443 {
		t.Fatalf("expected port 443 at rank 0, got %d", top[0].Port)
	}
}

func TestHeatmapTopNLimitsResults(t *testing.T) {
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(time.Now())
	h.Record(Diff{Opened: []int{80, 443, 8080, 9090}})
	top := h.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
}

func TestHeatmapResetClearsEntries(t *testing.T) {
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(time.Now())
	h.Record(Diff{Opened: []int{80}})
	h.Reset()
	if got := h.TopN(0); len(got) != 0 {
		t.Fatalf("expected empty after reset, got %d", len(got))
	}
}

func TestHeatmapLastSeenUpdated(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	h := NewHeatmapTracker()
	h.clock = fixedHeatmapClock(now)
	h.Record(Diff{Opened: []int{80}})
	top := h.TopN(1)
	if !top[0].LastSeen.Equal(now) {
		t.Fatalf("expected LastSeen %v, got %v", now, top[0].LastSeen)
	}
}
