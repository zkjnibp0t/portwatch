package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/ports"
)

// Notifier is the interface implemented by all alert backends.
type Notifier interface {
	Notify(diff ports.Diff) error
}

// Daemon polls open ports on a fixed interval and fires notifications on changes.
type Daemon struct {
	cfg       *config.Config
	scanner   *ports.Scanner
	notifiers []Notifier
}

// New constructs a Daemon from the provided config, wiring up the requested notifiers.
func New(cfg *config.Config) (*Daemon, error) {
	scanner := ports.NewScanner(cfg.PortRange.Start, cfg.PortRange.End)

	var notifiers []Notifier

	if cfg.Webhook.URL != "" {
		notifiers = append(notifiers, notify.NewWebhookNotifier(cfg.Webhook.URL, cfg.Webhook.Secret))
	}

	if cfg.Desktop.Enabled {
		notifiers = append(notifiers, notify.NewDesktopNotifier(cfg.Desktop.AppName))
	}

	return &Daemon{
		cfg:       cfg,
		scanner:   scanner,
		notifiers: notifiers,
	}, nil
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch starting: range %d-%d, interval %s",
		d.cfg.PortRange.Start, d.cfg.PortRange.End, d.cfg.Interval)

	previous, err := d.scanner.Scan()
	if err != nil {
		return err
	}
	log.Printf("initial scan complete: %d open port(s)", len(previous))

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch shutting down")
			return ctx.Err()
		case <-ticker.C:
			current, err := d.scanner.Scan()
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}

			diff := ports.Compare(ports.ToSet(previous), ports.ToSet(current))
			if diff.HasChanges() {
				log.Printf("changes detected: +%d opened, -%d closed",
					len(diff.Opened), len(diff.Closed))
				d.alert(diff)
			}

			previous = current
		}
	}
}

func (d *Daemon) alert(diff ports.Diff) {
	for _, n := range d.notifiers {
		if err := n.Notify(diff); err != nil {
			log.Printf("notifier error: %v", err)
		}
	}
}
