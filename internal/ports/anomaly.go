package ports

import "fmt"

// AnomalyKind classifies the type of anomaly detected on a port.
type AnomalyKind string

const (
	AnomalyUnknownProcess  AnomalyKind = "unknown_process"
	AnomalyWhitelistDenied AnomalyKind = "whitelist_denied"
	AnomalyBaselineDrift   AnomalyKind = "baseline_drift"
)

// Anomaly describes a suspicious port event.
type Anomaly struct {
	Port    int
	Kind    AnomalyKind
	Process ProcessInfo
	Detail  string
}

func (a Anomaly) String() string {
	return fmt.Sprintf("[%s] port %d process=%s detail=%q",
		a.Kind, a.Port, a.Process.String(), a.Detail)
}

// AnomalyDetector inspects a resolved port set and emits anomalies.
type AnomalyDetector struct {
	whitelist *Whitelist
	baseline  *BaselineManager
}

// NewAnomalyDetector creates a detector backed by the given whitelist and
// baseline manager. Either may be nil, in which case that check is skipped.
func NewAnomalyDetector(wl *Whitelist, bm *BaselineManager) *AnomalyDetector {
	return &AnomalyDetector{whitelist: wl, baseline: bm}
}

// Detect evaluates each entry in resolved and returns any anomalies found.
func (d *AnomalyDetector) Detect(resolved []ResolvedPort) []Anomaly {
	var anomalies []Anomaly

	for _, rp := range resolved {
		// Whitelist check.
		if d.whitelist != nil && !d.whitelist.Allow(rp.Port, rp.Info.Name) {
			anomalies = append(anomalies, Anomaly{
				Port:    rp.Port,
				Kind:    AnomalyWhitelistDenied,
				Process: rp.Info,
				Detail:  "port/process combination not in whitelist",
			})
			continue
		}

		// Unknown process check.
		if rp.Info.PID == 0 && rp.Info.Name == "" {
			anomalies = append(anomalies, Anomaly{
				Port:    rp.Port,
				Kind:    AnomalyUnknownProcess,
				Process: rp.Info,
				Detail:  "no owning process could be resolved",
			})
		}
	}

	// Baseline drift check.
	if d.baseline != nil {
		current := make(PortSet)
		for _, rp := range resolved {
			current[rp.Port] = struct{}{}
		}
		diff := d.baseline.Diff(current)
		for _, p := range diff.Opened {
			anomalies = append(anomalies, Anomaly{
				Port:   p,
				Kind:   AnomalyBaselineDrift,
				Detail: "port opened since baseline was recorded",
			})
		}
	}

	return anomalies
}
