package ports_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func startWatchListener(t *testing.T)int, func()) {
	t.Helper()
	lntcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestWatcherEmitsEvent(t *testing.T) {
	port, stop := startWatchListener(t)
	defer stop()

	scanner := ports.NewScanner(port, port)
	filter := ports.NewFilter(nil, nil)
	watcher := ports.NewWatcher(scanner, filter, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch := watcher.Run(ctx)
	var got ports.WatchEvent
	for ev := range ch {
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		got = ev
		break
	}
	if _, ok := got.Current[port]; !ok {
		t.Errorf("expected port %d in current set, got %v", port, got.Current)
	}
}

func TestWatcherDetectsClose(t *testing.T) {
	port, stop := startWatchListener(t)

	scanner := ports.NewScanner(port, port)
	filter := ports.NewFilter(nil, nil)
	watcher := ports.NewWatcher(scanner, filter, 60*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	ch := watcher.Run(ctx)

	// Consume first event (port open)
	<-ch
	// Close the listener so the port disappears
	stop()

	// Wait for an event that detects the closure
	for ev := range ch {
		if ev.Err != nil {
			continue
		}
		if len(ev.Diff.Closed) > 0 {
			return // success
		}
	}
	t.Error("expected a closed-port diff event")
}

func TestWatcherCancelsCleanly(t *testing.T) {
	scanner := ports.NewScanner(19900, 19910)
	filter := ports.NewFilter(nil, nil)
	watcher := ports.NewWatcher(scanner, filter, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	ch := watcher.Run(ctx)
	cancel()

	// Drain channel; it must close without blocking.
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("watcher channel did not close after cancel")
		}
	}
}
