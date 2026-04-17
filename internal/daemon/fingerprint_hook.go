package daemon

import (
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// FingerprintHook logs when the port-set fingerprint changes between cycles.
type FingerprintHook struct {
	fp   *ports.Fingerprinter
	last ports.Fingerprint
	log  *log.Logger
}

// NewFingerprintHook returns a FingerprintHook using the given writer.
// If w is nil, os.Stdout is used.
func NewFingerprintHook(w io.Writer) *FingerprintHook {
	if w == nil {
		w = os.Stdout
	}
	return &FingerprintHook{
		fp:  ports.NewFingerprinter(),
		log: log.New(w, "[fingerprint] ", 0),
	}
}

// BeforeScan is a no-op.
func (h *FingerprintHook) BeforeScan() {}

// AfterScan computes the fingerprint for the new port set and logs if it changed.
func (h *FingerprintHook) AfterScan(resolved []ports.ResolvedPort) {
	current := h.fp.Compute(resolved)
	if h.last == "" {
		h.last = current
		h.log.Printf("initial fingerprint: %s", current)
		return
	}
	if !h.fp.Equal(h.last, current) {
		h.log.Printf("fingerprint changed: %s -> %s", h.last, current)
		h.last = current
	}
}

// Current returns the most recently recorded fingerprint.
func (h *FingerprintHook) Current() ports.Fingerprint {
	return h.last
}
