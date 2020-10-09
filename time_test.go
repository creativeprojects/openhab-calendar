package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This is to demonstrate that Parse is NOT using the local timezone, but UTC instead
func TestParseWithNoTimezone(t *testing.T) {
	format := "2006-01-02T15:04:05"
	source := "2020-09-11T00:01:01"

	date, err := time.Parse(format, source)
	require.NoError(t, err)

	assert.Equal(t, "2020-09-11 00:01:01 +0000 UTC", date.String())
}

func TestParseLocationFromDifferentTimezone(t *testing.T) {
	format := "2006-01-02T15:04:05"
	source := "2020-09-11T00:01:01"
	locations := []string{"Local", "UTC", "Europe/London", "Europe/Paris"}

	for _, location := range locations {
		t.Run(location, func(t *testing.T) {
			loc, err := time.LoadLocation(location)
			require.NoError(t, err)

			date, err := time.ParseInLocation(format, source, loc)
			require.NoError(t, err)

			assert.Equal(t, source, date.Format(format))
		})
	}
}

func TestTomorrowFromDifferentTimezone(t *testing.T) {
	format := "2006-01-02T15:04:05"
	source := "2020-09-11T00:01:01"
	expected := "2020-09-12T00:01:01"
	locations := []string{"Local", "UTC", "Europe/London", "Europe/Paris"}

	for _, location := range locations {
		t.Run(location, func(t *testing.T) {
			loc, err := time.LoadLocation(location)
			require.NoError(t, err)

			date, err := time.ParseInLocation(format, source, loc)
			require.NoError(t, err)

			tomorrow := getTomorrow(date)

			assert.Equal(t, expected, tomorrow.Format(format))
		})
	}
}
