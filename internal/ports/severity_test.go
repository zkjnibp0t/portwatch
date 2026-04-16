package ports

import (
	"testing"
)

func TestSeverityString(t *testing.T) {
	if SeverityInfo.String() != "info" {
		t.Errorf("expected info")
	}
	if SeverityWarning.String() != "warning" {
		t.Errorf("expected warning")
	}
	if SeverityCritical.String() != "critical" {
		t.Errorf("expected critical")
	}
}

func TestClassifyMatchesRule(t *testing.T) {
	c := NewSeverityClassifier([]SeverityRule{
		{MinPort: 1, MaxPort: 1023, Level: SeverityCritical},
		{MinPort: 1024, MaxPort: 49151, Level: SeverityWarning},
	}, SeverityInfo)

	if got := c.Classify(80); got != SeverityCritical {
		t.Errorf("port 80: expected critical, got %s", got)
	}
	if got := c.Classify(8080); got != SeverityWarning {
		t.Errorf("port 8080: expected warning, got %s", got)
	}
	if got := c.Classify(60000); got != SeverityInfo {
		t.Errorf("port 60000: expected info, got %s", got)
	}
}

func TestClassifyDefaultWhenNoMatch(t *testing.T) {
	c := NewSeverityClassifier(nil, SeverityWarning)
	if got := c.Classify(9999); got != SeverityWarning {
		t.Errorf("expected default warning, got %s", got)
	}
}

func TestClassifySetReturnsHighest(t *testing.T) {
	c := NewSeverityClassifier([]SeverityRule{
		{MinPort: 22, MaxPort: 22, Level: SeverityCritical},
	}, SeverityInfo)

	got := c.ClassifySet([]int{8080, 22, 9000})
	if got != SeverityCritical {
		t.Errorf("expected critical, got %s", got)
	}
}

func TestClassifySetEmptyReturnsDefault(t *testing.T) {
	c := NewSeverityClassifier(nil, SeverityInfo)
	if got := c.ClassifySet(nil); got != SeverityInfo {
		t.Errorf("expected info, got %s", got)
	}
}
