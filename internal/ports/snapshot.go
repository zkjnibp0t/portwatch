package ports

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []string  `json:"ports"`
}

// NewSnapshot creates a Snapshot from a port set at the current time.
func NewSnapshot(set map[string]struct{}) Snapshot {
	ports := make([]string, 0, len(set))
	for p := range set {
		ports = append(ports, p)
	}
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
}

// ToSet converts the snapshot's port list back into a set for diffing.
func (s Snapshot) ToSet() map[string]struct{} {
	out := make(map[string]struct{}, len(s.Ports))
	for _, p := range s.Ports {
		out[p] = struct{}{}
	}
	return out
}

// SaveSnapshot writes a snapshot as JSON to the given file path.
func SaveSnapshot(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadSnapshot reads a snapshot from a JSON file.
// If the file does not exist, an empty snapshot is returned without error.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, err
	}
	return s, nil
}
