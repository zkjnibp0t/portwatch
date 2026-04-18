package ports

import (
	"testing"
	"time"
)

func TestJittererZeroMaxPctReturnsBase(t *testing.T) {
	j := NewJitterer(JitterConfig{Base: 10 * time.Second, MaxPct: 0})
	for i := 0; i < 20; i++ {
		if got := j.Next(); got != 10*time.Second {
			t.Fatalf("expected 10s, got %v", got)
		}
	}
}

func TestJittererZeroBaseReturnsZero(t *testing.T) {
	j := NewJitterer(JitterConfig{Base: 0, MaxPct: 0.5})
	if got := j.Next(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestJittererStaysWithinBounds(t *testing.T) {
	base := 10 * time.Second
	maxPct := 0.3
	j := NewJitterer(JitterConfig{Base: base, MaxPct: maxPct})
	j.SetSeed(42)

	low := time.Duration(float64(base) * (1 - maxPct))
	high := time.Duration(float64(base) * (1 + maxPct))

	for i := 0; i < 200; i++ {
		v := j.Next()
		if v < low || v > high {
			t.Fatalf("jitter %v out of [%v, %v]", v, low, high)
		}
	}
}

func TestJittererMaxPctClamped(t *testing.T) {
	j := NewJitterer(JitterConfig{Base: 5 * time.Second, MaxPct: 2.5})
	if j.cfg.MaxPct != 1.0 {
		t.Fatalf("expected MaxPct clamped to 1.0, got %v", j.cfg.MaxPct)
	}
}

func TestJittererNegativeMaxPctClamped(t *testing.T) {
	j := NewJitterer(JitterConfig{Base: 5 * time.Second, MaxPct: -0.5})
	if j.cfg.MaxPct != 0 {
		t.Fatalf("expected MaxPct clamped to 0, got %v", j.cfg.MaxPct)
	}
}

func TestJittererDifferentValuesProduced(t *testing.T) {
	j := NewJitterer(JitterConfig{Base: 10 * time.Second, MaxPct: 0.5})
	j.SetSeed(99)
	seen := map[time.Duration]bool{}
	for i := 0; i < 50; i++ {
		seen[j.Next()] = true
	}
	if len(seen) < 5 {
		t.Fatal("expected varied jitter values, got too few distinct results")
	}
}
