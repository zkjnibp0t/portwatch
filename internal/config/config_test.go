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
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadDefaults(t *testing.T) {
	path := writeTempConfig(t, "{}\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.PortStart != 1 {
		t.Errorf("expected default port_start 1, got %d", cfg.Scan.PortStart)
	}
	if cfg.Scan.PortEnd != 65535 {
		t.Errorf("expected default port_end 65535, got %d", cfg.Scan.PortEnd)
	}
	if cfg.Scan.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %s", cfg.Scan.Interval)
	}
	if !cfg.Notify.Desktop {
		t.Error("expected desktop notifications enabled by default")
	}
	if cfg.Notify.AppName != "portwatch" {
		t.Errorf("expected default app_name 'portwatch', got %q", cfg.Notify.AppName)
	}
}

func TestLoadCustomValues(t *testing.T) {
	content := `
scan:
  port_start: 1024
  port_end: 9000
  interval: 60s
notify:
  webhook_url: "https://example.com/hook"
  desktop: false
  app_name: "myapp"
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.PortStart != 1024 {
		t.Errorf("expected port_start 1024, got %d", cfg.Scan.PortStart)
	}
	if cfg.Notify.WebhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook_url: %s", cfg.Notify.WebhookURL)
	}
	if cfg.Notify.Desktop {
		t.Error("expected desktop to be false")
	}
}

func TestLoadInvalidPortRange(t *testing.T) {
	content := "scan:\n  port_start: 9000\n  port_end: 1024\n"
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid port range, got nil")
	}
}

func TestLoadIntervalTooShort(t *testing.T) {
	content := "scan:\n  interval: 500ms\n"
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for interval < 1s, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
