package ports

import (
	"path/filepath"
	"testing"
)

func TestBaselineManagerNoBaseline(t *testing.T) {
	dir := t.TempDir()
	m, err := NewBaselineManager(filepath.Join(dir, "baseline.json"))
	if err != nil {
		t.Fatalf("NewBaselineManager: %v", err)
	}
	if m.HasBaseline() {
		t.Error("expected HasBaseline to be false before recording")
	}
	if !m.CreatedAt().IsZero() {
		t.Error("expected zero CreatedAt when no baseline exists")
	}
}

func TestBaselineManagerRecord(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewBaselineManager(filepath.Join(dir, "baseline.json"))

	ports := map[string]bool{"0.0.0.0:80": true}
	if err := m.Record(ports); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if !m.HasBaseline() {
		t.Error("expected HasBaseline to be true after recording")
	}
	if m.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt after recording")
	}
}

func TestBaselineManagerDiffNoBaseline(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewBaselineManager(filepath.Join(dir, "baseline.json"))

	_, err := m.Diff(map[string]bool{})
	if err == nil {
		t.Fatal("expected error when diffing without a baseline")
	}
}

func TestBaselineManagerDiff(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewBaselineManager(filepath.Join(dir, "baseline.json"))
	_ = m.Record(map[string]bool{"0.0.0.0:22": true})

	current := map[string]bool{"0.0.0.0:22": true, "0.0.0.0:8080": true}
	diff, err := m.Diff(current)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(diff.Opened) != 1 || diff.Opened[0] != "0.0.0.0:8080" {
		t.Errorf("expected 0.0.0.0:8080 opened, got %v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diff.Closed)
	}
}

func TestBaselineManagerPersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	m1, _ := NewBaselineManager(path)
	_ = m1.Record(map[string]bool{"0.0.0.0:443": true})

	m2, err := NewBaselineManager(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !m2.HasBaseline() {
		t.Error("expected baseline to persist across reload")
	}
	diff, _ := m2.Diff(map[string]bool{})
	if len(diff.Closed) != 1 || diff.Closed[0] != "0.0.0.0:443" {
		t.Errorf("expected 0.0.0.0:443 to be closed, got %v", diff.Closed)
	}
}
