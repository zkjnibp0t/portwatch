package ports

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSuppressIsSuppressed(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	s.Suppress(8080, "", 5*time.Minute)
	if !s.IsSuppressed(8080, "nginx") {
		t.Fatal("expected port 8080 to be suppressed")
	}
}

func TestSuppressExpires(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	s.Suppress(9090, "", 1*time.Second)
	// advance clock past expiry
	s.now = fixedClock(now.Add(2 * time.Second))
	if s.IsSuppressed(9090, "") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestSuppressProcessSpecific(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	s.Suppress(443, "nginx", 10*time.Minute)
	if !s.IsSuppressed(443, "nginx") {
		t.Fatal("expected nginx on 443 to be suppressed")
	}
	if s.IsSuppressed(443, "sshd") {
		t.Fatal("expected sshd on 443 not to be suppressed")
	}
}

func TestSuppressUnrelatedPort(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	s.Suppress(8080, "", 5*time.Minute)
	if s.IsSuppressed(9090, "") {
		t.Fatal("expected port 9090 not to be suppressed")
	}
}

func TestSuppressPrune(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	s.Suppress(1111, "", 1*time.Second)
	s.Suppress(2222, "", 10*time.Minute)
	s.now = fixedClock(now.Add(2 * time.Second))
	s.Prune()
	if s.ActiveCount() != 1 {
		t.Fatalf("expected 1 active rule after prune, got %d", s.ActiveCount())
	}
}

func TestSuppressActiveCount(t *testing.T) {
	now := time.Now()
	s := &Suppressor{now: fixedClock(now)}
	if s.ActiveCount() != 0 {
		t.Fatal("expected 0 active rules")
	}
	s.Suppress(80, "", 5*time.Minute)
	s.Suppress(443, "", 5*time.Minute)
	if s.ActiveCount() != 2 {
		t.Fatalf("expected 2 active rules, got %d", s.ActiveCount())
	}
}
