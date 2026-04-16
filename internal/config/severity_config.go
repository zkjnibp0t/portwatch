package config

import "github.com/user/portwatch/internal/ports"

// SeverityRuleConfig holds a single severity rule from YAML config.
type SeverityRuleConfig struct {
	MinPort  int    `yaml:"min_port"`
	MaxPort  int    `yaml:"max_port"`
	Severity string `yaml:"severity"`
}

// BuildSeverityClassifier constructs a SeverityClassifier from config rules.
func BuildSeverityClassifier(rules []SeverityRuleConfig, defaultLevel string) *ports.SeverityClassifier {
	parsed := make([]ports.SeverityRule, 0, len(rules))
	for _, r := range rules {
		parsed = append(parsed, ports.SeverityRule{
			MinPort: r.MinPort,
			MaxPort: r.MaxPort,
			Level:   parseSeverity(r.Severity),
		})
	}
	return ports.NewSeverityClassifier(parsed, parseSeverity(defaultLevel))
}

func parseSeverity(s string) ports.Severity {
	switch s {
	case "critical":
		return ports.SeverityCritical
	case "warning":
		return ports.SeverityWarning
	default:
		return ports.SeverityInfo
	}
}
