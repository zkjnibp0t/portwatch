package daemon

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestShadowHook(window time.Duration) (*ShadowHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	h := NewShadowHook(window, logger)
	return h, buf
}

func TestShadowHookSilentOnOpenOnly(t *testing.T) {
	h, buf := newTestShadowHook(5 * time.Second)
	h.AfterScan(ports.Diff{Opened: []int{8080}})
	if buf.Len() != 0 {
		t.Errorf("expected no log on open, got: %s", buf.String())
	}
}

func TestShadowHookSilentWhenClosedWithoutOpen(t *testing.T) {
	h, buf := newTestShadowHook(5 * time.Second)
	h.AfterScan(ports.Diff{Closed: []int{9090}})
	if buf.Len() != 0 {
		t.Errorf("expected no log for unknown close, got: %s", buf.String())
	}
}

func TestShadowHookDetectsFlashPort(t *testing.T) {
	h, buf := newTestShadowHook(10 * time.Second)
	h.tracker.RecordOpened(443)
	// Simulate immediate close within window
	h.AfterScan(ports.Diff{Closed: []int{443}})
	if !strings.Contains(buf.String(), "shadow port detected") {
		t.Errorf("expected shadow log, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "443") {
		t.Errorf("expected port 443 in log, got: %s", buf.String())
	}
}

func TestShadowHookAccumulatesShadows(t *testing.T) {
	h, _ := newTestShadowHook(10 * time.Second)
	for _, p := range []int{80, 22} {
		h.AfterScan(ports.Diff{Opened: []int{p}})
		h.AfterScan(ports.Diff{Closed: []int{p}})
	}
	if len(h.Shadows()) != 2 {
		t.Errorf("expected 2 shadows, got %d", len(h.Shadows()))
	}
}

func TestShadowHookBeforeScanIsNoop(t *testing.T) {
	h, buf := newTestShadowHook(5 * time.Second)
	h.BeforeScan()
	if buf.Len() != 0 {
		t.Error("expected no log from BeforeScan")
	}
}
