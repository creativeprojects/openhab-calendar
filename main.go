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

	// make sure the post-rules are ordered by priority
	sort.Slice(config.PostRules, func(i, j int) bool {
		return config.PostRules[i].Priority < config.PostRules[j].Priority
	})

	// daemon mode?
	if flags.Daemon {
		startServer(config)
		return
	}

	// save mode?
	if flags.Save {
		saveCalendars(config)
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

func getCalendarResult(dateInput string, config Configuration, loader *Loader) (Result, error) {
	date, err := parseDate(dateInput)
	if err != nil {
		return ResultError, fmt.Errorf("cannot parse date input: %w", err)
	}
	result, err := loader.GetResultFromRules(date, config.Rules)
	if result.Calendar == "" {
		// no match: return the default
		result.Calendar = config.Default.Result
		// and an error if there was none
		if err == nil {
			err = errors.New("no match")
		}
		return result, err
	}
	// run the post-rules
	result, err = loader.PostRules(result, date, config.Rules, config.PostRules)
	return result, err
}

func parseDate(get string) (time.Time, error) {
	get = strings.TrimSpace(get)
	if get == "" || strings.ToLower(get) == "tomorrow" {
		// default to tomorrow
		return getTomorrow(time.Now()), nil
	}
	if strings.ToLower(get) == "today" {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, get)
}

func startServer(config Configuration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	notifyReady()

	servers := setupServices(config)
	if len(servers) > 0 {
		// systemd watchdog
		go setupWatchdog(servers)

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
