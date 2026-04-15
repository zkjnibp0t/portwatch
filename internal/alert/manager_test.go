package alert_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/portwatch/internal/alert"
	"github.com/yourusername/portwatch/internal/ports"
)

type mockNotifier struct {
	calls int
	err   error
}

func (m *mockNotifier) Notify(_ context.Context, _ ports.Diff) error {
	m.calls++
	return m.err
}

func TestDispatchEmptyDiffSkipped(t *testing.T) {
	n := &mockNotifier{}
	mgr := alert.NewManager(n, 5*time.Second)
	err := mgr.Dispatch(context.Background(), ports.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.calls != 0 {
		t.Errorf("expected 0 notifier calls, got %d", n.calls)
	}
}

func TestDispatchSendsAlert(t *testing.T) {
	n := &mockNotifier{}
	mgr := alert.NewManager(n, 5*time.Second)
	diff := ports.Diff{Opened: []int{8080}}
	if err := mgr.Dispatch(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.calls != 1 {
		t.Errorf("expected 1 notifier call, got %d", n.calls)
	}
}

func TestDispatchThrottlesSameDiff(t *testing.T) {
	n := &mockNotifier{}
	mgr := alert.NewManager(n, 1*time.Hour)
	diff := ports.Diff{Opened: []int{8080}}
	mgr.Dispatch(context.Background(), diff) //nolint
	mgr.Dispatch(context.Background(), diff) //nolint
	if n.calls != 1 {
		t.Errorf("expected 1 notifier call due to throttle, got %d", n.calls)
	}
}

func TestDispatchResetsThrottleOnError(t *testing.T) {
	n := &mockNotifier{err: errors.New("send failed")}
	mgr := alert.NewManager(n, 1*time.Hour)
	diff := ports.Diff{Opened: []int{9090}}
	mgr.Dispatch(context.Background(), diff) //nolint: first call fails
	// After failure the throttle should be reset, allowing a retry.
	n.err = nil
	if err := mgr.Dispatch(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if n.calls != 2 {
		t.Errorf("expected 2 notifier calls (1 fail + 1 retry), got %d", n.calls)
	}
}
