package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/creativeprojects/clog"
	"github.com/icholy/digest"
)

type httpClient struct {
	URL      string
	Username string
	Password string
	client   *http.Client
}

type Loader struct {
	httpClients []httpClient
	timeout     time.Duration
}

func NewLoader(config Configuration) *Loader {
	httpClients := make([]httpClient, len(config.Auth))
	for index, clientConfig := range config.Auth {
		clientConfig.URL = strings.ToLower(clientConfig.URL)
		clog.Debugf("add authentication for url %q", clientConfig.URL)
		httpClients[index] = httpClient{
			URL:      clientConfig.URL,
			Username: clientConfig.Username,
			Password: clientConfig.Password,
			client: &http.Client{
				Transport: &digest.Transport{
					Username: clientConfig.Username,
					Password: clientConfig.Password,
				},
			},
		}
	}
	return &Loader{
		httpClients: httpClients,
		timeout:     Timeout,
	}
}

func (l *Loader) LoadCalendar(filename, url string) (*ics.Calendar, error) {
	if filename != "" {
		return l.LoadLocalCalendar(filename)
	}
	if url != "" {
		return l.LoadRemoteCalendar(url)
	}
	return nil, errors.New("not enough parameter (need either file or url)")
}

func (l *Loader) LoadLocalCalendar(filename string) (*ics.Calendar, error) {
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

func (l *Loader) LoadRemoteCalendar(url string) (*ics.Calendar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := l.getClient(url).Do(request)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %s", response.Status)
	}
	return ics.ParseCalendar(response.Body)
}

func (l *Loader) GetResultFromCalendar(date time.Time, rules []RuleConfiguration) (string, error) {
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
		if events == nil || len(events) == 0 {
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

func (l *Loader) getClient(url string) *http.Client {
	if len(l.httpClients) == 0 {
		return http.DefaultClient
	}
	url = strings.ToLower(url)
	for _, client := range l.httpClients {
		if strings.HasPrefix(url, client.URL) {
			return client.client
		}
	}
	return http.DefaultClient
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
