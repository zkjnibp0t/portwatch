package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/ports"
)

func TestWebhookNotifierSuccess(t *testing.T) {
	var received notify.WebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := notify.NewWebhookNotifier(server.URL)
	diff := ports.Diff{
		Opened: []ports.Port{{Number: 8080, Protocol: "tcp"}},
		Closed: []ports.Port{{Number: 22, Protocol: "tcp"}},
	}

	if err := notifier.Notify(diff); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	if len(received.Opened) != 1 || received.Opened[0].Number != 8080 {
		t.Errorf("unexpected opened ports: %v", received.Opened)
	}
	if len(received.Closed) != 1 || received.Closed[0].Number != 22 {
		t.Errorf("unexpected closed ports: %v", received.Closed)
	}
	if received.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
}

func TestWebhookNotifierNon2xxError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := notify.NewWebhookNotifier(server.URL)
	err := notifier.Notify(ports.Diff{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestWebhookNotifierInvalidURL(t *testing.T) {
	notifier := notify.NewWebhookNotifier("http://127.0.0.1:0/no-server")
	err := notifier.Notify(ports.Diff{})
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
