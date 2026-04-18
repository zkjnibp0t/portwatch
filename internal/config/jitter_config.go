package config

import (
	"time"

	"github.com/user/portwatch/internal/ports"
)

// JitterEntry is the YAML-decoded representation of jitter settings.
type JitterEntry struct {
	Enabled bool    `yaml:"enabled"`
	MaxPct  float64 `yaml:"max_pct"`
}

// BuildJitterer constructs a Jitterer from the config and the resolved
// base scan interval. If jitter is disabled or MaxPct is zero a
// Jitterer with MaxPct=0 is returned (passthrough behaviour).
func BuildJitterer(entry JitterEntry, base time.Duration) *ports.Jitterer {
	if !entry.Enabled || entry.MaxPct <= 0 {
		return ports.NewJitterer(ports.JitterConfig{
			Base:   base,
			MaxPct: 0,
		})
	}
	return ports.NewJitterer(ports.JitterConfig{
		Base:   base,
		MaxPct: entry.MaxPct,
	})
}
