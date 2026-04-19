package daemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestSequenceHook() (*SequenceHook, *bytes.Buffer, *ports.SequenceTracker) {
	buf := &bytes.Buffer{}
	tracker := ports.NewSequenceTracker()
	h := NewSequenceHook(tracker, buf)
	return h, buf, tracker
}

func TestSequenceHookSilentOnNoDiff(t *testing.T) {
	h, buf, _ := newTestSequenceHook()
	h.AfterScan(ports.Diff{})
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestSequenceHookLogsOpened(t *testing.T) {
	h, buf, _ := newTestSequenceHook()
	h.AfterScan(ports.Diff{Opened: []int{80}})
	out := buf.String()
	if !strings.Contains(out, "opened") || !strings.Contains(out, "80") {
		t.Fatalf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "seq=1") {
		t.Fatalf("expected seq=1 in output: %s", out)
	}
}

func TestSequenceHookLogsClosed(t *testing.T) {
	h, buf, _ := newTestSequenceHook()
	h.AfterScan(ports.Diff{Closed: []int{443}})
	out := buf.String()
	if !strings.Contains(out, "closed") || !strings.Contains(out, "443") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestSequenceHookIncrementsAcrossEvents(t *testing.T) {
	h, _, tracker := newTestSequenceHook()
	h.AfterScan(ports.Diff{Opened: []int{22}, Closed: []int{80}})
	if tracker.LastFor(22) == tracker.LastFor(80) {
		t.Fatal("expected different sequence numbers for different ports")
	}
	if tracker.Len() != 2 {
		t.Fatalf("expected 2 tracked ports, got %d", tracker.Len())
	}
}

func TestSequenceHookBeforeScanNoop(t *testing.T) {
	h, buf, _ := newTestSequenceHook()
	h.BeforeScan()
	if buf.Len() != 0 {
		t.Fatal("BeforeScan should produce no output")
	}
}
