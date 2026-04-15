package ports

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcessInfo holds metadata about the process owning a port.
type ProcessInfo struct {
	PID  int
	Name string
	User string
}

func (p ProcessInfo) String() string {
	if p.PID == 0 {
		return "unknown"
	}
	parts := []string{strconv.Itoa(p.PID)}
	if p.Name != "" {
		parts = append(parts, p.Name)
	}
	if p.User != "" {
		parts = append(parts, "("+p.User+")")
	}
	return strings.Join(parts, " ")
}

// LookupProcess attempts to find the process name for a given PID by reading
// /proc/<pid>/comm on Linux. Returns an empty ProcessInfo when unavailable.
func LookupProcess(pid int) ProcessInfo {
	if pid <= 0 {
		return ProcessInfo{}
	}
	info := ProcessInfo{PID: pid}

	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	data, err := os.ReadFile(commPath)
	if err == nil {
		info.Name = strings.TrimSpace(string(data))
	}

	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	statusData, err := os.ReadFile(statusPath)
	if err == nil {
		for _, line := range strings.Split(string(statusData), "\n") {
			if strings.HasPrefix(line, "Uid:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					info.User = resolveUID(fields[1])
				}
				break
			}
		}
	}
	return info
}

// resolveUID returns the username for a UID string, falling back to the UID itself.
func resolveUID(uid string) string {
	// Avoid importing os/user to keep dependencies minimal; return raw UID.
	return uid
}
