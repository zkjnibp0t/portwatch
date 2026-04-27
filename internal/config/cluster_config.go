package config

import "github.com/user/portwatch/internal/ports"

// ClusterConfig holds configuration for the cluster detector.
type ClusterConfig struct {
	MinSupport int `yaml:"min_support"`
}

// defaultClusterConfig returns sensible defaults for cluster detection.
func defaultClusterConfig() ClusterConfig {
	return ClusterConfig{
		MinSupport: 3,
	}
}

// BuildClusterDetector constructs a ClusterDetector from the provided config,
// falling back to defaults for zero values.
func BuildClusterDetector(cfg ClusterConfig) *ports.ClusterDetector {
	if cfg.MinSupport <= 0 {
		cfg = defaultClusterConfig()
	}
	return ports.NewClusterDetector(cfg.MinSupport)
}
