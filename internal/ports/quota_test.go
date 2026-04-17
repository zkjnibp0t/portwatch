package ports

import (
	"testing"
	"time"
)

func fixedQuotaClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestQuotaAllowsUpToLimit(t *testing.T) {
	q := NewQuotaTracker(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := q.Allow(); err != nil {
			t.Fatalf("expected allow on call %d, got %v", i+1, err)
		}
	}
}

func TestQuotaBlocksOverLimit(t *testing.T) {
	q := NewQuotaTracker(2, time.Minute)
	_ = q.Allow()
	_ = q.Allow()
	err := q.Allow()
	if err == nil {
		t.Fatal("expected quota exceeded error")
	}
	qe, ok := err.(*QuotaExceededError)
	if !ok {
		t.Fatalf("expected QuotaExceededError, got %T", err)
	}
	if qe.Limit != 2 {
		t.Errorf("expected limit 2, got %d", qe.Limit)
	}
}

func TestQuotaResetsAfterWindow(t *testing.T) {
	now := time.Now()
	q := NewQuotaTracker(2, time.Minute)
	q.clock = fixedQuotaClock(now)
	_ = q.Allow()
	_ = q.Allow()

	// advance past window
	q.clock = fixedQuotaClock(now.Add(2 * time.Minute))
	if err := q.Allow(); err != nil {
		t.Fatalf("expected allow after window reset, got %v", err)
	}
}

func TestQuotaRemaining(t *testing.T) {
	q := NewQuotaTracker(5, time.Minute)
	_ = q.Allow()
	_ = q.Allow()
	if got := q.Remaining(); got != 3 {
		t.Errorf("expected 3 remaining, got %d", got)
	}
}

func TestQuotaReset(t *testing.T) {
	q := NewQuotaTracker(2, time.Minute)
	_ = q.Allow()
	_ = q.Allow()
	q.Reset()
	if got := q.Remaining(); got != 2 {
		t.Errorf("expected full quota after reset, got %d", got)
	}
}

func TestQuotaErrorMessage(t *testing.T) {
	e := &QuotaExceededError{Limit: 10, Window: time.Minute, Consumed: 10}
	if e.Error() == "" {
		t.Error("expected non-empty error message")
	}
}
