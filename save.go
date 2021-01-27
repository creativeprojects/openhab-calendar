package main

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/creativeprojects/clog"
)

func saveCalendars(config Configuration) {
	execPath, _ := os.Executable()
	savePath := filepath.Join(filepath.Dir(execPath), "files")
	err := os.MkdirAll(savePath, 0777)
	if err != nil {
		clog.Error(err)
		return
	}

	loader := NewLoader(config)

	for _, rule := range config.Rules {
		if rule.Calendar.URL == "" {
			continue
		}
		clog.Debugf("Loading %q", rule.Calendar.URL)
		URL, err := url.Parse(rule.Calendar.URL)
		if err != nil {
			clog.Error(err)
			continue
		}
		base := path.Base(URL.Path)
		if base == "" {
			continue
		}
		saveFile := filepath.Join(savePath, base) + ".ics"
		clog.Debugf("Saving %q", saveFile)
		err = loader.SaveRemoteCalendar(rule.Calendar.URL, saveFile)
		if err != nil {
			clog.Error(err)
		}
	}
}
