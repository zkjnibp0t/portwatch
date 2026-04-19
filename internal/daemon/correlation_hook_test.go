package daemon

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestCorrelationHook(minGroup int) (*CorrelationHook, *bytes.Buffer) {
	var buf bytes.Buffer
	h := NewCorrelationHook(10*time.Second, minGroup, &buf)
	return h, &buf
}

func TestCorrelationHookSilentBelowMinGroup(t *testing.T) {
	h, buf := newTestCorrelationHook(3)
	_ = h.AfterScan(ports.Diff{Opened: []int{80}})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestCorrelationHookLogsWhenMinGroupReached(t *testing.T) {
	h, buf := newTestCorrelationHook(2)
	_ = h.AfterScan(ports.Diff{Opened: []int{80, 443}})
	if !strings.Contains(buf.String(), "correlated group detected") {
		t.Errorf("expected correlation log, got: %s", buf.String())
	}
}

func TestCorrelationHookBeforeScanNoop(t *testing.T) {
	h, _ := newTestCorrelationHook(2)
	if err := h.BeforeScan(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCorrelationHookSeparatesOpenClose(t *testing.T) {
	h, buf := newTestCorrelationHook(1)
	_ = h.AfterScan(ports.Diff{Opened: []int{80}, Closed: []int{22}})
	out := buf.String()
	if strings.Count(out, "correlated group detected") < 2 {
		t.Errorf("expected two group logs, got: %s", out)
	}
}

func TestCorrelationHookDefaultWriterIsStdout(t *testing.T) {
	h := NewCorrelationHook(5*time.Second, 1, nil)
	if h.logger == nil {
		t.Error("expected logger to be set")
	}
}
