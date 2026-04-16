package ports

import (
	"strings"
	"testing"
)

// helpers

func makeResolved(port int, pid int, name string) ResolvedPort {
	return ResolvedPort{
		Port: port,
		Info: ProcessInfo{PID: pid, Name: name},
	}
}

func TestAnomalyString(t *testing.T) {
	a := Anomaly{Port: 8080, Kind: AnomalyWhitelistDenied, Detail: "not allowed"}
	s := a.String()
	if !strings.Contains(s, "whitelist_denied") {
		t.Errorf("expected kind in string, got %q", s)
	}
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port in string, got %q", s)
	}
}

func TestDetectNoAnomalies(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{{Port: 80, Process: "nginx"}})
	d := NewAnomalyDetector(wl, nil)
	resolved := []ResolvedPort{makeResolved(80, 100, "nginx")}
	anomalies := d.Detect(resolved)
	if len(anomalies) != 0 {
		t.Errorf("expected no anomalies, got %d", len(anomalies))
	}
}

func TestDetectWhitelistDenied(t *testing.T) {
	wl := NewWhitelist([]WhitelistEntry{{Port: 443, Process: "nginx"}})
	d := NewAnomalyDetector(wl, nil)
	resolved := []ResolvedPort{makeResolved(8080, 200, "rogue")}
	anomalies := d.Detect(resolved)
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Kind != AnomalyWhitelistDenied {
		t.Errorf("expected whitelist_denied, got %s", anomalies[0].Kind)
	}
}

func TestDetectUnknownProcess(t *testing.T) {
	d := NewAnomalyDetector(nil, nil)
	resolved := []ResolvedPort{makeResolved(9000, 0, "")}
	anomalies := d.Detect(resolved)
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Kind != AnomalyUnknownProcess {
		t.Errorf("expected unknown_process, got %s", anomalies[0].Kind)
	}
}

func TestDetectBaselineDrift(t *testing.T) {
	bm := NewBaselineManager("")
	// Record a baseline with only port 80.
	baseline := PortSet{80: struct{}{}}
	bm.Record(baseline)

	d := NewAnomalyDetector(nil, bm)
	// Now port 9999 is new — not in baseline.
	resolved := []ResolvedPort{
		makeResolved(80, 1, "nginx"),
		makeResolved(9999, 2, "mystery"),
	}
	anomalies := d.Detect(resolved)
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Kind != AnomalyBaselineDrift {
		t.Errorf("expected baseline_drift, got %s", anomalies[0].Kind)
	}
	if anomalies[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", anomalies[0].Port)
	}
}

func TestDetectNilWhitelistAndBaseline(t *testing.T) {
	d := NewAnomalyDetector(nil, nil)
	resolved := []ResolvedPort{makeResolved(80, 1, "nginx")}
	// Should not panic; known process, no whitelist check.
	anomalies := d.Detect(resolved)
	if len(anomalies) != 0 {
		t.Errorf("expected 0 anomalies, got %d", len(anomalies))
	}
}
