package daemon

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestHook(threshold int) (*TrendHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	tracker := ports.NewTrendTracker(time.Hour)
	return NewTrendHook(tracker, threshold, logger), buf
}

func TestTrendHookRecordsOpenedPorts(t *testing.T) {
	hook, _ := newTestHook(10)
	hook.Apply(ports.Diff{Opened: []int{8080, 9090}})

	if got := hook.tracker.FlapCount(8080); got != 1 {
		t.Fatalf("expected 1 for 8080, got %d", got)
	}
	if got := hook.tracker.FlapCount(9090); got != 1 {
		t.Fatalf("expected 1 for 9090, got %d", got)
	}
}

func TestTrendHookRecordsClosedPorts(t *testing.T) {
	hook, _ := newTestHook(10)
	hook.Apply(ports.Diff{Closed: []int{443}})

	if got := hook.tracker.FlapCount(443); got != 1 {
		t.Fatalf("expected 1 for 443, got %d", got)
	}
}

func TestTrendHookLogsWhenThresholdReached(t *testing.T) {
	hook, buf := newTestHook(3)
	for i := 0; i < 3; i++ {
		hook.Apply(ports.Diff{Opened: []int{7000}})
	}

	if buf.Len() == 0 {
		t.Fatal("expected warning log, got nothing")
	}
}

func TestTrendHookNoLogBelowThreshold(t *testing.T) {
	hook, buf := newTestHook(5)
	hook.Apply(ports.Diff{Opened: []int{7000}})

	if buf.Len() != 0 {
		t.Fatalf("expected no log, got: %s", buf.String())
	}
}

func TestTrendHookDefaultThreshold(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	tracker := ports.NewTrendTracker(time.Hour)
	hook := NewTrendHook(tracker, 0, logger)
	if hook.threshold != 5 {
		t.Fatalf("expected default threshold 5, got %d", hook.threshold)
	}
}
