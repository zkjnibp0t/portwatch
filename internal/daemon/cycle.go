package daemon

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/user/portwatch/internal/ports"
)

const snapshotFile = "portwatch_snapshot.json"

// runCycle performs one scan-diff-alert cycle.
// It loads the previous snapshot, scans current ports, computes the diff,
// dispatches alerts, and persists the new snapshot.
func (d *Daemon) runCycle(ctx context.Context) error {
	snapshotPath := filepath.Join(d.cfg.StateDir, snapshotFile)

	// Load previous snapshot (empty on first run).
	prev, err := ports.LoadSnapshot(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	// Scan current open ports.
	raw, err := d.scanner.Scan(ctx, d.cfg.PortRange)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	// Apply include/exclude filters.
	filtered := d.filter.Apply(raw)

	// Resolve process info if available.
	current := d.resolver.ResolveSet(filtered)

	// Diff against previous state.
	diff := ports.Compare(prev.ToSet(), current)

	// Dispatch alert (throttling handled inside manager).
	if alertErr := d.manager.Dispatch(ctx, diff); alertErr != nil {
		log.Printf("[portwatch] alert dispatch error: %v", alertErr)
	}

	// Record to history store.
	if err := d.store.Record(diff); err != nil {
		log.Printf("[portwatch] history record error: %v", err)
	}

	// Persist new snapshot.
	newSnap := ports.NewSnapshot(current)
	if err := ports.SaveSnapshot(snapshotPath, newSnap); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	return nil
}
