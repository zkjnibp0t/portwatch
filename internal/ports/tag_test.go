package ports

import (
	"testing"
)

func TestTaggerLabelKnownPort(t *testing.T) {
	tagger := NewTagger([]Tag{
		{Name: "web", Ports: []int{80, 443}},
	})
	if got := tagger.Label(80); got != "web" {
		t.Errorf("expected web, got %s", got)
	}
	if got := tagger.Label(443); got != "web" {
		t.Errorf("expected web, got %s", got)
	}
}

func TestTaggerLabelUnknownPort(t *testing.T) {
	tagger := NewTagger([]Tag{})
	if got := tagger.Label(9999); got != "untagged" {
		t.Errorf("expected untagged, got %s", got)
	}
}

func TestTaggerTagSet(t *testing.T) {
	tagger := NewTagger([]Tag{
		{Name: "db", Ports: []int{5432, 3306}},
	})
	set := map[int]struct{}{5432: {}, 8080: {}}
	result := tagger.TagSet(set)
	if result[5432] != "db" {
		t.Errorf("expected db for 5432")
	}
	if result[8080] != "untagged" {
		t.Errorf("expected untagged for 8080")
	}
}

func TestTaggerGroups(t *testing.T) {
	tagger := NewTagger([]Tag{
		{Name: "web", Ports: []int{80, 443}},
		{Name: "db", Ports: []int{5432}},
	})
	set := map[int]struct{}{80: {}, 443: {}, 5432: {}, 9999: {}}
	groups := tagger.Groups(set)
	if len(groups["web"]) != 2 {
		t.Errorf("expected 2 web ports, got %d", len(groups["web"]))
	}
	if len(groups["db"]) != 1 {
		t.Errorf("expected 1 db port")
	}
	if len(groups["untagged"]) != 1 {
		t.Errorf("expected 1 untagged port")
	}
}

func TestTaggerEmptySet(t *testing.T) {
	tagger := NewTagger([]Tag{{Name: "web", Ports: []int{80}}})
	result := tagger.TagSet(map[int]struct{}{})
	if len(result) != 0 {
		t.Errorf("expected empty result")
	}
}
