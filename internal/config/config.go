package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	Scan     ScanConfig     `yaml:"scan"`
	Notify   NotifyConfig   `yaml:"notify"`
}

// ScanConfig defines port scanning parameters.
type ScanConfig struct {
	PortStart int           `yaml:"port_start"`
	PortEnd   int           `yaml:"port_end"`
	Interval  time.Duration `yaml:"interval"`
}

// NotifyConfig defines notification targets.
type NotifyConfig struct {
	WebhookURL  string `yaml:"webhook_url"`
	Desktop     bool   `yaml:"desktop"`
	AppName     string `yaml:"app_name"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := defaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Scan: ScanConfig{
			PortStart: 1,
			PortEnd:   65535,
			Interval:  30 * time.Second,
		},
		Notify: NotifyConfig{
			Desktop: true,
			AppName: "portwatch",
		},
	}
}

func (c *Config) validate() error {
	if c.Scan.PortStart < 1 || c.Scan.PortStart > 65535 {
		return fmt.Errorf("scan.port_start must be between 1 and 65535, got %d", c.Scan.PortStart)
	}
	if c.Scan.PortEnd < 1 || c.Scan.PortEnd > 65535 {
		return fmt.Errorf("scan.port_end must be between 1 and 65535, got %d", c.Scan.PortEnd)
	}
	if c.Scan.PortStart > c.Scan.PortEnd {
		return fmt.Errorf("scan.port_start (%d) must be <= scan.port_end (%d)", c.Scan.PortStart, c.Scan.PortEnd)
	}
	if c.Scan.Interval < time.Second {
		return fmt.Errorf("scan.interval must be at least 1s, got %s", c.Scan.Interval)
	}
	return nil
}
