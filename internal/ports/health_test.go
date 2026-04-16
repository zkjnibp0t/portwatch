package ports

import (
	"errors"
	"testing"
	"time"
)

func fixedHealthClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestHealthTrackerInitialState(t *testing.T) {
	h := NewHealthTracker(nil)
	s := h.Status()
	if s.TotalScans != 0 || s.TotalErrors != 0 {
		t.Errorf("expected zero counts, got scans=%d errors=%d", s.TotalScans, s.TotalErrors)
	}
	if h.IsHealthy() {
		t.Error("expected unhealthy before any scan")
	}
}

func TestHealthTrackerRecordSuccess(t *testing.T) {
	now := time.Now()
	h := NewHealthTracker(fixedHealthClock(now))
	h.RecordSuccess()
	s := h.Status()
	if !s.LastScanOK {
		t.Error("expected LastScanOK true")
	}
	if s.TotalScans != 1 || s.TotalErrors != 0 {
		t.Errorf("unexpected counts: scans=%d errors=%d", s.TotalScans, s.TotalErrors)
	}
	if !s.LastScanAt.Equal(now) {
		t.Errorf("unexpected LastScanAt: %v", s.LastScanAt)
	}
}

func TestHealthTrackerRecordError(t *testing.T) {
	h := NewHealthTracker(nil)
	h.RecordError(errors.New("scan failed"))
	s := h.Status()
	if s.LastScanOK {
		t.Error("expected LastScanOK false after error")
	}
	if s.LastError != "scan failed" {
		t.Errorf("unexpected error string: %s", s.LastError)
	}
	if s.ConsecErrors != 1 || s.TotalErrors != 1 {
		t.Errorf("unexpected error counts: consec=%d total=%d", s.ConsecErrors, s.TotalErrors)
	}
}

func TestHealthTrackerConsecErrorsResetOnSuccess(t *testing.T) {
	h := NewHealthTracker(nil)
	h.RecordError(errors.New("e1"))
	h.RecordError(errors.New("e2"))
	h.RecordSuccess()
	s := h.Status()
	if s.ConsecErrors != 0 {
		t.Errorf("expected ConsecErrors=0 after success, got %d", s.ConsecErrors)
	}
	if s.TotalErrors != 2 {
		t.Errorf("expected TotalErrors=2, got %d", s.TotalErrors)
	}
}

func TestHealthTrackerIsHealthy(t *testing.T) {
	h := NewHealthTracker(nil)
	h.RecordSuccess()
	if !h.IsHealthy() {
		t.Error("expected healthy after success")
	}
	h.RecordError(errors.New("oops"))
	if h.IsHealthy() {
		t.Error("expected unhealthy after error")
	}
}
