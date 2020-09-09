package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/creativeprojects/clog"
)

func main() {
	flag.Parse()
	closeLogger := setupLogger(flags.Verbose)
	defer closeLogger()

	// set the configuration file relative to the the binary
	configFile := ConfigFile
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

	date, err := parseGetFlag(flags.Get)
	if err != nil {
		clog.Error(fmt.Errorf("cannot parse -get option: %w", err))
	}
	clog.Debug(date)
	result, err := GetResultFromCalendar(date, config.Rules)
	if err != nil {
		clog.Error(err)
	}
	if result == "" {
		// no match, return the default
		clog.Error("no match")
		result = config.Default.Result
	}
	fmt.Println(result)
}

func parseGetFlag(get string) (time.Time, error) {
	get = strings.TrimSpace(get)
	if get == "" || strings.ToLower(get) == "tomorrow" {
		// default to tomorrow
		return time.Now().AddDate(0, 0, 1), nil
	}
	return time.Parse(time.RFC3339, get)
}
