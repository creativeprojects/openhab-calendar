package main

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/creativeprojects/clog"
)

func (l *Loader) GetResultFromRules(date time.Time, rules []RuleConfiguration) (string, error) {
	for _, rule := range rules {
		if !HasMatchingDays(date, rule.Weekdays) {
			continue
		}
		if rule.Calendar.URL == "" && rule.Calendar.File == "" {
			// no calendar to check, this is a simple weekday match
			return rule.Result, nil
		}
		clog.Debugf("Loading %s...", rule.Name)
		cal, err := l.LoadCalendar(rule.Calendar.File, rule.Calendar.URL)
		if err != nil {
			return "ERROR", fmt.Errorf("cannot load calendar '%s': %w", rule.Name, err)
		}
		events := cal.Events()
		if len(events) == 0 {
			continue
		}
		clog.Debugf("  %s calendar has %d entries", rule.Name, len(events))
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

func HasMatchingEvent(day time.Time, events []*ics.VEvent) bool {
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
			return true
		}
	}
	clog.Debug("  --> no match")
	return false
}
