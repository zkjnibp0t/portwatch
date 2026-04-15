package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a persisted reference port set used to detect
// deviations from a known-good state.
type Baseline struct {
	CreatedAt time.Time        `json:"created_at"`
	Ports     map[string]bool  `json:"ports"`
}

// NewBaseline creates a Baseline from the given port set.
func NewBaseline(ports map[string]bool) *Baseline {
	copy := make(map[string]bool, len(ports))
	for k, v := range ports {
		copy[k] = v
	}
	return &Baseline{
		CreatedAt: time.Now().UTC(),
		Ports:     copy,
	}
}

// SaveBaseline writes a Baseline to the given file path as JSON.
func SaveBaseline(path string, b *Baseline) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadBaseline reads a Baseline from the given file path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("baseline: read: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// CompareToBaseline returns the diff between the current port set and the
// baseline. Ports present in current but not baseline are Opened; ports in
// baseline but not current are Closed.
func CompareToBaseline(b *Baseline, current map[string]bool) Diff {
	return Compare(b.Ports, current)
}
