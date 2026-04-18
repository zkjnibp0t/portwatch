package ports

import (
	"testing"
	"time"
)

func fixedEvictionClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestEvictionNotEvictedBeforeTTL(t *testing.T) {
	now := time.Now()
	ep := NewEvictionPolicy(5 * time.Second)
	ep.clock = fixedEvictionClock(now)
	ep.Touch(map[int]struct{}{80: {}, 443: {}})
	ep.clock = fixedEvictionClock(now.Add(3 * time.Second))
	evicted := ep.Evict()
	if len(evicted) != 0 {
		t.Fatalf("expected no evictions, got %v", evicted)
	}
}

func TestEvictionEvictsAfterTTL(t *testing.T) {
	now := time.Now()
	ep := NewEvictionPolicy(5 * time.Second)
	ep.clock = fixedEvictionClock(now)
	ep.Touch(map[int]struct{}{8080: {}})
	ep.clock = fixedEvictionClock(now.Add(6 * time.Second))
	evicted := ep.Evict()
	if len(evicted) != 1 || evicted[0] != 8080 {
		t.Fatalf("expected port 8080 evicted, got %v", evicted)
	}
	if ep.Len() != 0 {
		t.Fatalf("expected empty tracker after eviction")
	}
}

func TestEvictionTouchRefreshes(t *testing.T) {
	now := time.Now()
	ep := NewEvictionPolicy(5 * time.Second)
	ep.clock = fixedEvictionClock(now)
	ep.Touch(map[int]struct{}{22: {}})
	ep.clock = fixedEvictionClock(now.Add(4 * time.Second))
	ep.Touch(map[int]struct{}{22: {}})
	ep.clock = fixedEvictionClock(now.Add(8 * time.Second))
	evicted := ep.Evict()
	if len(evicted) != 0 {
		t.Fatalf("expected no evictions after refresh, got %v", evicted)
	}
}

func TestEvictionResetClearsAll(t *testing.T) {
	ep := NewEvictionPolicy(10 * time.Second)
	ep.Touch(map[int]struct{}{1: {}, 2: {}, 3: {}})
	ep.Reset()
	if ep.Len() != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", ep.Len())
	}
}

func TestEvictionLenTracksCount(t *testing.T) {
	ep := NewEvictionPolicy(10 * time.Second)
	ep.Touch(map[int]struct{}{100: {}, 200: {}, 300: {}})
	if ep.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", ep.Len())
	}
}
