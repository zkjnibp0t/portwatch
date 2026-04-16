package ports

import (
	"fmt"
	"testing"
)

func TestLabelerKnownPort(t *testing.T) {
	l := NewLabeler(nil)
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestLabelerUnknownPort(t *testing.T) {
	l := NewLabeler(nil)
	expected := fmt.Sprintf("port/%d", 9999)
	if got := l.Label(9999); got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestLabelerCustomOverridesDefault(t *testing.T) {
	l := NewLabeler(map[int]string{80: "my-http"})
	if got := l.Label(80); got != "my-http" {
		t.Fatalf("expected my-http, got %s", got)
	}
}

func TestLabelerCustomNewEntry(t *testing.T) {
	l := NewLabeler(map[int]string{9200: "elasticsearch"})
	if got := l.Label(9200); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %s", got)
	}
}

func TestLabelerSkipsEmptyValue(t *testing.T) {
	l := NewLabeler(map[int]string{22: ""})
	// empty string should not override default
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestLabelerLabelSet(t *testing.T) {
	l := NewLabeler(nil)
	set := map[int]struct{}{22: {}, 80: {}, 9999: {}}
	labels := l.LabelSet(set)
	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}
	byPort := make(map[int]string)
	for _, pl := range labels {
		byPort[pl.Port] = pl.Service
	}
	if byPort[22] != "ssh" {
		t.Errorf("expected ssh for 22, got %s", byPort[22])
	}
	if byPort[80] != "http" {
		t.Errorf("expected http for 80, got %s", byPort[80])
	}
	if byPort[9999] != "port/9999" {
		t.Errorf("expected port/9999 for 9999, got %s", byPort[9999])
	}
}

func TestLabelerEmptySet(t *testing.T) {
	l := NewLabeler(nil)
	labels := l.LabelSet(map[int]struct{}{})
	if len(labels) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(labels))
	}
}
