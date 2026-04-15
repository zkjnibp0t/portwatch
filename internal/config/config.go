package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	PortRange struct {
		Start int `yaml:"start"`
		End   int `yaml:"end"`
	} `yaml:"port_range"`

	Interval time.Duration `yaml:"interval"`
	Cooldown time.Duration `yaml:"cooldown"`

	Filter struct {
		Include []int `yaml:"include"`
		Exclude []int `yaml:"exclude"`
	} `yaml:"filter"`

	WebhookURL  string `yaml:"webhook_url"`
	AppName     string `yaml:"app_name"`
	HistoryFile string `yaml:"history_file"`
}

func defaultConfig() Config {
	var c Config
	c.PortRange.Start = 1
	c.PortRange.End = 65535
	c.Interval = 30 * time.Second
	c.Cooldown = 5 * time.Minute
	c.AppName = "portwatch"
	c.HistoryFile = "portwatch_history.json"
	return c
}

// Load reads a YAML config file from path, falling back to defaults for
// unset fields. If path is empty, only defaults are returned.
func Load(path string) (Config, error) {
	cfg := defaultConfig()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, validate(cfg)
}

func validate(c Config) error {
	if c.PortRange.Start < 1 || c.PortRange.End > 65535 {
		return errors.New("port range must be between 1 and 65535")
	}
	if c.PortRange.Start > c.PortRange.End {
		return errors.New("port_range.start must be <= port_range.end")
	}
	if c.Interval < 5*time.Second {
		return errors.New("interval must be at least 5s")
	}
	for _, p := range append(c.Filter.Include, c.Filter.Exclude...) {
		if p < 1 || p > 65535 {
			return errors.New("filter port out of valid range 1-65535")
		}
	}
	return nil
}
