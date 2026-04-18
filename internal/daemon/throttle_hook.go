package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// ThrottleHook skips scan cycles when the port fingerprint is unchanged within
// the configured minimum interval.
type ThrottleHook struct {
	throttle     *ports.ScanThrottle
	fingerprinter *ports.Fingerprinter
	logger       *log.Logger
}

// NewThrottleHook creates a ThrottleHook with the given minimum interval.
func NewThrottleHook(minInterval time.Duration, logger *log.Logger) *ThrottleHook {
	if logger == nil {
		logger = log.Default()
	}
	return &ThrottleHook{
		throttle:      ports.NewScanThrottle(minInterval),
		fingerprinter: ports.NewFingerprinter(),
		logger:        logger,
	}
}

// BeforeScan returns false (skip) when the fingerprint is unchanged within the
// throttle window, true otherwise.
func (h *ThrottleHook) BeforeScan(current ports.PortSet) bool {
	fp := h.fingerprinter.Compute(current)
	if h.throttle.ShouldSkip(fp) {
		h.logger.Println("[throttle] skipping scan: fingerprint unchanged within interval")
		return false
	}
	return true
}

// AfterScan records the fingerprint of the completed scan.
func (h *ThrottleHook) AfterScan(current ports.PortSet) {
	fp := h.fingerprinter.Compute(current)
	h.throttle.Record(fp)
}

// Reset clears throttle state, forcing the next scan to proceed.
func (h *ThrottleHook) Reset() {
	h.throttle.Reset()
}
