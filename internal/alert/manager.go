package alert

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/portwatch/internal/notify"
	"github.com/yourusername/portwatch/internal/ports"
)

// Manager coordinates throttling and dispatching of alerts
// when port changes are detected.
type Manager struct {
	notifier  notify.Notifier
	throttler *Throttler
}

// NewManager creates an alert Manager with the given notifier and
// cooldown duration used to throttle repeated alerts.
func NewManager(n notify.Notifier, cooldown time.Duration) *Manager {
	return &Manager{
		notifier:  n,
		throttler: NewThrottler(cooldown),
	}
}

// Dispatch sends an alert for the given diff if throttling allows it.
// The throttle key is derived from the set of changed ports so that
// identical back-to-back diffs are suppressed.
func (m *Manager) Dispatch(ctx context.Context, diff ports.Diff) error {
	if diff.IsEmpty() {
		return nil
	}

	key := diffKey(diff)
	if !m.throttler.Allow(key) {
		log.Printf("[alert] throttled: %s", key)
		return nil
	}

	if err := m.notifier.Notify(ctx, diff); err != nil {
		// Reset throttle so the next tick can retry.
		m.throttler.Reset(key)
		return fmt.Errorf("alert dispatch failed: %w", err)
	}
	return nil
}

// diffKey produces a stable string key representing the diff contents.
func diffKey(d ports.Diff) string {
	return fmt.Sprintf("opened=%v closed=%v", d.Opened, d.Closed)
}
