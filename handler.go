package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/creativeprojects/clog"
)

type CalendarResult struct {
	Calendar string `json:"calendar"`
	Error    string `json:"error,omitempty"`
}

func getCalendarHandler(config Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := cleanupPath(r.URL.Path)
		date := r.URL.Query().Get("date")
		clog.Debugf("%s %s %s", r.Method, path, date)

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var (
			result CalendarResult
			err    error
		)
		result.Calendar, err = getCalendarResult(date, config)
		if err != nil {
			clog.Error(err)
			result.Error = err.Error()
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(&result)
	}
}

func cleanupPath(path string) string {
	return strings.TrimSuffix(strings.TrimSpace(path), "/")
}
