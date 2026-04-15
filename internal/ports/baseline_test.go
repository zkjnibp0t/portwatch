package ports

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaselineStoresPorts(t *testing.T) {
	ports := map[string]bool{"127.0.0.1:80": true, "127.0.0.1:443": true}
	b := NewBaseline(ports)
	if len(b.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(b.Ports))
	}
	if b.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := NewBaseline(map[string]bool{"0.0.0.0:22": true})
	if err := SaveBaseline(path, original); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if !loaded.Ports["0.0.0.0:22"] {
		t.Error("expected port 0.0.0.0:22 to be present after reload")
	}
}

func TestLoadBaselineMissingFile(t *testing.T) {
	b, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if b != nil {
		t.Fatal("expected nil baseline for missing file")
	}
}

func TestBaselineFileIsValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = SaveBaseline(path, NewBaseline(map[string]bool{"0.0.0.0:8080": true}))

	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("baseline file is not valid JSON: %v", err)
	}
}

func TestCompareToBaseline(t *testing.T) {
	b := NewBaseline(map[string]bool{"0.0.0.0:22": true, "0.0.0.0:80": true})
	current := map[string]bool{"0.0.0.0:22": true, "0.0.0.0:443": true}

	diff := CompareToBaseline(b, current)
	if len(diff.Opened) != 1 || diff.Opened[0] != "0.0.0.0:443" {
		t.Errorf("expected 0.0.0.0:443 opened, got %v", diff.Opened)
	}
	if len(diff.Closed) != 1 || diff.Closed[0] != "0.0.0.0:80" {
		t.Errorf("expected 0.0.0.0:80 closed, got %v", diff.Closed)
	}
}

func TestBaselineManagerPrintSummary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	m, _ := NewBaselineManager(path)
	var buf bytes.Buffer
	m.PrintSummary(&buf)
	if buf.String() != "No baseline recorded.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}

	_ = m.Record(map[string]bool{"0.0.0.0:22": true})
	buf.Reset()
	m.PrintSummary(&buf)
	if buf.Len() == 0 {
		t.Error("expected non-empty summary after recording baseline")
	}
}
