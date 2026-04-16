package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// MetricsHook wires a MetricsCollector into the daemon cycle.
type MetricsHook struct {
	collector *ports.MetricsCollector
	logger    *log.Logger
}

// NewMetricsHook creates a MetricsHook backed by the given collector.
func NewMetricsHook(c *ports.MetricsCollector, logger *log.Logger) *MetricsHook {
	if logger == nil {
		logger = log.Default()
	}
	return &MetricsHook{collector: c, logger: logger}
}

// Record builds a ScanMetrics entry and stores it in the collector.
func (h *MetricsHook) Record(portsFound int, diff ports.Diff, anomalies int, dur time.Duration) {
	sm := ports.ScanMetrics{
		Timestamp:  time.Now(),
		Duration:   dur,
		PortsFound: portsFound,
		Opened:     len(diff.Opened),
		Closed:     len(diff.Closed),
		Anomalies:  anomalies,
	}
	h.collector.Record(sm)

	s := h.collector.Summary()
	h.logger.Printf("[metrics] scan=%v ports=%d opened=%d closed=%d anomalies=%d | totals: scans=%d opened=%d closed=%d",
		dur.Round(time.Millisecond), portsFound,
		sm.Opened, sm.Closed, sm.Anomalies,
		s.Count, s.TotalOpened, s.TotalClosed,
	)
}
