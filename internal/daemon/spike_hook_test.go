package daemon

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestSpikeHook(window time.Duration) (*SpikeHook, *bytes.Buffer) {
	var buf bytes.Buffer
	h := NewSpikeHook(window, &buf)
	return h, &buf
}

func TestSpikeHookSilentOnOpenOnly(t *testing.T) {
	h, buf := newTestSpikeHook(5 * time.Second)
	h.AfterScan(ports.Diff{Opened: []int{8080}})
	if buf.Len() != 0 {
		t.Fatalf("expected no output on open-only diff, got: %s", buf.String())
	}
}

func TestSpikeHookSilentWhenClosedWithoutOpen(t *testing.T) {
	h, buf := newTestSpikeHook(5 * time.Second)
	h.AfterScan(ports.Diff{Closed: []int{9090}})
	if buf.Len() != 0 {
		t.Fatalf("expected no output when closed port was never opened")
	}
}

func TestSpikeHookDetectsSpike(t *testing.T) {
	h, buf := newTestSpikeHook(10 * time.Second)
	h.AfterScan(ports.Diff{Opened: []int{1234}})
	// override clock so close happens within window
	h.detector = ports.NewSpikeDetector(10 * time.Second)
	h.AfterScan(ports.Diff{Opened: []int{1234}})
	h.AfterScan(ports.Diff{Closed: []int{1234}})
	if !strings.Contains(buf.String(), "spike detected") {
		t.Fatalf("expected spike log, got: %q", buf.String())
	}
}

func TestSpikeHookLogsPort(t *testing.T) {
	h, buf := newTestSpikeHook(1 * time.Hour)
	h.AfterScan(ports.Diff{Opened: []int{443}})
	h.AfterScan(ports.Diff{Closed: []int{443}})
	out := buf.String()
	if !strings.Contains(out, "443") {
		t.Fatalf("expected port 443 in log, got: %q", out)
	}
}

func TestSpikeHookBeforeScanNoop(t *testing.T) {
	h, buf := newTestSpikeHook(5 * time.Second)
	h.BeforeScan()
	if buf.Len() != 0 {
		t.Fatal("BeforeScan should produce no output")
	}
}
