package config

import (
	"fmt"

	"github.com/user/portwatch/internal/ports"
)

// PrivilegedContributor raises score for ports below 1024.
type PrivilegedContributor struct{ Delta float64 }

func (p PrivilegedContributor) Contribute(port int) (float64, string) {
	if port < 1024 {
		return p.Delta, fmt.Sprintf("privileged port %d", port)
	}
	return 0, ""
}

// WellKnownRiskyContributor raises score for commonly exploited ports.
type WellKnownRiskyContributor struct{ Delta float64 }

var riskyPorts = map[int]string{
	21:   "ftp",
	23:   "telnet",
	445:  "smb",
	3389: "rdp",
	5900: "vnc",
}

func (w WellKnownRiskyContributor) Contribute(port int) (float64, string) {
	if name, ok := riskyPorts[port]; ok {
		return w.Delta, fmt.Sprintf("risky service %s", name)
	}
	return 0, ""
}

// BuildScorer constructs a default Scorer from config knobs.
func BuildScorer(privilegedDelta, riskyDelta float64) *ports.Scorer {
	return ports.NewScorer(
		PrivilegedContributor{Delta: privilegedDelta},
		WellKnownRiskyContributor{Delta: riskyDelta},
	)
}
