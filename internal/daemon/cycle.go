package daemon

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/ports"
)

// cycleHandler processes a single WatchEvent produced by a ports.Watcher.
type cycleHandler struct {
	manager *alert.Manager
	store   *history.Store
}

newCycleHandler := func(m *alert.Manager, s *history.Store) *cycleHandler {
	return &cycleHandler{manager: m, store: s}
}

// handle processes one event: records history and dispatches alerts.
func (h *cycleHandler) handle(ctx context.Context, ev ports.WatchEvent) {
	if ev.Err != nil {
		log.Printf("[portwatch] scan error: %v", ev.Err)
		return
	}
	if len(ev.Diff.Opened) == 0 && len(ev.Diff.Closed) == 0 {
		return
	}
	if err := h.store.Record(ev.Diff); err != nil {
		log.Printf("[portwatch] history record error: %v", err)
	}
	if err := h.manager.Dispatch(ctx, ev.Diff); err != nil {
		log.Printf("[portwatch] alert dispatch error: %v", err)
	}
}

// runCycles drives the event loop, delegating each event to handle.
func runCycles(ctx context.Context, ch <-chan ports.WatchEvent, h *cycleHandler) {
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return
			}
			h.handle(ctx, ev)
		case <-ctx.Done():
			return
		}
	}
}
