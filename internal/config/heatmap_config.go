package config

import (
	"github.com/user/portwatch/internal/ports"
)

// HeatmapConfig holds configuration for the HeatmapHook.
type HeatmapConfig struct {
	// TopN is the number of top ports to report. Defaults to 10.
	TopN int `yaml:"top_n"`
	// Every controls how many scan cycles pass before logging. Defaults to 5.
	Every int `yaml:"every"`
}

// defaultHeatmapConfig returns sensible defaults.
func defaultHeatmapConfig() HeatmapConfig {
	return HeatmapConfig{
		TopN:  10,
		Every: 5,
	}
}

// BuildHeatmapTracker constructs a HeatmapTracker from the config.
// The returned tracker is ready to be passed to NewHeatmapHook.
func BuildHeatmapTracker(cfg HeatmapConfig) (*ports.HeatmapTracker, HeatmapConfig) {
	def := defaultHeatmapConfig()
	if cfg.TopN <= 0 {
		cfg.TopN = def.TopN
	}
	if cfg.Every <= 0 {
		cfg.Every = def.Every
	}
	return ports.NewHeatmapTracker(), cfg
}
