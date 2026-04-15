package ports

import (
	"fmt"
	"net"
	"strconv"
)

// ProcessInfo holds basic information about a process listening on a port.
type ProcessInfo struct {
	Port    int
	Proto   string
	Address string
	Host    string
}

// String returns a human-readable representation of ProcessInfo.
func (p ProcessInfo) String() string {
	if p.Host != "" {
		return fmt.Sprintf("%s:%d (%s) -> %s", p.Proto, p.Port, p.Address, p.Host)
	}
	return fmt.Sprintf("%s:%d (%s)", p.Proto, p.Port, p.Address)
}

// Resolver performs reverse DNS lookups for open ports.
type Resolver struct {
	lookup func(addr string) ([]string, error)
}

// NewResolver creates a Resolver using the system DNS.
func NewResolver() *Resolver {
	return &Resolver{
		lookup: net.LookupAddr,
	}
}

// ResolvePort attempts a reverse DNS lookup for the given port's address
// and returns a ProcessInfo. If lookup fails the Host field is left empty.
func (r *Resolver) ResolvePort(port int, proto string) ProcessInfo {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	info := ProcessInfo{
		Port:    port,
		Proto:   proto,
		Address: addr,
	}

	hosts, err := r.lookup("127.0.0.1")
	if err == nil && len(hosts) > 0 {
		info.Host = hosts[0]
	}
	return info
}

// ResolveSet resolves all ports in the given set and returns a slice of ProcessInfo.
func (r *Resolver) ResolveSet(ports map[int]struct{}, proto string) []ProcessInfo {
	results := make([]ProcessInfo, 0, len(ports))
	for port := range ports {
		results = append(results, r.ResolvePort(port, proto))
	}
	return results
}
