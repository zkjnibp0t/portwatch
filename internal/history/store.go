package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a recorded port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []string  `json:"opened,omitempty"`
	Closed    []string  `json:"closed,omitempty"`
}

// Store persists port change history to a JSON file.
type Store struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// NewStore creates a Store backed by the given file path.
// Existing entries are loaded if the file exists.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Record appends a new entry and flushes to disk.
func (s *Store) Record(opened, closed []string) error {
	if len(opened) == 0 && len(closed) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, Entry{
		Timestamp: time.Now().UTC(),
		Opened:    opened,
		Closed:    closed,
	})
	return s.flush()
}

// Entries returns a copy of all recorded entries.
func (s *Store) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
