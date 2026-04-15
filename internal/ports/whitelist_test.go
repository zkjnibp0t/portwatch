package ports

import (
	"testing"
)

func TestWhitelistAllowExactMatch(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 8080, Process: "nginx"},
	})
	if !wl.Allow(8080, "nginx") {
		t.Error("expected exact match to be allowed")
	}
}

func TestWhitelistAllowPortOnly(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 443, Process: ""},
	})
	if !wl.Allow(443, "anything") {
		t.Error("expected port-only rule to allow any process")
	}
}

func TestWhitelistAllowProcessOnly(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 0, Process: "sshd"},
	})
	if !wl.Allow(22, "sshd") {
		t.Error("expected process-only rule to allow any port")
	}
}

func TestWhitelistDenyUnknown(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 8080, Process: "nginx"},
	})
	if wl.Allow(9090, "unknown") {
		t.Error("expected unknown entry to be denied")
	}
}

func TestWhitelistEmpty(t *testing.T) {
	wl := NewWhitelist(nil)
	if wl.Allow(80, "httpd") {
		t.Error("empty whitelist should deny everything")
	}
}

func TestWhitelistFilterDiffRemovesAllowed(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 8080, Process: "nginx"},
	})
	diff := Diff{Opened: []int{8080, 9090}, Closed: []int{22}}
	lookup := func(p int) string {
		if p == 8080 {
			return "nginx"
		}
		return "unknown"
	}
	result := wl.FilterDiff(diff, lookup)
	if len(result.Opened) != 1 || result.Opened[0] != 9090 {
		t.Errorf("expected only 9090 in opened, got %v", result.Opened)
	}
	if len(result.Closed) != 1 || result.Closed[0] != 22 {
		t.Errorf("closed should be unchanged, got %v", result.Closed)
	}
}

func TestWhitelistFilterDiffEmptyWhitelist(t *testing.T) {
	wl := NewWhitelist(nil)
	diff := Diff{Opened: []int{80, 443}, Closed: []int{}}
	result := wl.FilterDiff(diff, nil)
	if len(result.Opened) != 2 {
		t.Errorf("expected all ports to pass through, got %v", result.Opened)
	}
}

func TestWhitelistFilterDiffNilLookup(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{
		{Port: 8080, Process: ""},
	})
	diff := Diff{Opened: []int{8080, 9090}}
	result := wl.FilterDiff(diff, nil)
	if len(result.Opened) != 1 || result.Opened[0] != 9090 {
		t.Errorf("expected 8080 filtered out, got %v", result.Opened)
	}
}
