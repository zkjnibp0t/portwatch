package ports

import (
	"testing"
)

func TestFingerprintEmptySet(t *testing.T) {
	fp := NewFingerprinter()
	result := fp.Compute(nil)
	if result != "empty" {
		t.Errorf("expected 'empty', got %q", result)
	}
}

func TestFingerprintDeterministic(t *testing.T) {
	fp := NewFingerprinter()
	ports := []ResolvedPort{
		{Port: 8080, PID: 100},
		{Port: 443, PID: 200},
	}
	a := fp.Compute(ports)
	b := fp.Compute(ports)
	if a != b {
		t.Errorf("fingerprint not deterministic: %q vs %q", a, b)
	}
}

func TestFingerprintOrderIndependent(t *testing.T) {
	fp := NewFingerprinter()
	a := fp.Compute([]ResolvedPort{
		{Port: 80, PID: 1},
		{Port: 443, PID: 2},
	})
	b := fp.Compute([]ResolvedPort{
		{Port: 443, PID: 2},
		{Port: 80, PID: 1},
	})
	if a != b {
		t.Errorf("fingerprint should be order-independent: %q vs %q", a, b)
	}
}

func TestFingerprintDiffersOnChange(t *testing.T) {
	fp := NewFingerprinter()
	a := fp.Compute([]ResolvedPort{{Port: 80, PID: 1}})
	b := fp.Compute([]ResolvedPort{{Port: 8080, PID: 1}})
	if fp.Equal(a, b) {
		t.Error("expected fingerprints to differ")
	}
}

func TestFingerprintEqualSameSet(t *testing.T) {
	fp := NewFingerprinter()
	a := fp.Compute([]ResolvedPort{{Port: 22, PID: 55}})
	b := fp.Compute([]ResolvedPort{{Port: 22, PID: 55}})
	if !fp.Equal(a, b) {
		t.Error("expected fingerprints to be equal")
	}
}

// TestFingerprintDiffersOnPIDChange verifies that two port sets with the same
// port numbers but different PIDs produce distinct fingerprints.
func TestFingerprintDiffersOnPIDChange(t *testing.T) {
	fp := NewFingerprinter()
	a := fp.Compute([]ResolvedPort{{Port: 80, PID: 1}})
	b := fp.Compute([]ResolvedPort{{Port: 80, PID: 2}})
	if fp.Equal(a, b) {
		t.Error("expected fingerprints to differ when PID changes")
	}
}
