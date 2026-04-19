package ports

import (
	"testing"
)

func TestSequenceTrackerInitialZero(t *testing.T) {
	st := NewSequenceTracker()
	if got := st.LastFor(80); got != 0 {
		t.Fatalf("expected 0 for unseen port, got %d", got)
	}
}

func TestSequenceTrackerNextIncrements(t *testing.T) {
	st := NewSequenceTracker()
	a := st.Next(80)
	b := st.Next(443)
	if a >= b {
		t.Fatalf("expected a < b, got a=%d b=%d", a, b)
	}
}

func TestSequenceTrackerLastForRecorded(t *testing.T) {
	st := NewSequenceTracker()
	seq := st.Next(8080)
	if got := st.LastFor(8080); got != seq {
		t.Fatalf("expected %d, got %d", seq, got)
	}
}

func TestSequenceTrackerMultiplePortsIndependent(t *testing.T) {
	st := NewSequenceTracker()
	st.Next(22)
	st.Next(22)
	st.Next(80)
	if st.LastFor(22) == st.LastFor(80) {
		t.Fatal("expected different sequence numbers for different ports")
	}
}

func TestSequenceTrackerLen(t *testing.T) {
	st := NewSequenceTracker()
	st.Next(22)
	st.Next(80)
	st.Next(22) // duplicate port
	if got := st.Len(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestSequenceTrackerReset(t *testing.T) {
	st := NewSequenceTracker()
	st.Next(22)
	st.Reset()
	if got := st.LastFor(22); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := st.Len(); got != 0 {
		t.Fatalf("expected len 0 after reset, got %d", got)
	}
	first := st.Next(22)
	if first != 1 {
		t.Fatalf("expected counter to restart at 1, got %d", first)
	}
}
