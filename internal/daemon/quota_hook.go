package daemon

import (
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// QuotaHook gates each scan cycle through a QuotaTracker and logs when
// the quota is exhausted.
type QuotaHook struct {
	quota  *ports.QuotaTracker
	logger *log.Logger
}

// NewQuotaHook creates a QuotaHook wrapping the given tracker.
func NewQuotaHook(q *ports.QuotaTracker, w io.Writer) *QuotaHook {
	if w == nil {
		w = os.Stdout
	}
	return &QuotaHook{
		quota:  q,
		logger: log.New(w, "[quota] ", 0),
	}
}

// BeforeScan returns false (blocking the scan) when the quota is exceeded.
func (h *QuotaHook) BeforeScan() bool {
	if err := h.quota.Allow(); err != nil {
		h.logger.Printf("scan blocked: %v (remaining: %d)", err, h.quota.Remaining())
		return false
	}
	return true
}

// AfterScan is a no-op for this hook but satisfies a potential hook interface.
func (h *QuotaHook) AfterScan() {}
