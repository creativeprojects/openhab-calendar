package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var days = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func (wd *Weekday) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, day := range days {
		if s == strings.ToLower(day) || s[:3] == strings.ToLower(day)[:3] {
			*wd = Weekday(i)
			return nil
		}
	}

	return fmt.Errorf("invalid weekday: %s", s)
}
