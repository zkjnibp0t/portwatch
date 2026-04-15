package notify

import "fmt"

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Notify(opened, closed []string) error
}

// MultiNotifier fans out notifications to multiple backends.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier from the provided notifiers.
// Nil entries are silently ignored.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	filtered := make([]Notifier, 0, len(notifiers))
	for _, n := range notifiers {
		if n != nil {
			filtered = append(filtered, n)
		}
	}
	return &MultiNotifier{notifiers: filtered}
}

// Notify calls every registered notifier and collects errors.
// All notifiers are attempted even if earlier ones fail.
func (m *MultiNotifier) Notify(opened, closed []string) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(opened, closed); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	msg := fmt.Sprintf("%d notifier(s) failed:", len(errs))
	for i, e := range errs {
		msg += fmt.Sprintf(" [%d] %s", i+1, e.Error())
	}
	return fmt.Errorf("%s", msg)
}

// Len returns the number of registered notifiers.
func (m *MultiNotifier) Len() int {
	return len(m.notifiers)
}
