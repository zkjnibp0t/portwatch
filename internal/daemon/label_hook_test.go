package daemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestLabelHook() (*LabelHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	labeler := ports.NewLabeler(map[int]string{
		80:  "http",
		443: "https",
	})
	return NewLabelHook(labeler, buf), buf
}

func TestLabelHookSilentOnNoDiff(t *testing.T) {
	hook, buf := newTestLabelHook()
	hook.OnCycle(ports.Diff{})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestLabelHookLogsOpenedWithLabel(t *testing.T) {
	hook, buf := newTestLabelHook()
	hook.OnCycle(ports.Diff{Opened: []int{80}})
	out := buf.String()
	if !strings.Contains(out, "opened") {
		t.Errorf("expected 'opened' in output, got: %s", out)
	}
	if !strings.Contains(out, "80(http)") {
		t.Errorf("expected '80(http)' in output, got: %s", out)
	}
}

func TestLabelHookLogsClosedWithLabel(t *testing.T) {
	hook, buf := newTestLabelHook()
	hook.OnCycle(ports.Diff{Closed: []int{443}})
	out := buf.String()
	if !strings.Contains(out, "closed") {
		t.Errorf("expected 'closed' in output, got: %s", out)
	}
	if !strings.Contains(out, "443(https)") {
		t.Errorf("expected '443(https)' in output, got: %s", out)
	}
}

func TestLabelHookUnknownPortFallback(t *testing.T) {
	hook, buf := newTestLabelHook()
	hook.OnCycle(ports.Diff{Opened: []int{9999}})
	out := buf.String()
	if !strings.Contains(out, "9999") {
		t.Errorf("expected port 9999 in output, got: %s", out)
	}
}

func TestLabelHookDefaultWriterIsStdout(t *testing.T) {
	labeler := ports.NewLabeler(nil)
	hook := NewLabelHook(labeler, nil)
	if hook == nil {
		t.Fatal("expected non-nil hook")
	}
}
