package config

import "testing"

func TestBuildWhitelistEmpty(t *testing.T) {
	wl := BuildWhitelist(nil)
	if wl == nil {
		t.Fatal("expected non-nil whitelist")
	}
	// empty whitelist denies everything
	if wl.Allow(80, "httpd") {
		t.Error("empty whitelist should not allow anything")
	}
}

func TestBuildWhitelistSingleEntry(t *testing.T) {
	entries := []WhitelistEntryConfig{
		{Port: 8080, Process: "nginx"},
	}
	wl := BuildWhitelist(entries)
	if !wl.Allow(8080, "nginx") {
		t.Error("expected 8080/nginx to be allowed")
	}
	if wl.Allow(8080, "other") {
		t.Error("expected 8080/other to be denied")
	}
}

func TestBuildWhitelistPortWildcard(t *testing.T) {
	entries := []WhitelistEntryConfig{
		{Port: 443, Process: ""},
	}
	wl := BuildWhitelist(entries)
	if !wl.Allow(443, "caddy") {
		t.Error("expected 443 with any process to be allowed")
	}
}

func TestBuildWhitelistProcessWildcard(t *testing.T) {
	entries := []WhitelistEntryConfig{
		{Port: 0, Process: "sshd"},
	}
	wl := BuildWhitelist(entries)
	if !wl.Allow(22, "sshd") {
		t.Error("expected sshd on any port to be allowed")
	}
	if wl.Allow(22, "other") {
		t.Error("expected non-sshd process to be denied")
	}
}

func TestBuildWhitelistMultipleEntries(t *testing.T) {
	entries := []WhitelistEntryConfig{
		{Port: 80, Process: "apache"},
		{Port: 443, Process: ""},
	}
	wl := BuildWhitelist(entries)
	if !wl.Allow(80, "apache") {
		t.Error("expected 80/apache allowed")
	}
	if !wl.Allow(443, "anything") {
		t.Error("expected 443/anything allowed")
	}
	if wl.Allow(8080, "unknown") {
		t.Error("expected 8080/unknown denied")
	}
}
