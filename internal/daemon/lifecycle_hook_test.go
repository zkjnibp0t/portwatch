package daemon

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestLifecycleHook(minLog time.Duration, buf *bytes.Buffer) (*LifecycleHook, *ports.LifecycleTracker) {
	var tick int64
	base := time.Unix(1000, 0)
	tracker := ports.NewLifecycleTracker(func() time.Time {
		t := base.Add(time.Duration(tick) * time.Second)
		tick++
		return t
	})
	hook := NewLifecycleHook(tracker, minLog, buf)
	return hook, tracker
}

func TestLifecycleHookBeforeScanNoop(t *testing.T) {
	hook, _ := newTestLifecycleHook(0, &bytes.Buffer{})
	if err := hook.BeforeScan(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLifecycleHookRecordsOpened(t *testing.T) {
	hook, tracker := newTestLifecycleHook(0, &bytes.Buffer{})
	diff := ports.Diff{Opened: []int{8080}}
	_ = hook.AfterScan(diff)

	_, ok := tracker.OpenSince(8080)
	if !ok {
		t.Error("expected port 8080 to be tracked as open")
	}
}

func TestLifecycleHookRecordsClosed(t *testing.T) {
	hook, tracker := newTestLifecycleHook(0, &bytes.Buffer{})
	_ = hook.AfterScan(ports.Diff{Opened: []int{9090}})
	_ = hook.AfterScan(ports.Diff{Closed: []int{9090}})

	_, ok := tracker.OpenSince(9090)
	if ok {
		t.Error("expected port 9090 to be closed")
	}
}

func TestLifecycleHookLogsLongLivedPort(t *testing.T) {
	var buf bytes.Buffer
	base := time.Unix(0, 0)
	var tick int64
	tracker := ports.NewLifecycleTracker(func() time.Time {
		t := base.Add(time.Duration(tick) * time.Minute)
		tick += 10
		return t
	})
	hook := NewLifecycleHook(tracker, 5*time.Minute, &buf)

	_ = hook.AfterScan(ports.Diff{Opened: []int{443}})
	_ = hook.AfterScan(ports.Diff{Closed: []int{443}})

	if !strings.Contains(buf.String(), "443") {
		t.Errorf("expected log to mention port 443, got: %q", buf.String())
	}
}

func TestLifecycleHookSilentBelowMinLog(t *testing.T) {
	var buf bytes.Buffer
	base := time.Unix(0, 0)
	var tick int64
	tracker := ports.NewLifecycleTracker(func() time.Time {
		t := base.Add(time.Duration(tick) * time.Second)
		tick++
		return t
	})
	hook := NewLifecycleHook(tracker, 10*time.Minute, &buf)

	_ = hook.AfterScan(ports.Diff{Opened: []int{80}})
	_ = hook.AfterScan(ports.Diff{Closed: []int{80}})

	if buf.Len() != 0 {
		t.Errorf("expected no log output, got: %q", buf.String())
	}
}
