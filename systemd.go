package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/creativeprojects/clog"
)

func notifyReady() {
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		clog.Errorf("cannot notify systemd: %s", err)
	}
}

func notifyLeaving() {
	_, _ = daemon.SdNotify(false, daemon.SdNotifyStopping)
}

func setupWatchdog(servers map[string]*HTTPServer) {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		clog.Errorf("cannot verify if systemd watchdog is enabled: %s", err)
		return
	}
	if interval == 0 {
		// watchdog not enabled
		return
	}
	for {
		// Check that the services are healthy
		healthy := true
		for _, server := range servers {
			if server == nil {
				continue
			}
			parts := strings.Split(server.listen, ":")
			if len(parts) != 2 {
				clog.Errorf("invalid listen address %q", server.listen)
				continue
			}
			url := "http"
			if server.tls {
				url = "https"
			}
			url += "://localhost:" + parts[1] + "/health"
			result, err := http.Get(url)
			if err != nil {
				clog.Errorf("cannot verify if service %q is healthy: %s", server.name, err)
				healthy = false
				continue
			}
			if result.StatusCode != http.StatusOK {
				clog.Errorf("/health endpoint on service %q returned HTTP %d", server.name, result.StatusCode)
				healthy = false
				continue
			}
		}
		if healthy {
			_, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			if err != nil {
				clog.Errorf("cannot notify systemd watchdog: %s", err)
			}
		}
		time.Sleep(interval / 3)
	}
}
