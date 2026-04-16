package daemon

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestSuppressHook() (*SuppressHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	sup := ports.NewSuppressor()
	return NewSuppressHook(sup, buf), buf
}

func TestSuppressHookAllowsUnsuppressed(t *testing.T) {
	h, _ := newTestSuppressHook()
	result := h.FilterOpened([]int{80, 443}, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
}

func TestSuppressHookFiltersOpened(t *testing.T) {
	h, buf := newTestSuppressHook()
	h.SuppressFor(8080, "", 5*time.Minute)
	buf.Reset()
	result := h.FilterOpened([]int{8080, 9090}, nil)
	if len(result) != 1 || result[0] != 9090 {
		t.Fatalf("expected [9090], got %v", result)
	}
	if buf.Len() == 0 {
		t.Fatal("expected log output for suppressed port")
	}
}

func TestSuppressHookFiltersClosed(t *testing.T) {
	h, _ := newTestSuppressHook()
	h.SuppressFor(22, "", 5*time.Minute)
	result := h.FilterClosed([]int{22, 80}, nil)
	if len(result) != 1 || result[0] != 80 {
		t.Fatalf("expected [80], got %v", result)
	}
}

func TestSuppressHookProcessAware(t *testing.T) {
	h, _ := newTestSuppressHook()
	h.SuppressFor(443, "nginx", 5*time.Minute)
	proc := func(p int) string { return "sshd" }
	result := h.FilterOpened([]int{443}, proc)
	if len(result) != 1 {
		t.Fatal("expected sshd on 443 to pass through suppression")
	}
}

func TestSuppressHookDefaultWriter(t *testing.T) {
	sup := ports.NewSuppressor()
	h := NewSuppressHook(sup, nil)
	if h.writer == nil {
		t.Fatal("expected default writer to be set")
	}
}
