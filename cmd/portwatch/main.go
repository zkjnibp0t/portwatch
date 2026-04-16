package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/ports"
)

var (
	cfgPath   = flag.String("config", "portwatch.yaml", "path to config file")
	reportCmd = flag.Bool("report", false, "print scan history report and exit")
	version   = flag.Bool("version", false, "print version and exit")

	buildVersion = "dev"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("portwatch %s\n", buildVersion)
		return
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Handle report-only mode.
	if *reportCmd {
		store, err := history.NewStore(cfg.HistoryFile)
		if err != nil {
			log.Fatalf("failed to open history store: %v", err)
		}
		history.PrintReport(os.Stdout, store)
		return
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)

	// Build notifiers.
	var notifiers []notify.Notifier
	if cfg.Webhook.URL != "" {
		wh, err := notify.NewWebhookNotifier(cfg.Webhook.URL, cfg.Webhook.Secret)
		if err != nil {
			logger.Fatalf("invalid webhook config: %v", err)
		}
		notifiers = append(notifiers, wh)
	}
	if cfg.Desktop.Enabled {
		notifiers = append(notifiers, notify.NewDesktopNotifier(cfg.Desktop.AppName))
	}
	multi := notify.NewMultiNotifier(notifiers...)

	// Build alert manager with throttle.
	throttler := alert.NewThrottler(cfg.Alert.CooldownDuration())
	alertMgr := alert.NewManager(multi, throttler)

	// Build port scanner and filter.
	scanner := ports.NewScanner(cfg.Ports.Low, cfg.Ports.High)
	filter := ports.NewFilter(cfg.Ports.Include, cfg.Ports.Exclude)
	whitelist := config.BuildWhitelist(cfg.Whitelist)

	// Build daemon.
	d, err := daemon.New(daemon.Options{
		Config:    cfg,
		Scanner:   scanner,
		Filter:    filter,
		Whitelist: whitelist,
		AlertMgr:  alertMgr,
		Logger:    logger,
	})
	if err != nil {
		logger.Fatalf("failed to create daemon: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Printf("portwatch %s starting (ports %d-%d, interval %s)",
		buildVersion, cfg.Ports.Low, cfg.Ports.High, cfg.Interval)

	if err := d.Run(ctx); err != nil && err != context.Canceled {
		logger.Fatalf("daemon exited with error: %v", err)
	}
	logger.Println("portwatch stopped")
}
