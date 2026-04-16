package ports

import (
	"testing"
	"time"
)

func baseMetric(opened, closed, anomalies int, dur time.Duration) ScanMetrics {
	return ScanMetrics{
		Timestamp: time.Now(),
		Duration:  dur,
		PortsFound: 10,
		Opened:    opened,
		Closed:    closed,
		Anomalies: anomalies,
	}
}

func TestMetricsLatestEmpty(t *testing.T) {
	c := NewMetricsCollector(10)
	_, ok := c.Latest()
	if ok {
		t.Fatal("expected no latest on empty collector")
	}
}

func TestMetricsLatestReturnsLast(t *testing.T) {
	c := NewMetricsCollector(10)
	c.Record(baseMetric(1, 0, 0, time.Millisecond))
	c.Record(baseMetric(2, 1, 0, 2*time.Millisecond))
	last, ok := c.Latest()
	if !ok {
		t.Fatal("expected a record")
	}
	if last.Opened != 2 {
		t.Errorf("expected Opened=2, got %d", last.Opened)
	}
}

func TestMetricsEvictsOldest(t *testing.T) {
	c := NewMetricsCollector(3)
	for i := 0; i < 5; i++ {
		c.Record(baseMetric(i, 0, 0, time.Millisecond))
	}
	if len(c.records) != 3 {
		t.Errorf("expected 3 records, got %d", len(c.records))
	}
	if c.records[0].Opened != 2 {
		t.Errorf("expected oldest Opened=2, got %d", c.records[0].Opened)
	}
}

func TestMetricsSummaryAggregates(t *testing.T) {
	c := NewMetricsCollector(10)
	c.Record(baseMetric(3, 1, 2, 10*time.Millisecond))
	c.Record(baseMetric(1, 2, 0, 20*time.Millisecond))
	s := c.Summary()
	if s.Count != 2 {
		t.Errorf("expected Count=2, got %d", s.Count)
	}
	if s.TotalOpened != 4 {
		t.Errorf("expected TotalOpened=4, got %d", s.TotalOpened)
	}
	if s.TotalClosed != 3 {
		t.Errorf("expected TotalClosed=3, got %d", s.TotalClosed)
	}
	if s.TotalAnomalies != 2 {
		t.Errorf("expected TotalAnomalies=2, got %d", s.TotalAnomalies)
	}
	if s.AvgDuration != 15*time.Millisecond {
		t.Errorf("expected AvgDuration=15ms, got %v", s.AvgDuration)
	}
}

func TestMetricsSummaryEmpty(t *testing.T) {
	c := NewMetricsCollector(10)
	s := c.Summary()
	if s.Count != 0 || s.AvgDuration != 0 {
		t.Error("expected zero summary for empty collector")
	}
}
