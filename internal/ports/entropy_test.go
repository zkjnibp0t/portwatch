package ports

import (
	"math"
	"testing"
)

func makePortSet(ports ...int) map[int]struct{} {
	s := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		s[p] = struct{}{}
	}
	return s
}

func TestEntropyEmptyHistory(t *testing.T) {
	ec := NewEntropyCalculator(5)
	if got := ec.Entropy(); got != 0 {
		t.Errorf("expected 0, got %f", got)
	}
}

func TestEntropySingleSnapshot(t *testing.T) {
	ec := NewEntropyCalculator(5)
	ec.Record(makePortSet(80, 443))
	// Two equally frequent ports: entropy = log2(2) = 1.0
	got := ec.Entropy()
	if math.Abs(got-1.0) > 1e-9 {
		t.Errorf("expected 1.0, got %f", got)
	}
}

func TestEntropyWindowEvictsOldest(t *testing.T) {
	ec := NewEntropyCalculator(2)
	ec.Record(makePortSet(80))
	ec.Record(makePortSet(443))
	ec.Record(makePortSet(8080)) // evicts first
	if len(ec.history) != 2 {
		t.Errorf("expected window size 2, got %d", len(ec.history))
	}
	// Only 443 and 8080 remain
	top := ec.TopPorts(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top ports, got %d", len(top))
	}
}

func TestEntropyTopPorts(t *testing.T) {
	ec := NewEntropyCalculator(10)
	ec.Record(makePortSet(80, 443))
	ec.Record(makePortSet(80, 8080))
	ec.Record(makePortSet(80))
	top := ec.TopPorts(1)
	if len(top) != 1 || top[0] != 80 {
		t.Errorf("expected top port 80, got %v", top)
	}
}

func TestEntropyTopPortsFewerThanN(t *testing.T) {
	ec := NewEntropyCalculator(5)
	ec.Record(makePortSet(22))
	top := ec.TopPorts(5)
	if len(top) != 1 {
		t.Errorf("expected 1 port, got %d", len(top))
	}
}

func TestEntropyZeroWindowSizeClamped(t *testing.T) {
	ec := NewEntropyCalculator(0)
	ec.Record(makePortSet(80))
	ec.Record(makePortSet(443))
	if len(ec.history) != 1 {
		t.Errorf("expected window clamped to 1, got %d", len(ec.history))
	}
}
