package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/creativeprojects/clog"
	"github.com/icholy/digest"
)

func LoadCalendar(filename, URL, username, password string) (*ics.Calendar, error) {
	if filename != "" {
		return LoadLocalCalendar(filename)
	}
	if URL != "" {
		return LoadRemoteCalendar(URL, username, password)
	}
	return nil, errors.New("not enough parameter (need either file or url)")
}

func LoadLocalCalendar(filename string) (*ics.Calendar, error) {
	// path relative to the binary
	if !path.IsAbs(filename) {
		me, err := os.Executable()
		if err == nil {
			dir := path.Dir(me)
			filename = path.Join(dir, filename)
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ics.ParseCalendar(file)
}

func LoadRemoteCalendar(URL, username, password string) (*ics.Calendar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &digest.Transport{
			Username: username,
			Password: password,
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %s", response.Status)
	}
	return ics.ParseCalendar(response.Body)
}

func GetResultFromCalendar(date time.Time, rules []RuleConfiguration) (string, error) {
	for _, rule := range rules {
		if !HasMatchingDays(date, rule.Weekdays) {
			continue
		}
		if rule.Calendar.URL == "" && rule.Calendar.File == "" {
			// no calendar to check, this is a simple weekday match
			return rule.Result, nil
		}
		clog.Debugf("Loading %s...", rule.Name)
		cal, err := LoadCalendar(rule.Calendar.File, rule.Calendar.URL, rule.Calendar.Username, rule.Calendar.Password)
		if err != nil {
			return "ERROR", fmt.Errorf("cannot load calendar '%s': %w", rule.Name, err)
		}
		events := cal.Events()
		if events == nil || len(events) == 0 {
			continue
		}
		if HasMatchingEvent(date, events) {
			return rule.Result, nil
		}
	}
	// empty value will return the default
	return "", nil
}

// HasMatchingDays returns true if the date is in the specified weekdays.
// PLEASE NOTE the function returns TRUE when weekdays slice is empty (or nil)
func HasMatchingDays(day time.Time, weekdays []Weekday) bool {
	if weekdays == nil || len(weekdays) == 0 {
		return true
	}
	for _, weekday := range weekdays {
		if int(day.Weekday()) == int(weekday) {
			return true
		}
	}
	return false
}

func HasMatchingEvent(day time.Time, events []*ics.VEvent) bool {
	dateFormat := "20060102"
	for _, event := range events {
		start := event.GetProperty(ics.ComponentPropertyDtStart)
		end := event.GetProperty(ics.ComponentPropertyDtEnd)
		if start == nil || end == nil {
			continue
		}
		startTime, err := time.Parse(dateFormat, start.Value)
		if err != nil {
			continue
		}
		endTime, err := time.Parse(dateFormat, end.Value)
		if err != nil {
			continue
		}
		if day.After(startTime) && day.Before(endTime) {
			return true
		}
	}
	return false
}
