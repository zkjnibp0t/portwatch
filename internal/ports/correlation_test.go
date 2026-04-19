package ports

import (
	"testing"
	"time"
)

func fixedCorrelationClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCorrelatorEmptyGroups(t *testing.T) {
	c := NewCorrelator(5 * time.Second)
	if len(c.Groups()) != 0 {
		t.Fatal("expected no groups")
	}
}

func TestCorrelatorGroupsOpenedPorts(t *testing.T) {
	now := time.Now()
	c := NewCorrelator(10 * time.Second)
	c.clock = fixedCorrelationClock(now)
	c.Record(80, true)
	c.Record(443, true)
	groups := c.Groups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if !groups[0].Opened {
		t.Error("expected opened group")
	}
	if len(groups[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(groups[0].Ports))
	}
}

func TestCorrelatorSeparatesOpenClose(t *testing.T) {
	now := time.Now()
	c := NewCorrelator(10 * time.Second)
	c.clock = fixedCorrelationClock(now)
	c.Record(80, true)
	c.Record(22, false)
	groups := c.Groups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestCorrelatorEvictsOldEvents(t *testing.T) {
	now := time.Now()
	c := NewCorrelator(5 * time.Second)
	c.clock = fixedCorrelationClock(now.Add(-10 * time.Second))
	c.Record(80, true)
	c.clock = fixedCorrelationClock(now)
	groups := c.Groups()
	if len(groups) != 0 {
		t.Fatal("expected old events to be evicted")
	}
}

func TestCorrelatorFlushClearsEvents(t *testing.T) {
	now := time.Now()
	c := NewCorrelator(10 * time.Second)
	c.clock = fixedCorrelationClock(now)
	c.Record(80, true)
	c.Flush()
	if len(c.Groups()) != 0 {
		t.Fatal("expected groups to be empty after flush")
	}
}

func TestCorrelatorGroupString(t *testing.T) {
	g := CorrelationGroup{Ports: []int{80, 443}, Opened: true, At: time.Now()}
	s := g.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
