package ports

// Tag represents a user-defined label applied to a port range.
type Tag struct {
	Name  string
	Ports []int
}

// Tagger maps ports to human-readable tags.
type Tagger struct {
	tags []Tag
	index map[int]string
}

// NewTagger builds a Tagger from a slice of Tag definitions.
func NewTagger(tags []Tag) *Tagger {
	index := make(map[int]string, len(tags)*4)
	for _, t := range tags {
		for _, p := range t.Ports {
			index[p] = t.Name
		}
	}
	return &Tagger{tags: tags, index: index}
}

// Label returns the tag name for a port, or "untagged" if none matches.
func (t *Tagger) Label(port int) string {
	if name, ok := t.index[port]; ok {
		return name
	}
	return "untagged"
}

// TagSet annotates a set of ports, returning a map of port -> tag.
func (t *Tagger) TagSet(ports map[int]struct{}) map[int]string {
	result := make(map[int]string, len(ports))
	for p := range ports {
		result[p] = t.Label(p)
	}
	return result
}

// Groups returns ports grouped by tag name.
func (t *Tagger) Groups(ports map[int]struct{}) map[string][]int {
	groups := make(map[string][]int)
	for p := range ports {
		label := t.Label(p)
		groups[label] = append(groups[label], p)
	}
	return groups
}
