package ports

import "fmt"

// Filter holds include/exclude rules for port filtering.
type Filter struct {
	Include map[int]struct{}
	Exclude map[int]struct{}
}

// NewFilter builds a Filter from include and exclude port lists.
// If includeList is empty, all ports are considered included by default.
func NewFilter(includeList, excludeList []int) *Filter {
	f := &Filter{
		Include: make(map[int]struct{}, len(includeList)),
		Exclude: make(map[int]struct{}, len(excludeList)),
	}
	for _, p := range includeList {
		f.Include[p] = struct{}{}
	}
	for _, p := range excludeList {
		f.Exclude[p] = struct{}{}
	}
	return f
}

// Allow returns true if the given port should be monitored.
// Exclude takes precedence over include.
func (f *Filter) Allow(port int) bool {
	if _, excluded := f.Exclude[port]; excluded {
		return false
	}
	if len(f.Include) == 0 {
		return true
	}
	_, included := f.Include[port]
	return included
}

// Apply returns a new PortSet containing only ports allowed by the filter.
func (f *Filter) Apply(ps PortSet) PortSet {
	filtered := make(PortSet, len(ps))
	for port := range ps {
		if f.Allow(port) {
			filtered[port] = struct{}{}
		}
	}
	return filtered
}

// String returns a human-readable description of the filter rules.
func (f *Filter) String() string {
	return fmt.Sprintf("Filter{include=%d ports, exclude=%d ports}",
		len(f.Include), len(f.Exclude))
}
