package main

import "flag"

type Flags struct {
	Verbose bool
	Get     string
}

var (
	flags Flags
)

func init() {
	flag.BoolVar(&flags.Verbose, "v", false, "display debugging information (verbose)")
	flag.StringVar(&flags.Get, "get", "tomorrow", "type of request")
}
