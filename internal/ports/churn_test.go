package ports

import (
	"testing"
	"time"
)

func fixedChurnClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestChurnNotUnstableBelowThreshold(t *testing.T) {
	now := time.Now()
	ct := NewChurnTracker(time.Minute, 3, fixedChurnClock(now))
	ct.RecordOpen(8080)
	ct.RecordClose(8080)
	// churn rate = 2, threshold = 3 → not unstable
	if got := ct.Unstable(); len(got) != 0 {
		t.Fatalf("expected no unstable ports, got %v", got)
	}
}

func TestChurnUnstableAtThreshold(t *testing.T) {
	now := time.Now()
	ct := NewChurnTracker(time.Minute, 3, fixedChurnClock(now))
	ct.RecordOpen(8080)
	ct.RecordClose(8080)
	ct.RecordOpen(8080)
	// churn rate = 3, threshold = 3 → unstable
	recs := ct.Unstable()
	if len(recs) != 1 {
		t.Fatalf("expected 1 unstable port, got %d", len(recs))
	}
	if recs[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", recs[0].Port)
	}
	if recs[0].Opens != 2 || recs[0].Closes != 1 {
		t.Errorf("unexpected opens/closes: %+v", recs[0])
	}
}

func TestChurnEventsExpireOutsideWindow(t *testing.T) {
	base := time.Now()
	clock := base
	ct := NewChurnTracker(30*time.Second, 2, func() time.Time { return clock })

	// Record two events at base time.
	ct.RecordOpen(443)
	ct.RecordClose(443)

	// Advance clock beyond window.
	clock = base.Add(31 * time.Second)
	// Now add one fresh event — churn rate should be 1, below threshold.
	ct.RecordOpen(443)

	if got := ct.Unstable(); len(got) != 0 {
		t.Fatalf("expected events to have expired, got %v", got)
	}
}

func TestChurnDifferentPortsAreIndependent(t *testing.T) {
	now := time.Now()
	ct := NewChurnTracker(time.Minute, 2, fixedChurnClock(now))
	ct.RecordOpen(80)
	ct.RecordClose(80)
	ct.RecordOpen(443)
	// port 80: rate=2 → unstable; port 443: rate=1 → stable
	recs := ct.Unstable()
	if len(recs) != 1 {
		t.Fatalf("expected 1 unstable port, got %d", len(recs))
	}
	if recs[0].Port != 80 {
		t.Errorf("expected port 80, got %d", recs[0].Port)
	}
}

func TestChurnResetClearsPort(t *testing.T) {
	now := time.Now()
	ct := NewChurnTracker(time.Minute, 2, fixedChurnClock(now))
	ct.RecordOpen(9000)
	ct.RecordClose(9000)
	ct.Reset(9000)
	if got := ct.Unstable(); len(got) != 0 {
		t.Fatalf("expected no unstable ports after reset, got %v", got)
	}
}

func TestChurnRateMethod(t *testing.T) {
	r := ChurnRecord{Opens: 4, Closes: 3}
	if r.ChurnRate() != 7 {
		t.Errorf("expected ChurnRate 7, got %d", r.ChurnRate())
	}
}
