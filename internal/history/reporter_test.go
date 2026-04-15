package history_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestPrintReportEmpty(t *testing.T) {
	var buf bytes.Buffer
	history.PrintReport(&buf, nil)
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrintReportContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	entries := []history.Entry{
		{Timestamp: time.Now(), Opened: []string{"8080"}, Closed: []string{"22"}},
	}
	history.PrintReport(&buf, entries)
	out := buf.String()
	for _, h := range []string{"TIMESTAMP", "OPENED", "CLOSED"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q in output", h)
		}
	}
}

func TestPrintReportDashForEmpty(t *testing.T) {
	var buf bytes.Buffer
	entries := []history.Entry{
		{Timestamp: time.Now(), Opened: []string{"443"}, Closed: nil},
	}
	history.PrintReport(&buf, entries)
	if !strings.Contains(buf.String(), "-") {
		t.Errorf("expected dash for empty closed ports")
	}
}

func TestPrintReportMultiplePorts(t *testing.T) {
	var buf bytes.Buffer
	entries := []history.Entry{
		{Timestamp: time.Now(), Opened: []string{"80", "443", "8080"}, Closed: nil},
	}
	history.PrintReport(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "80,443,8080") {
		t.Errorf("expected comma-separated ports, got: %s", out)
	}
}
