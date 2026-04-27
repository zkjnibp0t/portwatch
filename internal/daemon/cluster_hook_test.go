package daemon

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func newTestClusterHook(minSupport int) (*ClusterHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	d := ports.NewClusterDetector(minSupport)
	h := NewClusterHook(d, minSupport, buf)
	return h, buf
}

func TestClusterHookSilentBelowSupport(t *testing.T) {
	h, buf := newTestClusterHook(3)
	h.AfterScan(ports.Diff{Opened: []int{80, 443}})
	h.AfterScan(ports.Diff{Opened: []int{80, 443}})
	if buf.Len() != 0 {
		t.Errorf("expected no output below support, got: %s", buf.String())
	}
}

func TestClusterHookLogsWhenSupportReached(t *testing.T) {
	h, buf := newTestClusterHook(2)
	h.AfterScan(ports.Diff{Opened: []int{80, 443}})
	h.AfterScan(ports.Diff{Opened: []int{80, 443}})
	if buf.Len() == 0 {
		t.Error("expected cluster log output, got none")
	}
	if got := buf.String(); !containsStr(got, "cluster detected") {
		t.Errorf("expected 'cluster detected' in output, got: %s", got)
	}
}

func TestClusterHookSilentOnSinglePortChange(t *testing.T) {
	h, buf := newTestClusterHook(1)
	h.AfterScan(ports.Diff{Opened: []int{80}})
	if buf.Len() != 0 {
		t.Errorf("expected no output for single port, got: %s", buf.String())
	}
}

func TestClusterHookBeforeScanNoop(t *testing.T) {
	h, buf := newTestClusterHook(1)
	h.BeforeScan()
	if buf.Len() != 0 {
		t.Errorf("expected no output from BeforeScan, got: %s", buf.String())
	}
}

func TestClusterHookMixedOpenClose(t *testing.T) {
	h, buf := newTestClusterHook(2)
	h.AfterScan(ports.Diff{Opened: []int{8080}, Closed: []int{9090}})
	h.AfterScan(ports.Diff{Opened: []int{8080}, Closed: []int{9090}})
	if buf.Len() == 0 {
		t.Error("expected cluster log for mixed open/close, got none")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
