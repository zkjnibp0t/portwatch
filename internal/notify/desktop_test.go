package notify

import (
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func TestNewDesktopNotifierDefaults(t *testing.T) {
	d := NewDesktopNotifier("")
	if d.AppName != "portwatch" {
		t.Errorf("expected default app name 'portwatch', got %q", d.AppName)
	}
}

func TestNewDesktopNotifierCustomName(t *testing.T) {
	d := NewDesktopNotifier("myapp")
	if d.AppName != "myapp" {
		t.Errorf("expected app name 'myapp', got %q", d.AppName)
	}
}

func TestBuildMessageNoDiff(t *testing.T) {
	diff := ports.Diff{}
	title, body := buildMessage(diff)
	if title != "" || body != "" {
		t.Errorf("expected empty title and body for empty diff, got %q / %q", title, body)
	}
}

func TestBuildMessageOpenedOnly(t *testing.T) {
	diff := ports.Diff{
		Opened: []string{"tcp:8080", "tcp:9090"},
	}
	title, body := buildMessage(diff)
	if title == "" {
		t.Error("expected non-empty title")
	}
	for _, p := range diff.Opened {
		if !containsStr(body, p) {
			t.Errorf("expected body to contain %q", p)
		}
	}
}

func TestBuildMessageClosedOnly(t *testing.T) {
	diff := ports.Diff{
		Closed: []string{"tcp:3306"},
	}
	title, body := buildMessage(diff)
	if title == "" {
		t.Error("expected non-empty title")
	}
	if !containsStr(body, "tcp:3306") {
		t.Errorf("expected body to contain closed port, got %q", body)
	}
}

func TestBuildMessageBothOpenedAndClosed(t *testing.T) {
	diff := ports.Diff{
		Opened: []string{"tcp:8080"},
		Closed: []string{"tcp:22"},
	}
	title, _ := buildMessage(diff)
	if !containsStr(title, "opened") || !containsStr(title, "closed") {
		t.Errorf("expected title to mention both opened and closed, got %q", title)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
