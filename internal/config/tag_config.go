package config

import "github.com/user/portwatch/internal/ports"

// TagEntry represents a tag definition in the config file.
type TagEntry struct {
	Name  string `yaml:"name"`
	Ports []int  `yaml:"ports"`
}

// BuildTagger constructs a Tagger from config tag entries.
func BuildTagger(entries []TagEntry) *ports.Tagger {
	tags := make([]ports.Tag, 0, len(entries))
	for _, e := range entries {
		if e.Name == "" || len(e.Ports) == 0 {
			continue
		}
		tags = append(tags, ports.Tag{
			Name:  e.Name,
			Ports: e.Ports,
		})
	}
	return ports.NewTagger(tags)
}
