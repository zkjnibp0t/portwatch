package ports

import (
	"context"
	"time"
)

// WatchEvent holds the result of a single scan cycle.
type WatchEvent struct {
	Current PortSet
	Previous PortSet
	Diff     Diff
	Err      error
}

// Watcher periodically scans ports and emits WatchEvents on a channel.
type Watcher struct {
	scanner  *Scanner
	filter   *Filter
	interval time.Duration
	previous PortSet
}

// NewWatcher creates a Watcher using the provided Scanner, Filter, and poll interval.
func NewWatcher(scanner *Scanner, filter *Filter, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  scanner,
		filter:   filter,
		interval: interval,
		previous: PortSet{},
	}
}

// Run starts the watch loop and sends events to the returned channel.
// The channel is closed when ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.tick(ctx, ch)
			}
		}
	}()
	return ch
}

func (w *Watcher) tick(ctx context.Context, ch chan<- WatchEvent) {
	raw, err := w.scanner.Scan()
	if err != nil {
		select {
		case ch <- WatchEvent{Err: err}:
		case <-ctx.Done():
		}
		return
	}
	current := w.filter.Apply(ToSet(raw))
	diff := Compare(w.previous, current)
	event := WatchEvent{
		Current:  current,
		Previous: w.previous,
		Diff:     diff,
	}
	w.previous = current
	select {
	case ch <- event:
	case <-ctx.Done():
	}
}
