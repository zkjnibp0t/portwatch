package daemon

import (
	"bytes"
	"log"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestFingerprintHook() (*FingerprintHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	fp := ports.NewFingerprinter()
	hook := NewFingerprintHook(fp, logger)
	return hook, buf
}

func TestFingerprintHookSilentOnFirstScan(t *testing.T) {
	hook, buf := newTestFingerprintHook()

	set := ports.Set{80: {}, 443: {}}
	hook.AfterScan(set, ports.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output on first scan, got: %s", buf.String())
	}
}

func TestFingerprintHookSilentWhenUnchanged(t *testing.T) {
	hook, buf := newTestFingerprintHook()

	set := ports.Set{80: {}, 443: {}}
	hook.AfterScan(set, ports.Diff{})
	buf.Reset()

	// Same set again — fingerprint should not change
	hook.AfterScan(set, ports.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output when fingerprint unchanged, got: %s", buf.String())
	}
}

func TestFingerprintHookLogsOnChange(t *testing.T) {
	hook, buf := newTestFingerprintHook()

	set1 := ports.Set{80: {}, 443: {}}
	hook.AfterScan(set1, ports.Diff{})
	buf.Reset()

	set2 := ports.Set{80: {}, 8080: {}}
	hook.AfterScan(set2, ports.Diff{Opened: []int{8080}, Closed: []int{443}})

	if buf.Len() == 0 {
		t.Error("expected log output when fingerprint changes")
	}
	out := buf.String()
	if !containsStr(out, "fingerprint") {
		t.Errorf("expected 'fingerprint' in output, got: %s", out)
	}
}

func TestFingerprintHookLogsOldAndNew(t *testing.T) {
	hook, buf := newTestFingerprintHook()

	set1 := ports.Set{22: {}}
	hook.AfterScan(set1, ports.Diff{})
	buf.Reset()

	set2 := ports.Set{22: {}, 3306: {}}
	hook.AfterScan(set2, ports.Diff{Opened: []int{3306}})

	out := buf.String()
	if !containsStr(out, "old=") || !containsStr(out, "new=") {
		t.Errorf("expected old= and new= in output, got: %s", out)
	}
}

func TestFingerprintHookDefaultLogger(t *testing.T) {
	fp := ports.NewFingerprinter()
	hook := NewFingerprintHook(fp, nil)
	if hook == nil {
		t.Fatal("expected non-nil hook with nil logger")
	}
	// Should not panic
	hook.AfterScan(ports.Set{80: {}}, ports.Diff{})
	hook.AfterScan(ports.Set{443: {}}, ports.Diff{Opened: []int{443}, Closed: []int{80}})
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
