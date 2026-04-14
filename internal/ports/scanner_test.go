package ports

import (
	"net"
	"testing"
)

// startTCPListener opens a random local TCP port and returns its port number and a closer func.
func startTCPListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScanFindsOpenPort(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	s := NewScanner(port, port)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
}

func TestScanInvalidRange(t *testing.T) {
	s := NewScanner(9000, 8000)
	_, err := s.Scan()
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestCompareDetectsOpenedAndClosed(t *testing.T) {
	prev := []PortState{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}}
	curr := []PortState{{Port: 443, Protocol: "tcp"}, {Port: 8080, Protocol: "tcp"}}

	diff := Compare(prev, curr)

	if len(diff.Opened) != 1 || diff.Opened[0].Port != 8080 {
		t.Errorf("expected port 8080 opened, got %+v", diff.Opened)
	}
	if len(diff.Closed) != 1 || diff.Closed[0].Port != 80 {
		t.Errorf("expected port 80 closed, got %+v", diff.Closed)
	}
	if !diff.HasChanges() {
		t.Error("expected HasChanges() == true")
	}
}

func TestCompareNoChanges(t *testing.T) {
	ports := []PortState{{Port: 22, Protocol: "tcp"}, {Port: 80, Protocol: "tcp"}}
	diff := Compare(ports, ports)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}
