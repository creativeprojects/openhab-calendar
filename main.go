package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"

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
	// make sure the calendars are in the right order
	sort.Slice(config.Calendars, func(i, j int) bool {
		return config.Calendars[i].Priority < config.Calendars[j].Priority
	})

	result, err := GetResultFromCalendar(config.Calendars)
	if err != nil {
		clog.Error(err)
	}
	if result == "" {
		// no match, return the default
		result = config.Default.Result
	}
	fmt.Println(result)
}
