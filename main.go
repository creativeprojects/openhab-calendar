package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creativeprojects/clog"
)

func main() {
	flag.Parse()
	closeLogger := setupLogger(flags.Verbose)
	defer closeLogger()

	// set the configuration file relative to the binary
	configFile := flags.ConfigFile
	if !path.IsAbs(configFile) {
		me, err := os.Executable()
		if err == nil {
			dir := path.Dir(me)
			configFile = path.Join(dir, configFile)
			clog.Debugf("configuration file: %s", configFile)
		}
	}

	config, err := LoadFileConfiguration(configFile)
	if err != nil {
		clog.Errorf("cannot load configuration: %v", err)
		fmt.Println("ERROR")
		return
	}

	// make sure the rules are ordered by priority
	sort.Slice(config.Rules, func(i, j int) bool {
		return config.Rules[i].Priority < config.Rules[j].Priority
	})

	// daemon mode?
	if flags.Daemon {
		startServer(config)
		return
	}
	// Legacy CLI mode

	loader := NewLoader(config)

	date := flags.Get
	if flags.Get == "" && flags.Date != "" {
		date = flags.Date
	}
	result, err := getCalendarResult(date, config, loader)
	if err != nil {
		clog.Error(err)
	}
	fmt.Println(result)
}

func getCalendarResult(dateInput string, config Configuration, loader *Loader) (string, error) {
	date, err := parseDate(dateInput)
	if err != nil {
		return "ERROR", fmt.Errorf("cannot parse date input: %w", err)
	}
	result, err := loader.GetResultFromCalendar(date, config.Rules)
	if result == "" {
		// no match: return the default
		result = config.Default.Result
		// and an error if there was none
		if err == nil {
			err = errors.New("no match")
		}
	}
	return result, err
}

func parseDate(get string) (time.Time, error) {
	get = strings.TrimSpace(get)
	if get == "" || strings.ToLower(get) == "tomorrow" {
		// default to tomorrow
		return getTomorrow(time.Now()), nil
	}
	return time.Parse(time.RFC3339, get)
}

func startServer(config Configuration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	notifyReady()

	// systemd watchdog
	go setupWatchdog()

	servers := setupServices(config)
	if len(servers) > 0 {
		// Wait until we're politely asked to leave
		<-stop
	}

	notifyLeaving()
	shutdownServices(servers)
}

func setupServices(config Configuration) map[string]*HTTPServer {
	httpServers := make(map[string]*HTTPServer, len(config.Servers))
	for name, s := range config.Servers {
		httpServer, err := NewHTTPServer(name, s)
		if err != nil {
			clog.Errorf("cannot start server %q: %v", name, err)
			continue
		}
		httpServers[name] = httpServer
		go httpServer.Start(config)
	}
	return httpServers
}

func shutdownServices(httpServers map[string]*HTTPServer) {
	if len(httpServers) == 0 {
		return
	}
	clog.Debug("shutting down...")
	var wg sync.WaitGroup
	wg.Add(len(httpServers))
	for _, s := range httpServers {
		if s == nil {
			wg.Done()
			continue
		}
		go s.Shutdown(&wg, 1*time.Minute)
	}
	wg.Wait()
}
