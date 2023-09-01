package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/creativeprojects/clog"
)

type CalendarResult struct {
	Calendar string              `json:"calendar"`
	Override []map[string]string `json:"override,omitempty"`
	Error    string              `json:"error,omitempty"`
}

func getCalendarHandler(config Configuration) http.HandlerFunc {
	loader := NewLoader(config)

	return func(w http.ResponseWriter, r *http.Request) {
		path := cleanupPath(r.URL.Path)
		date := r.URL.Query().Get("date")
		clog.Debugf("%s %s %s", r.Method, path, date)

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		result, err := getCalendarResult(date, config, loader)
		calendarResult := CalendarResult{
			Calendar: result.Calendar,
			Override: result.Metadata,
		}
		if err != nil {
			clog.Error(err)
			calendarResult.Error = err.Error()
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(&calendarResult)
		if err != nil {
			clog.Error(err)
		}
	}
}

func cleanupPath(path string) string {
	return strings.TrimSuffix(strings.TrimSpace(path), "/")
}
