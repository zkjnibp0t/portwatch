package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	PortRange string        `yaml:"port_range"`
	Interval  time.Duration `yaml:"interval"`
	Webhook   WebhookConfig `yaml:"webhook"`
	Desktop   DesktopConfig `yaml:"desktop"`
	HistoryFile string     `yaml:"history_file"`
}

// WebhookConfig holds webhook notifier settings.
type WebhookConfig struct {
	URL     string `yaml:"url"`
	Enabled bool   `yaml:"enabled"`
}

// DesktopConfig holds desktop notifier settings.
type DesktopConfig struct {
	Enabled bool   `yaml:"enabled"`
	AppName string `yaml:"app_name"`
}

func defaultConfig() Config {
	return Config{
		PortRange:   "1-1024",
		Interval:    30 * time.Second,
		HistoryFile: "portwatch_history.json",
		Desktop: DesktopConfig{
			AppName: "portwatch",
		},
	}
}

// Load reads configuration from path, falling back to defaults.
func Load(path string) (Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Interval < 5*time.Second {
		return errors.New("interval must be at least 5s")
	}
	if cfg.PortRange == "" {
		return errors.New("port_range must not be empty")
	}
	return nil
}
