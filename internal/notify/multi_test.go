package notify

import (
	"errors"
	"testing"
)

// stubNotifier records the last call made to it.
type stubNotifier struct {
	called bool
	opened []string
	closed []string
	errToReturn error
}

func (s *stubNotifier) Notify(opened, closed []string) error {
	s.called = true
	s.opened = opened
	s.closed = closed
	return s.errToReturn
}

func TestMultiNotifierCallsAll(t *testing.T) {
	a := &stubNotifier{}
	b := &stubNotifier{}
	m := NewMultiNotifier(a, b)

	if err := m.Notify([]string{"8080"}, []string{"22"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both notifiers to be called")
	}
}

func TestMultiNotifierSkipsNil(t *testing.T) {
	a := &stubNotifier{}
	m := NewMultiNotifier(nil, a, nil)
	if m.Len() != 1 {
		t.Fatalf("expected 1 notifier, got %d", m.Len())
	}
}

func TestMultiNotifierSingleError(t *testing.T) {
	a := &stubNotifier{errToReturn: errors.New("boom")}
	m := NewMultiNotifier(a)
	err := m.Notify(nil, nil)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected 'boom', got %v", err)
	}
}

func TestMultiNotifierMultipleErrors(t *testing.T) {
	a := &stubNotifier{errToReturn: errors.New("err-a")}
	b := &stubNotifier{errToReturn: errors.New("err-b")}
	m := NewMultiNotifier(a, b)
	err := m.Notify(nil, nil)
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
	msg := err.Error()
	if !contains(msg, "2 notifier(s) failed") {
		t.Errorf("unexpected error message: %s", msg)
	}
}

func TestMultiNotifierNoNotifiers(t *testing.T) {
	m := NewMultiNotifier()
	if err := m.Notify([]string{"443"}, nil); err != nil {
		t.Fatalf("unexpected error with empty notifier list: %v", err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
