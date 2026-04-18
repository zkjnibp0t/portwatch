package ports

import (
	"testing"
)

func TestRollupEmptyEntries(t *testing.T) {
	r := NewRollup()
	if len(r.Entries()) != 0 {
		t.Fatal("expected empty entries")
	}
}

func TestRollupRecordOpened(t *testing.T) {
	r := NewRollup()
	r.Record(Diff{Opened: []int{80, 443}})
	entries := r.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 80 || entries[0].Opened != 1 || entries[0].Net != 1 {
		t.Errorf("unexpected entry for port 80: %+v", entries[0])
	}
}

func TestRollupRecordClosed(t *testing.T) {
	r := NewRollup()
	r.Record(Diff{Closed: []int{22}})
	entries := r.Entries()
	if entries[0].Closed != 1 || entries[0].Net != -1 {
		t.Errorf("unexpected closed entry: %+v", entries[0])
	}
}

func TestRollupAccumulatesMultipleDiffs(t *testing.T) {
	r := NewRollup()
	r.Record(Diff{Opened: []int{8080}})
	r.Record(Diff{Opened: []int{8080}})
	r.Record(Diff{Closed: []int{8080}})
	e := r.Entries()[0]
	if e.Opened != 2 || e.Closed != 1 || e.Net != 1 {
		t.Errorf("unexpected accumulated entry: %+v", e)
	}
}

func TestRollupSortedByPort(t *testing.T) {
	r := NewRollup()
	r.Record(Diff{Opened: []int{9000, 80, 443}})
	entries := r.Entries()
	for i := 1; i < len(entries); i++ {
		if entries[i].Port < entries[i-1].Port {
			t.Fatal("entries not sorted by port")
		}
	}
}

func TestRollupReset(t *testing.T) {
	r := NewRollup()
	r.Record(Diff{Opened: []int{80}})
	r.Reset()
	if len(r.Entries()) != 0 {
		t.Fatal("expected empty after reset")
	}
}
