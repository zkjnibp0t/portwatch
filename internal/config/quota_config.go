package config

import (
	"time"

	"github.com/user/portwatch/internal/ports"
)

// QuotaConfig holds quota settings from the YAML config.
type QuotaConfig struct {
	Limit  int    `yaml:"limit"`
	Window string `yaml:"window"`
}

// BuildQuotaTracker constructs a QuotaTracker from config, returning nil when
// quota is disabled (limit <= 0).
func BuildQuotaTracker(cfg QuotaConfig) *ports.QuotaTracker {
	if cfg.Limit <= 0 {
		return nil
	}
	window := 1 * time.Minute
	if cfg.Window != "" {
		if d, err := time.ParseDuration(cfg.Window); err == nil && d > 0 {
			window = d
		}
	}
	return ports.NewQuotaTracker(cfg.Limit, window)
}
