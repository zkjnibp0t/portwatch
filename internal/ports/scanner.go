package ports

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single open port.
type PortState struct {
	Port     int
	Protocol string
	Address  string
}

// String returns a human-readable representation of a PortState.
func (p PortState) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(start, end int) *Scanner {
	return &Scanner{
		StartPort: start,
		EndPort:   end,
		Timeout:   500 * time.Millisecond,
	}
}

// Scan probes TCP ports in the configured range and returns open ones.
func (s *Scanner) Scan() ([]PortState, error) {
	if s.StartPort < 1 || s.EndPort > 65535 || s.StartPort > s.EndPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", s.StartPort, s.EndPort)
	}

	var open []PortState
	for port := s.StartPort; port <= s.EndPort; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, PortState{
				Port:     port,
				Protocol: "tcp",
				Address:  "127.0.0.1",
			})
		}
	}
	return open, nil
}

// ToSet converts a slice of PortState into a map keyed by port number for O(1) lookup.
func ToSet(states []PortState) map[int]PortState {
	set := make(map[int]PortState, len(states))
	for _, s := range states {
		set[s.Port] = s
	}
	return set
}
