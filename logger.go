package main

import (
	"log"

	"github.com/creativeprojects/clog"
)

func setupLogger(verbose bool) func() {
	if verbose {
		clog.SetDefaultLogger(clog.NewConsoleLogger())
		return func() {}
	}
	// otherwise use a file log
	fileHandler, err := clog.NewFileHandler(LogFile, "", log.LstdFlags)
	if err != nil {
		// just forget it
		return func() {}
	}
	clog.SetDefaultLogger(clog.NewLogger(clog.NewLevelFilter(clog.LevelInfo, fileHandler)))
	return fileHandler.Close
}
