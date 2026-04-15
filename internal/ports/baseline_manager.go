package ports

import (
	"fmt"
	"io"
	"time"
)

// BaselineManager handles loading, saving and comparing against a port baseline.
type BaselineManager struct {
	path    string
	baseline *Baseline
}

// NewBaselineManager creates a BaselineManager backed by the given file path.
// If a baseline file already exists it is loaded automatically.
func NewBaselineManager(path string) (*BaselineManager, error) {
	b, err := LoadBaseline(path)
	if err != nil {
		return nil, fmt.Errorf("baseline manager: %w", err)
	}
	return &BaselineManager{path: path, baseline: b}, nil
}

// HasBaseline reports whether a baseline has been established.
func (m *BaselineManager) HasBaseline() bool {
	return m.baseline != nil
}

// CreatedAt returns the time the baseline was recorded, or the zero time if
// no baseline exists.
func (m *BaselineManager) CreatedAt() time.Time {
	if m.baseline == nil {
		return time.Time{}
	}
	return m.baseline.CreatedAt
}

// Record captures the given port set as the new baseline and persists it.
func (m *BaselineManager) Record(ports map[string]bool) error {
	b := NewBaseline(ports)
	if err := SaveBaseline(m.path, b); err != nil {
		return fmt.Errorf("baseline manager: record: %w", err)
	}
	m.baseline = b
	return nil
}

// Diff returns the diff between the stored baseline and the current port set.
// Returns an error if no baseline has been recorded yet.
func (m *BaselineManager) Diff(current map[string]bool) (Diff, error) {
	if m.baseline == nil {
		return Diff{}, fmt.Errorf("baseline manager: no baseline recorded")
	}
	return CompareToBaseline(m.baseline, current), nil
}

// PrintSummary writes a human-readable summary of the baseline to w.
func (m *BaselineManager) PrintSummary(w io.Writer) {
	if m.baseline == nil {
		fmt.Fprintln(w, "No baseline recorded.")
		return
	}
	fmt.Fprintf(w, "Baseline recorded at: %s\n", m.baseline.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Baseline port count:  %d\n", len(m.baseline.Ports))
}
