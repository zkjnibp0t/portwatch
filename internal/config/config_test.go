package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", cfg.Interval)
	}
	if cfg.AppName != "portwatch" {
		t.Errorf("expected default app_name 'portwatch', got %q", cfg.AppName)
	}
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 65535 {
		t.Errorf("unexpected default port range: %d-%d", cfg.PortRange.Start, cfg.PortRange.End)
	}
}

func TestLoadCustomValues(t *testing.T) {
	path := writeTempConfig(t, `
port_range:
  start: 1024
  end: 9000
interval: 60s
app_name: myapp
webhook_url: https://example.com/hook
filter:
  include: [80, 443]
  exclude: [8080]
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PortRange.Start != 1024 {
		t.Errorf("expected start 1024, got %d", cfg.PortRange.Start)
	}
	if cfg.AppName != "myapp" {
		t.Errorf("expected app_name 'myapp', got %q", cfg.AppName)
	}
	if len(cfg.Filter.Include) != 2 {
		t.Errorf("expected 2 include ports, got %d", len(cfg.Filter.Include))
	}
	if len(cfg.Filter.Exclude) != 1 || cfg.Filter.Exclude[0] != 8080 {
		t.Errorf("expected exclude [8080], got %v", cfg.Filter.Exclude)
	}
}

func TestLoadInvalidPortRange(t *testing.T) {
	path := writeTempConfig(t, "port_range:\n  start: 9000\n  end: 1000\n")
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid port range")
	}
}

func TestLoadIntervalTooShort(t *testing.T) {
	path := writeTempConfig(t, "interval: 1s\n")
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for interval < 5s")
	}
}

func TestLoadInvalidFilterPort(t *testing.T) {
	path := writeTempConfig(t, "filter:\n  include: [0]\n")
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for out-of-range filter port 0")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Error("expected error when loading non-existent config file")
	}
}
