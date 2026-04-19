package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
)

// mockNotifier records every Diff it receives.
type mockNotifier struct {
	calls []ports.Diff
}

func (m *mockNotifier) Notify(diff ports.Diff) error {
	m.calls = append(m.calls, diff)
	return nil
}

func baseCfg() *config.Config {
	return &config.Config{
		PortRange: config.PortRange{Start: 10000, End: 10010},
		Interval:  50 * time.Millisecond,
		Desktop:   config.DesktopConfig{Enabled: false},
	}
}

func TestNewDaemonNoNotifiers(t *testing.T) {
	cfg := baseCfg()
	d, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.notifiers) != 0 {
		t.Errorf("expected 0 notifiers, got %d", len(d.notifiers))
	}
}

func TestNewDaemonWithWebhook(t *testing.T) {
	cfg := baseCfg()
	cfg.Webhook.URL = "http://example.com/hook"

	d, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.notifiers) != 1 {
		t.Errorf("expected 1 notifier, got %d", len(d.notifiers))
	}
}

func TestRunCancelsCleanly(t *testing.T) {
	cfg := baseCfg()
	d, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestAlertInvokesAllNotifiers(t *testing.T) {
	cfg := baseCfg()
	d, _ := New(cfg)

	m1 := &mockNotifier{}
	m2 := &mockNotifier{}
	d.notifiers = []Notifier{m1, m2}

	diff := ports.Diff{Opened: []int{8080}, Closed: []int{}}
	d.alert(diff)

	if len(m1.calls) != 1 {
		t.Errorf("m1: expected 1 call, got %d", len(m1.calls))
	}
	if len(m2.calls) != 1 {
		t.Errorf("m2: expected 1 call, got %d", len(m2.calls))
	}
}

func TestAlertNoNotifiers(t *testing.T) {
	cfg := baseCfg()
	d, _ := New(cfg)

	// alert with no notifiers registered should not panic
	diff := ports.Diff{Opened: []int{9090}, Closed: []int{}}
	d.alert(diff)
}
