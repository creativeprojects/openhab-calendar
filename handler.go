package main

import (
	"net/http"
	"strings"

	"github.com/creativeprojects/clog"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := cleanupPath(r.URL.Path)
	clog.Debugf("%s %s", r.Method, path)

	clog.Debugf("%s %s finished", r.Method, path)
}

func cleanupPath(path string) string {
	return strings.TrimSuffix(strings.TrimSpace(path), "/")
}
