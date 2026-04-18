package ports

import (
	"fmt"
	"testing"
)

func TestResolvePortNoHost(t *testing.T) {
	r := &Resolver{
		lookup: func(addr string) ([]string, error) {
			return nil, fmt.Errorf("no such host")
		},
	}

	info := r.ResolvePort(8080, "tcp")
	if info.Port != 8080 {
		t.Errorf("expected port 8080, got %d", info.Port)
	}
	if info.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", info.Proto)
	}
	if info.Host != "" {
		t.Errorf("expected empty host, got %s", info.Host)
	}
}

func TestResolvePortWithHost(t *testing.T) {
	r := &Resolver{
		lookup: func(addr string) ([]string, error) {
			return []string{"localhost."}, nil
		},
	}

	info := r.ResolvePort(443, "tcp")
	if info.Host != "localhost." {
		t.Errorf("expected host 'localhost.', got %s", info.Host)
	}
}

func TestResolvePortMultipleHosts(t *testing.T) {
	r := &Resolver{
		lookup: func(addr string) ([]string, error) {
			return []string{"first.local.", "second.local."}, nil
		},
	}

	// Only the first hostname should be used
	info := r.ResolvePort(80, "tcp")
	if info.Host != "first.local." {
		t.Errorf("expected host 'first.local.', got %s", info.Host)
	}
}

func TestProcessInfoString(t *testing.T) {
	p := ProcessInfo{Port: 80, Proto: "tcp", Address: "127.0.0.1:80", Host: "localhost."}
	s := p.String()
	if s == "" {
		t.Error("expected non-empty string")
	}

	p2 := ProcessInfo{Port: 80, Proto: "tcp", Address: "127.0.0.1:80"}
	s2 := p2.String()
	if s2 == "" {
		t.Error("expected non-empty string for no-host case")
	}
}

func TestResolveSet(t *testing.T) {
	r := &Resolver{
		lookup: func(addr string) ([]string, error) {
			return nil, fmt.Errorf("no host")
		},
	}

	ports := map[int]struct{}{
		80:   {},
		443:  {},
		8080: {},
	}

	results := r.ResolveSet(ports, "tcp")
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	seen := make(map[int]bool)
	for _, info := range results {
		seen[info.Port] = true
		if info.Proto != "tcp" {
			t.Errorf("expected proto tcp, got %s", info.Proto)
		}
	}
	for port := range ports {
		if !seen[port] {
			t.Errorf("port %d missing from results", port)
		}
	}
}

func TestResolveSetEmpty(t *testing.T) {
	r := &Resolver{
		lookup: func(addr string) ([]string, error) {
			return nil, fmt.Errorf("no host")
		},
	}

	results := r.ResolveSet(map[int]struct{}{}, "tcp")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestNewResolver(t *testing.T) {
	r := NewResolver()
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
	if r.lookup == nil {
		t.Error("expected lookup func to be set")
	}
}
