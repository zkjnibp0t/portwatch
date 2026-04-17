package daemon

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestQuotaHook(limit int) (*QuotaHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	q := ports.NewQuotaTracker(limit, time.Minute)
	return NewQuotaHook(q, buf), buf
}

func TestQuotaHookAllowsInitially(t *testing.T) {
	h, _ := newTestQuotaHook(5)
	if !h.BeforeScan() {
		t.Error("expected scan to be allowed initially")
	}
}

func TestQuotaHookBlocksAfterLimit(t *testing.T) {
	h, buf := newTestQuotaHook(2)
	h.BeforeScan()
	h.BeforeScan()
	if h.BeforeScan() {
		t.Error("expected scan to be blocked after limit")
	}
	if buf.Len() == 0 {
		t.Error("expected log output when blocked")
	}
}

func TestQuotaHookLogsRemaining(t *testing.T) {
	h, buf := newTestQuotaHook(1)
	h.BeforeScan()
	h.BeforeScan() // triggers log
	if buf.Len() == 0 {
		t.Error("expected log message")
	}
	got := buf.String()
	if len(got) == 0 {
		t.Error("log message should not be empty")
	}
}

func TestQuotaHookAfterScanIsNoop(t *testing.T) {
	h, _ := newTestQuotaHook(5)
	// should not panic
	h.AfterScan()
}

func TestQuotaHookDefaultWriter(t *testing.T) {
	q := ports.NewQuotaTracker(3, time.Minute)
	h := NewQuotaHook(q, nil)
	if h.logger == nil {
		t.Error("expected default logger to be set")
	}
}
