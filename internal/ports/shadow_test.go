package ports

import (
	"testing"
	"time"
)

func fixedShadowClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestShadowNotDetectedBelowWindow(t *testing.T) {
	now := time.Now()
	st := NewShadowTracker(2 * time.Second)
	st.clock = fixedShadowClock(now)
	st.RecordOpened(8080)
	st.clock = fixedShadowClock(now.Add(3 * time.Second))
	_, ok := st.RecordClosed(8080)
	if ok {
		t.Fatal("expected no shadow for port closed after window")
	}
}

func TestShadowDetectedWithinWindow(t *testing.T) {
	now := time.Now()
	st := NewShadowTracker(5 * time.Second)
	st.clock = fixedShadowClock(now)
	st.RecordOpened(9090)
	st.clock = fixedShadowClock(now.Add(1 * time.Second))
	e, ok := st.RecordClosed(9090)
	if !ok {
		t.Fatal("expected shadow entry")
	}
	if e.Port != 9090 {
		t.Errorf("expected port 9090, got %d", e.Port)
	}
	if e.Duration != time.Second {
		t.Errorf("expected 1s duration, got %v", e.Duration)
	}
}

func TestShadowUnknownPortIgnored(t *testing.T) {
	st := NewShadowTracker(5 * time.Second)
	_, ok := st.RecordClosed(1234)
	if ok {
		t.Fatal("expected no shadow for unknown port")
	}
}

func TestShadowsAccumulate(t *testing.T) {
	now := time.Now()
	st := NewShadowTracker(10 * time.Second)
	for _, p := range []int{80, 443, 8080} {
		st.clock = fixedShadowClock(now)
		st.RecordOpened(p)
		st.clock = fixedShadowClock(now.Add(2 * time.Second))
		st.RecordClosed(p)
	}
	if len(st.Shadows()) != 3 {
		t.Errorf("expected 3 shadows, got %d", len(st.Shadows()))
	}
}

func TestShadowResetClearsAll(t *testing.T) {
	now := time.Now()
	st := NewShadowTracker(10 * time.Second)
	st.clock = fixedShadowClock(now)
	st.RecordOpened(22)
	st.clock = fixedShadowClock(now.Add(1 * time.Second))
	st.RecordClosed(22)
	st.Reset()
	if len(st.Shadows()) != 0 {
		t.Error("expected empty shadows after reset")
	}
}
