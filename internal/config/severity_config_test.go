package config

import (
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func TestBuildSeverityClassifierEmpty(t *testing.T) {
	c := BuildSeverityClassifier(nil, "info")
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
	if got := c.Classify(80); got != ports.SeverityInfo {
		t.Errorf("expected info, got %s", got)
	}
}

func TestBuildSeverityClassifierSingleRule(t *testing.T) {
	rules := []SeverityRuleConfig{
		{MinPort: 1, MaxPort: 1023, Severity: "critical"},
	}
	c := BuildSeverityClassifier(rules, "info")
	if got := c.Classify(443); got != ports.SeverityCritical {
		t.Errorf("expected critical, got %s", got)
	}
	if got := c.Classify(3000); got != ports.SeverityInfo {
		t.Errorf("expected info, got %s", got)
	}
}

func TestBuildSeverityClassifierDefaultWarning(t *testing.T) {
	c := BuildSeverityClassifier(nil, "warning")
	if got := c.Classify(50000); got != ports.SeverityWarning {
		t.Errorf("expected warning, got %s", got)
	}
}

func TestBuildSeverityClassifierMultipleRules(t *testing.T) {
	rules := []SeverityRuleConfig{
		{MinPort: 1, MaxPort: 1023, Severity: "critical"},
		{MinPort: 1024, MaxPort: 49151, Severity: "warning"},
	}
	c := BuildSeverityClassifier(rules, "info")
	if got := c.Classify(22); got != ports.SeverityCritical {
		t.Errorf("port 22: expected critical")
	}
	if got := c.Classify(8443); got != ports.SeverityWarning {
		t.Errorf("port 8443: expected warning")
	}
	if got := c.Classify(55000); got != ports.SeverityInfo {
		t.Errorf("port 55000: expected info")
	}
}

func TestParseSeverityUnknownFallsBackToInfo(t *testing.T) {
	if got := parseSeverity("unknown"); got != ports.SeverityInfo {
		t.Errorf("expected info fallback")
	}
}
