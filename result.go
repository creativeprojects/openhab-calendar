package main

type Result struct {
	Calendar string
	Metadata []map[string]string
}

var (
	ResultError = Result{
		Calendar: "ERROR",
		Metadata: nil,
	}
)
