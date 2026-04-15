package ports

import (
	"fmt"
	"testing"
)

func TestProcessInfoStringUnknown(t *testing.T) {
	p := ProcessInfo{}
	if got := p.String(); got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestProcessInfoStringPIDOnly(t *testing.T) {
	p := ProcessInfo{PID: 42}
	if got := p.String(); got != "42" {
		t.Errorf("expected '42', got %q", got)
	}
}

func TestProcessInfoStringFull(t *testing.T) {
	p := ProcessInfo{PID: 1234, Name: "nginx", User: "www"}
	want := "1234 nginx (www)"
	if got := p.String(); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestLookupProcessInvalidPID(t *testing.T) {
	info := LookupProcess(-1)
	if info.PID != 0 {
		t.Errorf("expected zero PID for invalid input, got %d", info.PID)
	}
}

func TestEnricherNilLookupUsesDefault(t *testing.T) {
	e := NewEnricher(nil)
	if e.lookup == nil {
		t.Fatal("expected non-nil lookup function")
	}
}

func TestEnricherEnrichNoPID(t *testing.T) {
	e := NewEnricher(func(pid int) ProcessInfo {
		return ProcessInfo{PID: pid, Name: fmt.Sprintf("proc%d", pid)}
	})
	ports := []ResolvedPort{
		{Port: 8080, Proto: "tcp", PID: 0},
	}
	entries := e.Enrich(ports)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Process.PID != 0 {
		t.Errorf("expected empty ProcessInfo for PID=0")
	}
}

func TestEnricherEnrichWithPID(t *testing.T) {
	e := NewEnricher(func(pid int) ProcessInfo {
		return ProcessInfo{PID: pid, Name: "myapp"}
	})
	ports := []ResolvedPort{
		{Port: 9090, Proto: "tcp", PID: 555},
	}
	entries := e.Enrich(ports)
	if entries[0].Process.Name != "myapp" {
		t.Errorf("expected process name 'myapp', got %q", entries[0].Process.Name)
	}
	if entries[0].Process.PID != 555 {
		t.Errorf("expected PID 555, got %d", entries[0].Process.PID)
	}
}
