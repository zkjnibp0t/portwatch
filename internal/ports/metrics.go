package ports

import (
	"sync"
	"time"
)

// ScanMetrics holds statistics for a single scan cycle.
type ScanMetrics struct {
	Timestamp   time.Time
	Duration    time.Duration
	PortsFound  int
	Opened      int
	Closed      int
	Anomalies   int
}

// MetricsCollector accumulates scan metrics over time.
type MetricsCollector struct {
	mu      sync.Mutex
	records []ScanMetrics
	max     int
}

// NewMetricsCollector creates a collector that retains at most maxRecords entries.
func NewMetricsCollector(maxRecords int) *MetricsCollector {
	if maxRecords <= 0 {
		maxRecords = 100
	}
	return &MetricsCollector{max: maxRecords}
}

// Record appends a new metric entry, evicting the oldest if at capacity.
func (m *MetricsCollector) Record(sm ScanMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.records) >= m.max {
		m.records = m.records[1:]
	}
	m.records = append(m.records, sm)
}

// Latest returns the most recent ScanMetrics, or false if none recorded.
func (m *MetricsCollector) Latest() (ScanMetrics, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.records) == 0 {
		return ScanMetrics{}, false
	}
	return m.records[len(m.records)-1], true
}

// Summary returns aggregate stats across all retained records.
func (m *MetricsCollector) Summary() MetricsSummary {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := MetricsSummary{Count: len(m.records)}
	for _, r := range m.records {
		s.TotalOpened += r.Opened
		s.TotalClosed += r.Closed
		s.TotalAnomalies += r.Anomalies
		s.TotalDuration += r.Duration
	}
	if s.Count > 0 {
		s.AvgDuration = s.TotalDuration / time.Duration(s.Count)
	}
	return s
}

// MetricsSummary holds aggregated values across multiple scan cycles.
type MetricsSummary struct {
	Count          int
	TotalOpened    int
	TotalClosed    int
	TotalAnomalies int
	TotalDuration  time.Duration
	AvgDuration    time.Duration
}
