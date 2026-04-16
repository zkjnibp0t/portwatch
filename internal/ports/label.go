package ports

import "fmt"

// PortLabel associates a port number with a human-readable service name.
type PortLabel struct {
	Port    int
	Service string
}

// Labeler maps port numbers to well-known service names.
type Labeler struct {
	known map[int]string
}

var defaultLabels = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// NewLabeler creates a Labeler seeded with default well-known port mappings.
// Additional entries from extras are merged, overriding defaults.
func NewLabeler(extras map[int]string) *Labeler {
	m := make(map[int]string, len(defaultLabels)+len(extras))
	for k, v := range defaultLabels {
		m[k] = v
	}
	for k, v := range extras {
		if v != "" {
			m[k] = v
		}
	}
	return &Labeler{known: m}
}

// Label returns the service name for a port, or a numeric fallback.
func (l *Labeler) Label(port int) string {
	if name, ok := l.known[port]; ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// LabelSet returns a slice of PortLabel for every port in the set.
func (l *Labeler) LabelSet(ports map[int]struct{}) []PortLabel {
	out := make([]PortLabel, 0, len(ports))
	for p := range ports {
		out = append(out, PortLabel{Port: p, Service: l.Label(p)})
	}
	return out
}
