package daemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestRollupHook(every int) (*RollupHook, *bytes.Buffer) {
	var buf bytes.Buffer
	return NewRollupHook(every, &buf), &buf
}

func TestRollupHookSilentBeforeWindow(t *testing.T) {
	h, buf := newTestRollupHook(3)
	h.AfterScan(ports.Diff{Opened: []int{80}})
	h.AfterScan(ports.Diff{Opened: []int{443}})
	if buf.Len() != 0 {
		t.Fatal("expected no output before window completes")
	}
}

func TestRollupHookLogsOnWindow(t *testing.T) {
	h, buf := newTestRollupHook(2)
	h.AfterScan(ports.Diff{Opened: []int{8080}})
	h.AfterScan(ports.Diff{Closed: []int{8080}})
	if !strings.Contains(buf.String(), "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", buf.String())
	}
}

func TestRollupHookResetsAfterWindow(t *testing.T) {
	h, buf := newTestRollupHook(1)
	h.AfterScan(ports.Diff{Opened: []int{22}})
	buf.Reset()
	h.AfterScan(ports.Diff{})
	if !strings.Contains(buf.String(), "no port activity") {
		t.Errorf("expected no-activity message, got: %s", buf.String())
	}
}

func TestRollupHookAggregatesNet(t *testing.T) {
	h, buf := newTestRollupHook(1)
	h.AfterScan(ports.Diff{Opened: []int{9000}, Closed: []int{9000}})
	if !strings.Contains(buf.String(), "net=0") {
		t.Errorf("expected net=0, got: %s", buf.String())
	}
}

func TestRollupHookDefaultEveryClampedToOne(t *testing.T) {
	h, buf := newTestRollupHook(0)
	h.AfterScan(ports.Diff{Opened: []int{80}})
	if buf.Len() == 0 {
		t.Fatal("expected output on first cycle when every=0 clamped to 1")
	}
}
