package config

import "github.com/user/portwatch/internal/ports"

// LabelEntry represents a single port-to-name mapping in the config file.
type LabelEntry struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

// BuildLabeler constructs a ports.Labeler from config-level label entries.
func BuildLabeler(entries []LabelEntry) *ports.Labeler {
	extras := make(map[int]string, len(entries))
	for _, e := range entries {
		if e.Port > 0 && e.Name != "" {
			extras[e.Port] = e.Name
		}
	}
	return ports.NewLabeler(extras)
}
