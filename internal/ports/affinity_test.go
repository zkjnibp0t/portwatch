package ports

import (
	"testing"
)

func TestAffinityPairString(t *testing.T) {
	p := AffinityPair{PortA: 80, PortB: 443, Count: 5}
	got := p.String()
	if got != "80<->443 (count=5)" {
		t.Errorf("unexpected string: %s", got)
	}
}

func TestAffinityNoPairsWhenBelowMinCount(t *testing.T) {
	at := NewAffinityTracker(3)
	at.Observe([]int{80, 443})
	at.Observe([]int{80, 443})
	pairs := at.Pairs()
	if len(pairs) != 0 {
		t.Errorf("expected no pairs below minCount, got %d", len(pairs))
	}
}

func TestAffinityPairsAtMinCount(t *testing.T) {
	at := NewAffinityTracker(2)
	at.Observe([]int{80, 443})
	at.Observe([]int{80, 443})
	pairs := at.Pairs()
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0].PortA != 80 || pairs[0].PortB != 443 {
		t.Errorf("unexpected pair: %v", pairs[0])
	}
	if pairs[0].Count != 2 {
		t.Errorf("expected count 2, got %d", pairs[0].Count)
	}
}

func TestAffinityMultiplePairsSortedByCount(t *testing.T) {
	at := NewAffinityTracker(1)
	at.Observe([]int{22, 80, 443})
	at.Observe([]int{80, 443})
	at.Observe([]int{80, 443})
	pairs := at.Pairs()
	if len(pairs) < 1 {
		t.Fatal("expected at least one pair")
	}
	// 80<->443 should have count=3, highest
	if pairs[0].PortA != 80 || pairs[0].PortB != 443 {
		t.Errorf("expected 80<->443 first, got %v", pairs[0])
	}
	if pairs[0].Count != 3 {
		t.Errorf("expected count 3, got %d", pairs[0].Count)
	}
}

func TestAffinityOrderIndependentKey(t *testing.T) {
	at := NewAffinityTracker(1)
	at.Observe([]int{443, 80})
	at.Observe([]int{80, 443})
	pairs := at.Pairs()
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0].Count != 2 {
		t.Errorf("expected count 2 for order-independent key, got %d", pairs[0].Count)
	}
}

func TestAffinityResetClearsData(t *testing.T) {
	at := NewAffinityTracker(1)
	at.Observe([]int{80, 443})
	at.Reset()
	pairs := at.Pairs()
	if len(pairs) != 0 {
		t.Errorf("expected no pairs after reset, got %d", len(pairs))
	}
}

func TestAffinitySinglePortNoPairs(t *testing.T) {
	at := NewAffinityTracker(1)
	at.Observe([]int{80})
	pairs := at.Pairs()
	if len(pairs) != 0 {
		t.Errorf("expected no pairs from single port, got %d", len(pairs))
	}
}
