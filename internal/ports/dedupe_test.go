package ports

import (
	"testing"
)

func TestDeduplicatorFirstCallNotDuplicate(t *testing.T) {
	d := NewDeduplicator()
	if d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestDeduplicatorSecondCallIsDuplicate(t *testing.T) {
	d := NewDeduplicator()
	d.IsDuplicate("port:8080:opened")
	if !d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected second call with same key to be a duplicate")
	}
}

func TestDeduplicatorDifferentKeysAreIndependent(t *testing.T) {
	d := NewDeduplicator()
	d.IsDuplicate("port:8080:opened")
	if d.IsDuplicate("port:9090:opened") {
		t.Fatal("different keys should not be duplicates of each other")
	}
}

func TestDeduplicatorResetClearsState(t *testing.T) {
	d := NewDeduplicator()
	d.IsDuplicate("port:8080:opened")
	d.Reset()
	if d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected key to be cleared after Reset")
	}
}

func TestDeduplicatorLen(t *testing.T) {
	d := NewDeduplicator()
	if d.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", d.Len())
	}
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	d.IsDuplicate("a") // duplicate, should not increase count
	if d.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", d.Len())
	}
}

func TestDeduplicatorResetResetsLen(t *testing.T) {
	d := NewDeduplicator()
	d.IsDuplicate("x")
	d.Reset()
	if d.Len() != 0 {
		t.Fatalf("expected Len 0 after reset, got %d", d.Len())
	}
}
