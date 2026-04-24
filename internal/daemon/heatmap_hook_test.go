package daemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestHeatmapHook(topN, every int) (*HeatmapHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	tracker := ports.NewHeatmapTracker()
	hook := NewHeatmapHook(tracker, topN, every, buf)
	return hook, buf
}

func TestHeatmapHookSilentOnNoDiff(t *testing.T) {
	hook, buf := newTestHeatmapHook(5, 1)
	hook.AfterScan(ports.Diff{})
	if buf.Len() != 0 {
		t.Fatalf("expected no output on empty diff, got: %s", buf.String())
	}
}

func TestHeatmapHookLogsOnCadence(t *testing.T) {
	hook, buf := newTestHeatmapHook(3, 1)
	hook.AfterScan(ports.Diff{Opened: []int{80}})
	if !strings.Contains(buf.String(), "port=80") {
		t.Fatalf("expected port=80 in output, got: %s", buf.String())
	}
}

func TestHeatmapHookSilentBeforeCadence(t *testing.T) {
	hook, buf := newTestHeatmapHook(3, 3)
	hook.AfterScan(ports.Diff{Opened: []int{80}})
	hook.AfterScan(ports.Diff{Opened: []int{443}})
	if buf.Len() != 0 {
		t.Fatalf("expected no output before cadence, got: %s", buf.String())
	}
	hook.AfterScan(ports.Diff{Opened: []int{8080}})
	if !strings.Contains(buf.String(), "top") {
		t.Fatalf("expected heatmap output on 3rd cycle, got: %s", buf.String())
	}
}

func TestHeatmapHookRanksCorrectly(t *testing.T) {
	hook, buf := newTestHeatmapHook(2, 1)
	hook.AfterScan(ports.Diff{Opened: []int{443}})
	hook.AfterScan(ports.Diff{Opened: []int{443}})
	hook.AfterScan(ports.Diff{Opened: []int{80}})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// last log block: first ranked line should mention 443
	found := false
	for _, l := range lines {
		if strings.Contains(l, "#1") && strings.Contains(l, "port=443") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected port=443 at rank 1, output: %s", buf.String())
	}
}

func TestHeatmapHookDefaultWriterIsStdout(t *testing.T) {
	tracker := ports.NewHeatmapTracker()
	hook := NewHeatmapHook(tracker, 5, 1, nil)
	if hook.logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestHeatmapHookBeforeScanIsNoop(t *testing.T) {
	hook, _ := newTestHeatmapHook(5, 1)
	// should not panic
	hook.BeforeScan()
}
