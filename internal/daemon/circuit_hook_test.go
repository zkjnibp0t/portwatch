package daemon

import (
	"bytes"
	"errors"
	"log"
	"testing"
	"time"
)

func newTestCircuitHook(threshold int) (*CircuitHook, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	h := NewCircuitHook(threshold, 500*time.Millisecond, logger)
	return h, buf
}

func TestCircuitHookAllowsInitially(t *testing.T) {
	h, _ := newTestCircuitHook(3)
	if !h.Allow() {
		t.Fatal("expected Allow() true initially")
	}
}

func TestCircuitHookBlocksAfterThreshold(t *testing.T) {
	h, buf := newTestCircuitHook(2)
	h.RecordFailure(errors.New("timeout"))
	h.RecordFailure(errors.New("timeout"))
	if h.Allow() {
		t.Fatal("expected Allow() false after threshold")
	}
	if buf.Len() == 0 {
		t.Fatal("expected log output when circuit opens")
	}
}

func TestCircuitHookLogsBlockedScan(t *testing.T) {
	h, buf := newTestCircuitHook(1)
	h.RecordFailure(errors.New("err"))
	buf.Reset()
	h.Allow()
	if buf.Len() == 0 {
		t.Fatal("expected blocked message in log")
	}
}

func TestCircuitHookRecoveryLogged(t *testing.T) {
	h, buf := newTestCircuitHook(1)
	h.RecordFailure(errors.New("err"))
	buf.Reset()
	h.RecordSuccess()
	if buf.Len() == 0 {
		t.Fatal("expected recovery log message")
	}
	if h.State() != "closed" {
		t.Fatalf("expected closed, got %s", h.State())
	}
}

func TestCircuitHookDefaultLogger(t *testing.T) {
	h := NewCircuitHook(3, time.Second, nil)
	if h.logger == nil {
		t.Fatal("expected default logger to be set")
	}
}

func TestCircuitHookStateString(t *testing.T) {
	h, _ := newTestCircuitHook(3)
	if h.State() != "closed" {
		t.Fatalf("expected closed, got %s", h.State())
	}
}
