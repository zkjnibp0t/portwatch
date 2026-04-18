package ports

import (
	"testing"
	"time"
)

func fixedThrottleClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestThrottleAllowsFirstScan(t *testing.T) {
	th := NewScanThrottle(10 * time.Second)
	if th.ShouldSkip("abc") {
		t.Fatal("expected first scan to be allowed")
	}
}

func TestThrottleSkipsWithinInterval(t *testing.T) {
	now := time.Now()
	th := NewScanThrottle(30 * time.Second)
	th.clock = fixedThrottleClock(now)
	th.Record("abc")
	if !th.ShouldSkip("abc") {
		t.Fatal("expected scan to be skipped within interval")
	}
}

func TestThrottleAllowsAfterInterval(t *testing.T) {
	now := time.Now()
	th := NewScanThrottle(5 * time.Second)
	th.clock = fixedThrottleClock(now)
	th.Record("abc")
	th.clock = fixedThrottleClock(now.Add(10 * time.Second))
	if th.ShouldSkip("abc") {
		t.Fatal("expected scan to be allowed after interval")
	}
}

func TestThrottleAllowsOnFingerprintChange(t *testing.T) {
	now := time.Now()
	th := NewScanThrottle(60 * time.Second)
	th.clock = fixedThrottleClock(now)
	th.Record("abc")
	if th.ShouldSkip("xyz") {
		t.Fatal("expected scan to be allowed when fingerprint changes")
	}
}

func TestThrottleResetAllowsImmediate(t *testing.T) {
	now := time.Now()
	th := NewScanThrottle(60 * time.Second)
	th.clock = fixedThrottleClock(now)
	th.Record("abc")
	th.Reset()
	if th.ShouldSkip("abc") {
		t.Fatal("expected scan to be allowed after reset")
	}
}

func TestThrottleLastScanUpdated(t *testing.T) {
	now := time.Now()
	th := NewScanThrottle(10 * time.Second)
	th.clock = fixedThrottleClock(now)
	th.Record("abc")
	if !th.LastScan().Equal(now) {
		t.Fatalf("expected last scan %v, got %v", now, th.LastScan())
	}
}
