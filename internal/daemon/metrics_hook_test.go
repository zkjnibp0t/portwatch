package daemon

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func newTestMetricsHook() (*MetricsHook, *ports.MetricsCollector, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	c := ports.NewMetricsCollector(50)
	h := NewMetricsHook(c, logger)
	return h, c, buf
}

func TestMetricsHookRecords(t *testing.T) {
	h, c, _ := newTestMetricsHook()
	diff := ports.Diff{Opened: []int{8080, 9090}, Closed: []int{22}}
	h.Record(15, diff, 1, 5*time.Millisecond)

	last, ok := c.Latest()
	if !ok {
		t.Fatal("expected a record")
	}
	if last.PortsFound != 15 {
		t.Errorf("expected PortsFound=15, got %d", last.PortsFound)
	}
	if last.Opened != 2 {
		t.Errorf("expected Opened=2, got %d", last.Opened)
	}
	if last.Closed != 1 {
		t.Errorf("expected Closed=1, got %d", last.Closed)
	}
	if last.Anomalies != 1 {
		t.Errorf("expected Anomalies=1, got %d", last.Anomalies)
	}
}

func TestMetricsHookLogsOutput(t *testing.T) {
	h, _, buf := newTestMetricsHook()
	diff := ports.Diff{Opened: []int{443}}
	h.Record(5, diff, 0, 3*time.Millisecond)

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected log output, got none")
	}
}

func TestMetricsHookDefaultLogger(t *testing.T) {
	c := ports.NewMetricsCollector(10)
	h := NewMetricsHook(c, nil)
	if h.logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestMetricsHookSummaryAccumulates(t *testing.T) {
	h, c, _ := newTestMetricsHook()
	h.Record(10, ports.Diff{Opened: []int{80}}, 0, time.Millisecond)
	h.Record(10, ports.Diff{Opened: []int{443}, Closed: []int{80}}, 0, time.Millisecond)

	s := c.Summary()
	if s.Count != 2 {
		t.Errorf("expected Count=2, got %d", s.Count)
	}
	if s.TotalOpened != 2 {
		t.Errorf("expected TotalOpened=2, got %d", s.TotalOpened)
	}
	if s.TotalClosed != 1 {
		t.Errorf("expected TotalClosed=1, got %d", s.TotalClosed)
	}
}
