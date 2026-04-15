package ports

import (
	"testing"
)

func TestFilterAllowNoRules(t *testing.T) {
	f := NewFilter(nil, nil)
	for _, port := range []int{80, 443, 8080, 9999} {
		if !f.Allow(port) {
			t.Errorf("expected port %d to be allowed with no rules", port)
		}
	}
}

func TestFilterExcludeTakesPrecedence(t *testing.T) {
	f := NewFilter([]int{80, 443}, []int{80})
	if f.Allow(80) {
		t.Error("expected port 80 to be excluded even though it is in include list")
	}
	if !f.Allow(443) {
		t.Error("expected port 443 to be allowed")
	}
}

func TestFilterIncludeListRestrictsOthers(t *testing.T) {
	f := NewFilter([]int{443, 8080}, nil)
	if !f.Allow(443) {
		t.Error("expected 443 to be allowed")
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be blocked (not in include list)")
	}
}

func TestFilterApply(t *testing.T) {
	ps := PortSet{
		80:   {},
		443:  {},
		8080: {},
		9090: {},
	}
	f := NewFilter([]int{443, 8080}, []int{8080})
	result := f.Apply(ps)

	if _, ok := result[443]; !ok {
		t.Error("expected 443 in filtered result")
	}
	if _, ok := result[8080]; ok {
		t.Error("expected 8080 to be excluded from result")
	}
	if _, ok := result[80]; ok {
		t.Error("expected 80 to be excluded (not in include list)")
	}
	if _, ok := result[9090]; ok {
		t.Error("expected 9090 to be excluded (not in include list)")
	}
}

func TestFilterApplyEmptySet(t *testing.T) {
	f := NewFilter(nil, []int{80})
	result := f.Apply(PortSet{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestFilterString(t *testing.T) {
	f := NewFilter([]int{80, 443}, []int{22})
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string from Filter.String()")
	}
}
