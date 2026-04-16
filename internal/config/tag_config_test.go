package config

import (
	"testing"
)

func TestBuildTaggerEmpty(t *testing.T) {
	tagger := BuildTagger([]TagEntry{})
	if got := tagger.Label(80); got != "untagged" {
		t.Errorf("expected untagged, got %s", got)
	}
}

func TestBuildTaggerSingleEntry(t *testing.T) {
	tagger := BuildTagger([]TagEntry{
		{Name: "web", Ports: []int{80, 443}},
	})
	if got := tagger.Label(80); got != "web" {
		t.Errorf("expected web, got %s", got)
	}
}

func TestBuildTaggerSkipsEmptyName(t *testing.T) {
	tagger := BuildTagger([]TagEntry{
		{Name: "", Ports: []int{80}},
		{Name: "db", Ports: []int{5432}},
	})
	if got := tagger.Label(80); got != "untagged" {
		t.Errorf("empty-name entry should be skipped")
	}
	if got := tagger.Label(5432); got != "db" {
		t.Errorf("expected db, got %s", got)
	}
}

func TestBuildTaggerSkipsEmptyPorts(t *testing.T) {
	tagger := BuildTagger([]TagEntry{
		{Name: "empty", Ports: []int{}},
	})
	if got := tagger.Label(0); got != "untagged" {
		t.Errorf("entry with no ports should be skipped")
	}
}

func TestBuildTaggerMultipleEntries(t *testing.T) {
	tagger := BuildTagger([]TagEntry{
		{Name: "web", Ports: []int{80}},
		{Name: "db", Ports: []int{5432}},
		{Name: "cache", Ports: []int{6379}},
	})
	for port, want := range map[int]string{80: "web", 5432: "db", 6379: "cache", 9999: "untagged"} {
		if got := tagger.Label(port); got != want {
			t.Errorf("port %d: expected %s, got %s", port, want, got)
		}
	}
}
