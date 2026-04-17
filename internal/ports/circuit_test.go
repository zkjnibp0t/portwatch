package ports

import (
	"testing"
	"time"
)

func TestCircuitInitiallyClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	if cb.State() != CircuitClosed {
		t.Fatalf("expected closed, got %s", cb.State())
	}
	if !cb.Allow() {
		t.Fatal("expected Allow() true when closed")
	}
}

func TestCircuitOpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitClosed {
		t.Fatal("should still be closed after 2 failures")
	}
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected open after threshold, got %s", cb.State())
	}
	if cb.Allow() {
		t.Fatal("expected Allow() false when open")
	}
}

func TestCircuitHalfOpenAfterReset(t *testing.T) {
	now := time.Now()
	cb := NewCircuitBreaker(1, 500*time.Millisecond)
	cb.clock = func() time.Time { return now }
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatal("expected open")
	}
	// advance clock past resetAfter
	cb.clock = func() time.Time { return now.Add(time.Second) }
	if !cb.Allow() {
		t.Fatal("expected Allow() true in half-open")
	}
	if cb.State() != CircuitHalfOpen {
		t.Fatalf("expected half-open, got %s", cb.State())
	}
}

func TestCircuitClosedAfterSuccess(t *testing.T) {
	cb := NewCircuitBreaker(1, time.Second)
	cb.RecordFailure()
	cb.RecordSuccess()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected closed after success, got %s", cb.State())
	}
	if !cb.Allow() {
		t.Fatal("expected Allow() true after success")
	}
}

func TestCircuitStateString(t *testing.T) {
	if CircuitClosed.String() != "closed" {
		t.Error("wrong string for closed")
	}
	if CircuitOpen.String() != "open" {
		t.Error("wrong string for open")
	}
	if CircuitHalfOpen.String() != "half-open" {
		t.Error("wrong string for half-open")
	}
}

func TestCircuitStillBlocksBeforeReset(t *testing.T) {
	now := time.Now()
	cb := NewCircuitBreaker(1, time.Second)
	cb.clock = func() time.Time { return now }
	cb.RecordFailure()
	cb.clock = func() time.Time { return now.Add(100 * time.Millisecond) }
	if cb.Allow() {
		t.Fatal("should still block before resetAfter")
	}
}
