package ports

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiterFirstCallImmediate(t *testing.T) {
	rl := NewScanRateLimiter(10 * time.Second)
	start := time.Now()
	rl.Wait()
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("first Wait took too long: %v", elapsed)
	}
}

func TestRateLimiterEnforcesGap(t *testing.T) {
	gap := 100 * time.Millisecond
	rl := NewScanRateLimiter(gap)
	rl.Wait() // prime the last-scan time

	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)

	if elapsed < gap-5*time.Millisecond {
		t.Fatalf("second Wait returned too early: %v (want >= %v)", elapsed, gap)
	}
}

func TestRateLimiterZeroGapNeverBlocks(t *testing.T) {
	rl := NewScanRateLimiter(0)
	for i := 0; i < 5; i++ {
		start := time.Now()
		rl.Wait()
		if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
			t.Fatalf("Wait blocked with zero gap: %v", elapsed)
		}
	}
}

func TestRateLimiterResetAllowsImmediate(t *testing.T) {
	rl := NewScanRateLimiter(10 * time.Second)
	rl.Wait() // prime
	rl.Reset()

	start := time.Now()
	rl.Wait()
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("Wait after Reset blocked: %v", elapsed)
	}
}

func TestRateLimiterLastScanUpdated(t *testing.T) {
	rl := NewScanRateLimiter(0)
	if !rl.LastScan().IsZero() {
		t.Fatal("LastScan should be zero before first Wait")
	}
	before := time.Now()
	rl.Wait()
	after := time.Now()
	ls := rl.LastScan()
	if ls.Before(before) || ls.After(after) {
		t.Fatalf("LastScan %v not in expected range [%v, %v]", ls, before, after)
	}
}

func TestRateLimiterConcurrentSafe(t *testing.T) {
	rl := NewScanRateLimiter(0)
	var count int64
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			rl.Wait()
			atomic.AddInt64(&count, 1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if atomic.LoadInt64(&count) != 10 {
		t.Fatalf("expected 10 completions, got %d", count)
	}
}
