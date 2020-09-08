package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/creativeprojects/clog"
	"github.com/icholy/digest"
)

// func LoadCalendar(URL, username, password string) (*bytes.Buffer, error) {
func LoadCalendar(URL, username, password string) (*ics.Calendar, error) {
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

func GetResultFromCalendar(calendars []CalendarConfiguration) (string, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	clog.Debug(tomorrow)
	for _, calendar := range calendars {
		clog.Debugf("Loading %s...", calendar.Name)
		cal, err := LoadCalendar(calendar.URL, calendar.Username, calendar.Password)
		if err != nil {
			return "ERROR", fmt.Errorf("cannot load calendar '%s': %w", calendar.Name, err)
		}
		events := cal.Events()
		if events == nil || len(events) == 0 {
			continue
		}
		if HasMatchingEvent(tomorrow, events) {
			return calendar.Result, nil
		}
	}
	// empty value will return the default
	return "", nil
}

func HasMatchingEvent(day time.Time, events []*ics.VEvent) bool {
	for _, event := range events {
		start := event.GetProperty(ics.ComponentPropertyDtStart)
		end := event.GetProperty(ics.ComponentPropertyDtEnd)
		if start == nil || end == nil {
			continue
		}
		startTime, err := time.Parse("20060201", start.Value)
		if err != nil {
			continue
		}
		endTime, err := time.Parse("20060201", end.Value)
		if err != nil {
			continue
		}
		if day.After(startTime) && day.Before(endTime) {
			return true
		}
	}
	return false
}
