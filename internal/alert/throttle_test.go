package alert_test

import (
	"testing"
	"time"

	"github.com/yourusername/portwatch/internal/alert"
)

func TestAllowFirstAlert(t *testing.T) {
	th := alert.NewThrottler(5 * time.Second)
	if !th.Allow("port:8080") {
		t.Error("expected first alert to be allowed")
	}
}

func TestAllowSuppressesWithinCooldown(t *testing.T) {
	th := alert.NewThrottler(5 * time.Second)
	th.Allow("port:8080") // first call sets the timestamp
	if th.Allow("port:8080") {
		t.Error("expected second alert within cooldown to be suppressed")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	th := alert.NewThrottler(5 * time.Second)
	th.Allow("port:8080")
	if !th.Allow("port:9090") {
		t.Error("expected different key to be allowed independently")
	}
}

func TestAllowAfterCooldownExpires(t *testing.T) {
	th := alert.NewThrottler(10 * time.Millisecond)
	th.Allow("port:8080")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("port:8080") {
		t.Error("expected alert to be allowed after cooldown expires")
	}
}

func TestResetAllowsImmediately(t *testing.T) {
	th := alert.NewThrottler(1 * time.Hour)
	th.Allow("port:8080")
	th.Reset("port:8080")
	if !th.Allow("port:8080") {
		t.Error("expected alert to be allowed after reset")
	}
}

func TestResetAllClearsEverything(t *testing.T) {
	th := alert.NewThrottler(1 * time.Hour)
	th.Allow("port:8080")
	th.Allow("port:9090")
	th.ResetAll()
	if !th.Allow("port:8080") {
		t.Error("expected port:8080 to be allowed after ResetAll")
	}
	if !th.Allow("port:9090") {
		t.Error("expected port:9090 to be allowed after ResetAll")
	}
}
