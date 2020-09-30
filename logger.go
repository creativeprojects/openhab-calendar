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
		// console - no verbose
		clog.SetDefaultLogger(clog.NewLogger(clog.NewLevelFilter(clog.LevelInfo, clog.NewConsoleHandler("", log.LstdFlags))))
		return func() {}
	}
	clog.SetDefaultLogger(clog.NewLogger(clog.NewLevelFilter(clog.LevelInfo, fileHandler)))
	return fileHandler.Close
}
