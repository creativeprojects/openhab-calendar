package main

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/creativeprojects/clog"
	"gopkg.in/yaml.v3"
)

func (l *Loader) GetResultFromRules(date time.Time, rules []RuleConfiguration) (Result, error) {
	for _, rule := range rules {
		if !HasMatchingDays(date, rule.Weekdays) {
			continue
		}
		if rule.Calendar.URL == "" && rule.Calendar.File == "" {
			// no calendar to check, this is a simple weekday match
			return Result{Calendar: rule.Result}, nil
		}
		clog.Debugf("Loading %s...", rule.Name)
		cal, err := l.LoadCalendar(rule.Calendar.File, rule.Calendar.URL)
		if err != nil {
			return ResultError, fmt.Errorf("cannot load calendar '%s': %w", rule.Name, err)
		}
		events := cal.Events()
		if len(events) == 0 {
			continue
		}
		clog.Debugf("  %s calendar has %d entries", rule.Name, len(events))
		if event := FindMatchingEvent(date, events); event != nil {
			return Result{Calendar: rule.Result, Metadata: getEventMetadata(event)}, nil
		}
	}
	// empty value will return the default
	return Result{}, nil
}

// HasMatchingDays returns true if the date is in the specified weekdays.
// PLEASE NOTE the function returns TRUE when weekdays slice is empty (or nil)
func HasMatchingDays(day time.Time, weekdays []Weekday) bool {
	if len(weekdays) == 0 {
		return true
	}
	for _, weekday := range weekdays {
		if int(day.Weekday()) == int(weekday) {
			return true
		}
	}
	return false
}

func FindMatchingEvent(day time.Time, events []*ics.VEvent) *ics.VEvent {
	dateFormat := "20060102"
	loc, _ := time.LoadLocation("Local")

	for _, event := range events {
		start := event.GetProperty(ics.ComponentPropertyDtStart)
		end := event.GetProperty(ics.ComponentPropertyDtEnd)
		if start == nil || end == nil {
			continue
		}
		startTime, err := time.ParseInLocation(dateFormat, start.Value, loc)
		if err != nil {
			continue
		}
		endTime, err := time.ParseInLocation(dateFormat, end.Value, loc)
		if err != nil {
			continue
		}
		clog.Debugf("    event from %v to %v", startTime, endTime)
		if day.After(startTime) && day.Before(endTime) {
			return event
		}
	}
	clog.Debug("  --> no match")
	return nil
}

func getEventMetadata(event *ics.VEvent) []map[string]string {
	description := event.GetProperty(ics.ComponentPropertyDescription)
	if description == nil {
		return nil
	}
	value := description.Value
	if value == "" {
		return nil
	}
	value = strings.ReplaceAll(value, `\n`, "\n")
	clog.Debugf("Metadata:\n%s\n", value)

	metadata := make([]map[string]string, 0)
	err := yaml.NewDecoder(strings.NewReader(value)).Decode(&metadata)
	if err != nil {
		clog.Errorf("cannot parse metadata: %v", err)
	}
	return metadata
}
