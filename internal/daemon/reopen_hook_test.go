package daemon

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestReopenHook(window time.Duration) (*ReopenHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	h := NewReopenHook(window, buf)
	return h, buf
}

func TestReopenHookSilentOnOpenOnly(t *testing.T) {
	h, buf := newTestReopenHook(10 * time.Second)

	_ = h.AfterScan(ports.Diff{Opened: []int{8080}})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestReopenHookSilentWhenClosedWithoutReopen(t *testing.T) {
	h, buf := newTestReopenHook(10 * time.Second)

	_ = h.AfterScan(ports.Diff{Closed: []int{8080}})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestReopenHookDetectsReopen(t *testing.T) {
	h, buf := newTestReopenHook(10 * time.Second)

	// Close the port in one cycle.
	_ = h.AfterScan(ports.Diff{Closed: []int{9090}})

	// Reopen in the next cycle within window.
	_ = h.AfterScan(ports.Diff{Opened: []int{9090}})

	if !strings.Contains(buf.String(), "9090") {
		t.Errorf("expected reopen log for port 9090, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "reopened after") {
		t.Errorf("expected 'reopened after' in log, got: %s", buf.String())
	}
}

func TestReopenHookBeforeScanIsNoop(t *testing.T) {
	h, _ := newTestReopenHook(5 * time.Second)
	if err := h.BeforeScan(); err != nil {
		t.Errorf("BeforeScan should return nil, got %v", err)
	}
}

func TestReopenHookDefaultWriterIsStdout(t *testing.T) {
	h := NewReopenHook(5*time.Second, nil)
	if h.logger == nil {
		t.Error("expected non-nil logger")
	}
}
