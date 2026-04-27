package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/user/portwatch/internal/ports"
)

// ClusterHook records co-occurring port changes and logs detected clusters.
type ClusterHook struct {
	detector   *ports.ClusterDetector
	minSupport int
	log        *log.Logger
	w          io.Writer
}

// NewClusterHook returns a ClusterHook that logs port clusters once they
// reach minSupport co-occurrences.
func NewClusterHook(detector *ports.ClusterDetector, minSupport int, w io.Writer) *ClusterHook {
	if w == nil {
		w = os.Stdout
	}
	return &ClusterHook{
		detector:   detector,
		minSupport: minSupport,
		log:        log.New(w, "[cluster] ", 0),
		w:          w,
	}
}

// BeforeScan is a no-op for this hook.
func (h *ClusterHook) BeforeScan() {}

// AfterScan records opened and closed ports as a co-occurrence group and
// logs any clusters that have reached the support threshold.
func (h *ClusterHook) AfterScan(diff ports.Diff) {
	changed := append(diff.Opened, diff.Closed...)
	if len(changed) < 2 {
		return
	}
	h.detector.Record(changed)
	clusters := h.detector.Clusters()
	for _, cl := range clusters {
		if cl.Count == h.minSupport {
			h.log.Printf("cluster detected: ports [%s] co-occurred %d times",
				formatClusterPorts(cl.Ports), cl.Count)
		}
	}
}

func formatClusterPorts(ports []int) string {
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ", ")
}
