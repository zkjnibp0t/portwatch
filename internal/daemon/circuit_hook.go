package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// CircuitHook wraps a CircuitBreaker and integrates it into the daemon scan cycle.
// It records scan health and logs when the circuit opens or recovers.
type CircuitHook struct {
	cb     *ports.CircuitBreaker
	logger *log.Logger
}

// NewCircuitHook creates a CircuitHook with the given failure threshold and reset window.
func NewCircuitHook(threshold int, resetAfter time.Duration, logger *log.Logger) *CircuitHook {
	if logger == nil {
		logger = log.Default()
	}
	return &CircuitHook{
		cb:     ports.NewCircuitBreaker(threshold, resetAfter),
		logger: logger,
	}
}

// Allow returns true when the circuit permits a scan to proceed.
func (h *CircuitHook) Allow() bool {
	if !h.cb.Allow() {
		h.logger.Println("[circuit] scan blocked: circuit is open")
		return false
	}
	return true
}

// RecordSuccess notifies the circuit that a scan completed successfully.
func (h *CircuitHook) RecordSuccess() {
	prev := h.cb.State()
	h.cb.RecordSuccess()
	if prev != ports.CircuitClosed {
		h.logger.Println("[circuit] circuit recovered: state -> closed")
	}
}

// RecordFailure notifies the circuit that a scan failed.
func (h *CircuitHook) RecordFailure(err error) {
	h.cb.RecordFailure()
	state := h.cb.State()
	if state == ports.CircuitOpen {
		h.logger.Printf("[circuit] circuit opened after error: %v", err)
	}
}

// State returns the current circuit state string for diagnostics.
func (h *CircuitHook) State() string {
	return h.cb.State().String()
}
