package ports

import (
	"sync"
	"time"
)

// Debouncer suppresses rapid repeated port change events for a given port,
// only emitting a stable event after the port state has been unchanged for
// the configured quiet period.
type Debouncer struct {
	mu      sync.Mutex
	quiet   time.Duration
	timers  map[int]*time.Timer
	callback func(port int, opened bool)
}

// NewDebouncer creates a Debouncer that waits quiet before invoking cb.
func NewDebouncer(quiet time.Duration, cb func(port int, opened bool)) *Debouncer {
	return &Debouncer{
		quiet:    quiet,
		timers:   make(map[int]*time.Timer),
		callback: cb,
	}
}

// Push schedules a debounced event for port. If a pending timer exists for
// the port it is reset, so only the final state within the quiet window fires.
func (d *Debouncer) Push(port int, opened bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[port]; ok {
		t.Stop()
	}

	d.timers[port] = time.AfterFunc(d.quiet, func() {
		d.mu.Lock()
		delete(d.timers, port)
		d.mu.Unlock()
		d.callback(port, opened)
	})
}

// Flush cancels all pending timers without firing them.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for port, t := range d.timers {
		t.Stop()
		delete(d.timers, port)
	}
}

// Pending returns the number of ports currently waiting to fire.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
