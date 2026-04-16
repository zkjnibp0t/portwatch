package ports

import (
	"sync"
	"testing"
	"time"
)

func TestDebouncerFiresAfterQuiet(t *testing.T) {
	var mu sync.Mutex
	fired := map[int]bool{}

	d := NewDebouncer(30*time.Millisecond, func(port int, opened bool) {
		mu.Lock()
		fired[port] = opened
		mu.Unlock()
	})

	d.Push(8080, true)
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if !fired[8080] {
		t.Fatal("expected callback to fire for port 8080")
	}
}

func TestDebouncerResetsOnRapidPush(t *testing.T) {
	count := 0
	var mu sync.Mutex

	d := NewDebouncer(40*time.Millisecond, func(port int, opened bool) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Push three times rapidly; only one callback should fire.
	d.Push(9000, true)
	d.Push(9000, false)
	d.Push(9000, true)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 callback, got %d", count)
	}
}

func TestDebouncerFlushCancelsPending(t *testing.T) {
	fired := false

	d := NewDebouncer(50*time.Millisecond, func(port int, opened bool) {
		fired = true
	})

	d.Push(1234, true)
	if d.Pending() != 1 {
		t.Fatal("expected 1 pending timer")
	}
	d.Flush()
	if d.Pending() != 0 {
		t.Fatal("expected 0 pending timers after flush")
	}

	time.Sleep(80 * time.Millisecond)
	if fired {
		t.Fatal("callback should not fire after flush")
	}
}

func TestDebouncerMultiplePorts(t *testing.T) {
	var mu sync.Mutex
	results := map[int]bool{}

	d := NewDebouncer(20*time.Millisecond, func(port int, opened bool) {
		mu.Lock()
		results[port] = opened
		mu.Unlock()
	})

	d.Push(80, true)
	d.Push(443, false)

	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if !results[80] {
		t.Error("expected port 80 opened=true")
	}
	if results[443] {
		t.Error("expected port 443 opened=false")
	}
}
