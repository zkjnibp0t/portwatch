package daemon

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/ports"
)

func newTestReportHook(t *testing.T) (*ReportHook, *bytes.Buffer, *history.Store) {
	t.Helper()
	dir := t.TempDir()
	store, err := history.NewStore(filepath.Join(dir, "history.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	var buf bytes.Buffer
	hook := NewReportHook(store, &buf)
	return hook, &buf, store
}

func TestReportHookSilentOnNoDiff(t *testing.T) {
	hook, buf, _ := newTestReportHook(t)
	hook.OnCycle(ports.Diff{})
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestReportHookLogsOpenedPorts(t *testing.T) {
	hook, buf, _ := newTestReportHook(t)
	hook.OnCycle(ports.Diff{Opened: []int{8080, 9090}})
	out := buf.String()
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("expected opened ports in output, got: %s", out)
	}
}

func TestReportHookLogsClosedPorts(t *testing.T) {
	hook, buf, _ := newTestReportHook(t)
	hook.OnCycle(ports.Diff{Closed: []int{443}})
	out := buf.String()
	if !strings.Contains(out, "443") {
		t.Errorf("expected closed port in output, got: %s", out)
	}
}

func TestReportHookDefaultWriterIsStdout(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(filepath.Join(dir, "h.json"))
	hook := NewReportHook(store, nil)
	if hook.out != os.Stdout {
		t.Error("expected default writer to be os.Stdout")
	}
}

func TestFormatPortsEmpty(t *testing.T) {
	if got := formatPorts(nil); got != "-" {
		t.Errorf("expected dash for empty, got %q", got)
	}
}

func TestFormatPortsMultiple(t *testing.T) {
	got := formatPorts([]int{22, 80, 443})
	if !strings.Contains(got, "22") || !strings.Contains(got, "80") || !strings.Contains(got, "443") {
		t.Errorf("unexpected format: %s", got)
	}
}
