package ports

// Severity represents the urgency level of a port change event.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityCritical:
		return "critical"
	default:
		return "info"
	}
}

// SeverityRule maps a port range to a severity level.
type SeverityRule struct {
	MinPort int
	MaxPort int
	Level   Severity
}

// SeverityClassifier assigns severity levels to ports based on rules.
type SeverityClassifier struct {
	rules []SeverityRule
	defaultLevel Severity
}

// NewSeverityClassifier creates a classifier with the given rules.
// Rules are evaluated in order; first match wins.
func NewSeverityClassifier(rules []SeverityRule, defaultLevel Severity) *SeverityClassifier {
	return &SeverityClassifier{rules: rules, defaultLevel: defaultLevel}
}

// Classify returns the severity for a given port number.
func (c *SeverityClassifier) Classify(port int) Severity {
	for _, r := range c.rules {
		if port >= r.MinPort && port <= r.MaxPort {
			return r.Level
		}
	}
	return c.defaultLevel
}

// ClassifySet returns the highest severity found across a set of ports.
func (c *SeverityClassifier) ClassifySet(ports []int) Severity {
	max := c.defaultLevel
	for _, p := range ports {
		if s := c.Classify(p); s > max {
			max = s
		}
	}
	return max
}
