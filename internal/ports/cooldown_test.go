package ports

import (
	"testing"
	"time"
)

func fixedCooldownClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCooldownNotActiveInitially(t *testing.T) {
	c := NewCooldownTracker(5 * time.Second)
	if c.IsActive(8080) {
		t.Fatal("expected port 8080 to not be in cooldown initially")
	}
}

func TestCooldownActiveAfterRecord(t *testing.T) {
	now := time.Now()
	c := NewCooldownTracker(5 * time.Second)
	c.clock = fixedCooldownClock(now)
	c.Record(8080)
	if !c.IsActive(8080) {
		t.Fatal("expected port 8080 to be in cooldown after record")
	}
}

func TestCooldownExpiresAfterWindow(t *testing.T) {
	now := time.Now()
	c := NewCooldownTracker(5 * time.Second)
	c.clock = fixedCooldownClock(now)
	c.Record(8080)
	c.clock = fixedCooldownClock(now.Add(6 * time.Second))
	if c.IsActive(8080) {
		t.Fatal("expected cooldown to have expired")
	}
}

func TestCooldownResetClearsPort(t *testing.T) {
	now := time.Now()
	c := NewCooldownTracker(5 * time.Second)
	c.clock = fixedCooldownClock(now)
	c.Record(9090)
	c.Reset(9090)
	if c.IsActive(9090) {
		t.Fatal("expected port 9090 cooldown to be cleared after reset")
	}
}

func TestCooldownPruneRemovesExpired(t *testing.T) {
	now := time.Now()
	c := NewCooldownTracker(2 * time.Second)
	c.clock = fixedCooldownClock(now)
	c.Record(1234)
	c.Record(5678)
	c.clock = fixedCooldownClock(now.Add(3 * time.Second))
	c.Prune()
	if len(c.entries) != 0 {
		t.Fatalf("expected all entries pruned, got %d", len(c.entries))
	}
}

func TestCooldownDifferentPortsIndependent(t *testing.T) {
	now := time.Now()
	c := NewCooldownTracker(5 * time.Second)
	c.clock = fixedCooldownClock(now)
	c.Record(80)
	if c.IsActive(443) {
		t.Fatal("expected port 443 to not be affected by port 80 cooldown")
	}
}
