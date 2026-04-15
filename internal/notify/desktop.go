package notify

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/user/portwatch/internal/ports"
)

// DesktopNotifier sends desktop notifications using OS-native tools.
type DesktopNotifier struct {
	AppName string
}

// NewDesktopNotifier creates a new DesktopNotifier with the given app name.
func NewDesktopNotifier(appName string) *DesktopNotifier {
	if appName == "" {
		appName = "portwatch"
	}
	return &DesktopNotifier{AppName: appName}
}

// Notify sends a desktop notification summarising port changes.
func (d *DesktopNotifier) Notify(diff ports.Diff) error {
	title, body := buildMessage(diff)
	if title == "" {
		return nil
	}
	return sendDesktopNotification(d.AppName+": "+title, body)
}

func buildMessage(diff ports.Diff) (title, body string) {
	var parts []string
	if len(diff.Opened) > 0 {
		parts = append(parts, fmt.Sprintf("%d port(s) opened", len(diff.Opened)))
	}
	if len(diff.Closed) > 0 {
		parts = append(parts, fmt.Sprintf("%d port(s) closed", len(diff.Closed)))
	}
	if len(parts) == 0 {
		return "", ""
	}
	title = strings.Join(parts, ", ")

	var lines []string
	for _, p := range diff.Opened {
		lines = append(lines, fmt.Sprintf("+ %s", p))
	}
	for _, p := range diff.Closed {
		lines = append(lines, fmt.Sprintf("- %s", p))
	}
	body = strings.Join(lines, "\n")
	return title, body
}

func sendDesktopNotification(title, body string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("notify-send", title, body).Run()
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, body, title)
		return exec.Command("osascript", "-e", script).Run()
	case "windows":
		// PowerShell toast-style message box fallback
		ps := fmt.Sprintf(
			`[System.Windows.Forms.MessageBox]::Show('%s','%s')`,
			strings.ReplaceAll(body, "'", "`'"),
			strings.ReplaceAll(title, "'", "`'"),
		)
		return exec.Command("powershell", "-Command", ps).Run()
	default:
		return fmt.Errorf("desktop notifications not supported on %s", runtime.GOOS)
	}
}
