package config

import "github.com/user/portwatch/internal/ports"

// DecayConfig holds configuration for the decay tracker.
type DecayConfig struct {
	DecayRate      float64 `yaml:"decay_rate"`
	PruneThreshold float64 `yaml:"prune_threshold"`
}

// BuildDecayTracker constructs a DecayTracker from config.
// Sensible defaults are applied when values are zero.
func BuildDecayTracker(cfg DecayConfig) *ports.DecayTracker {
	rate := cfg.DecayRate
	if rate <= 0 {
		rate = 0.05 // 5% decay per second by default
	}
	return ports.NewDecayTracker(rate)
}

// DefaultDecayConfig returns a DecayConfig with recommended defaults.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		DecayRate:      0.05,
		PruneThreshold: 0.01,
	}
}
