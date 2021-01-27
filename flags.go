package main

import "flag"

type Flags struct {
	Verbose    bool
	Get        string
	Date       string
	Daemon     bool
	ConfigFile string
	Save       bool
}

var (
	flags Flags
)

func init() {
	flag.BoolVar(&flags.Verbose, "v", false, "display debugging information (verbose)")
	flag.StringVar(&flags.Get, "get", "", "deprecated: use date instead")
	flag.StringVar(&flags.Date, "date", "tomorrow", "type of request")
	flag.BoolVar(&flags.Daemon, "d", false, "demonize and answer http requests")
	flag.StringVar(&flags.ConfigFile, "c", ConfigFile, "configuration file")
	flag.BoolVar(&flags.Save, "s", false, "save remote calendars into ./files/*.ics")
}
