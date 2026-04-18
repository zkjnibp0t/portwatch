package ports

import (
	"math/rand"
	"time"
)

// JitterConfig holds configuration for scan interval jitter.
type JitterConfig struct {
	Base   time.Duration
	MaxPct float64 // 0.0–1.0, e.g. 0.2 for ±20%
}

// Jitterer adds randomised jitter to a base scan interval so that
// multiple portwatch instances do not hammer the host in lockstep.
type Jitterer struct {
	cfg JitterConfig
	rng *rand.Rand
}

// NewJitterer creates a Jitterer with the given config.
// MaxPct is clamped to [0, 1].
func NewJitterer(cfg JitterConfig) *Jitterer {
	if cfg.MaxPct < 0 {
		cfg.MaxPct = 0
	}
	if cfg.MaxPct > 1 {
		cfg.MaxPct = 1
	}
	return &Jitterer{
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Next returns the base interval plus a random jitter in
// [-MaxPct*Base, +MaxPct*Base].
func (j *Jitterer) Next() time.Duration {
	if j.cfg.MaxPct == 0 || j.cfg.Base == 0 {
		return j.cfg.Base
	}
	max := float64(j.cfg.Base) * j.cfg.MaxPct
	// random in [-max, +max]
	delta := (j.rng.Float64()*2 - 1) * max
	result := time.Duration(float64(j.cfg.Base) + delta)
	if result < 0 {
		return 0
	}
	return result
}

// SetSeed replaces the random source; useful for deterministic tests.
func (j *Jitterer) SetSeed(seed int64) {
	j.rng = rand.New(rand.NewSource(seed))
}
