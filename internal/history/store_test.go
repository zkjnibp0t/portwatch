package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestRecordAndRetrieve(t *testing.T) {
	p := tempPath(t)
	s, err := history.NewStore(p)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := s.Record([]string{"8080"}, nil); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := s.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Opened) != 1 || entries[0].Opened[0] != "8080" {
		t.Errorf("unexpected opened: %v", entries[0].Opened)
	}
}

func TestRecordNoChangesSkipped(t *testing.T) {
	p := tempPath(t)
	s, _ := history.NewStore(p)
	_ = s.Record(nil, nil)
	if len(s.Entries()) != 0 {
		t.Error("expected no entries for empty diff")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	s1, _ := history.NewStore(p)
	_ = s1.Record([]string{"443"}, []string{"80"})

	s2, err := history.NewStore(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(s2.Entries()) != 1 {
		t.Fatalf("expected 1 persisted entry")
	}
}

func TestFileContainsValidJSON(t *testing.T) {
	p := tempPath(t)
	s, _ := history.NewStore(p)
	_ = s.Record([]string{"22"}, nil)

	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var entries []history.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
}
