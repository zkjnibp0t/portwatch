package ports_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func TestNewSnapshotRoundTrip(t *testing.T) {
	input := map[string]struct{}{
		"127.0.0.1:8080": {},
		"127.0.0.1:9090": {},
	}
	s := ports.NewSnapshot(input)
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	got := s.ToSet()
	if len(got) != len(input) {
		t.Fatalf("expected %d ports, got %d", len(input), len(got))
	}
	for k := range input {
		if _, ok := got[k]; !ok {
			t.Errorf("missing port %s in round-trip set", k)
		}
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := ports.Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Ports:     []string{"0.0.0.0:22", "0.0.0.0:80"},
	}
	if err := ports.SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	loaded, err := ports.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if !loaded.Timestamp.Equal(orig.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", loaded.Timestamp, orig.Timestamp)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: got %d want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoadSnapshotMissingFile(t *testing.T) {
	s, err := ports.LoadSnapshot("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.Ports) != 0 {
		t.Errorf("expected empty snapshot, got %+v", s)
	}
}

func TestSnapshotFileIsValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	s := ports.NewSnapshot(map[string]struct{}{"127.0.0.1:443": {}})
	if err := ports.SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("snapshot file is not valid JSON: %v", err)
	}
	if _, ok := raw["timestamp"]; !ok {
		t.Error("missing 'timestamp' key in JSON output")
	}
	if _, ok := raw["ports"]; !ok {
		t.Error("missing 'ports' key in JSON output")
	}
}
