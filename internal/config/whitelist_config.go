package config

import "github.com/user/portwatch/internal/ports"

// WhitelistEntryConfig is the YAML-serialisable form of a whitelist rule.
type WhitelistEntryConfig struct {
	Port    int    `yaml:"port"`
	Process string `yaml:"process"`
}

// BuildWhitelist converts the slice of config entries into a *ports.Whitelist
// that can be used at runtime.
func BuildWhitelist(entries []WhitelistEntryConfig) *ports.Whitelist {
	pe := make([]ports.WhitelistEntry, len(entries))
	for i, e := range entries {
		pe[i] = ports.WhitelistEntry{
			Port:    e.Port,
			Process: e.Process,
		}
	}
	return ports.NewWhitelist(pe)
}
