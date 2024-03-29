package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchingDays(t *testing.T) {
	dateFormat := "20060102"
	testData := []struct {
		day      string
		weekdays []Weekday
		expected bool
	}{
		{"20200907", []Weekday{0}, false},
		{"20200908", []Weekday{2}, true},
		{"20200909", []Weekday{1, 2, 3, 4, 5}, true},
		{"20200910", []Weekday{0, 6}, false},
		{"20200911", []Weekday{1, 2, 3, 4, 5}, true},
		{"20200912", []Weekday{0, 6}, true},
		{"20200913", []Weekday{0, 6}, true},
		{"20200914", []Weekday{}, true},
		{"20200915", nil, true},
	}

	loc, err := time.LoadLocation("Local")
	require.NoError(t, err)

	for _, testItem := range testData {
		t.Run(testItem.day, func(t *testing.T) {
			day, err := time.ParseInLocation(dateFormat, testItem.day, loc)
			require.NoError(t, err)
			assert.Equal(t, testItem.expected, HasMatchingDays(day, testItem.weekdays))
		})
	}
}

func TestGetResultFromRules(t *testing.T) {
	dateFormat := "20060102"
	rules := []RuleConfiguration{
		{
			Result:   "Weekday",
			Weekdays: []Weekday{1, 2, 3, 4, 5},
		},
		{
			Result:   "Weekend",
			Weekdays: []Weekday{0, 6},
		},
	}
	testData := []struct {
		day    string
		result string
	}{
		{"20200907", "Weekday"},
		{"20200908", "Weekday"},
		{"20200909", "Weekday"},
		{"20200910", "Weekday"},
		{"20200911", "Weekday"},
		{"20200912", "Weekend"},
		{"20200913", "Weekend"},
		{"20200914", "Weekday"},
		{"20200915", "Weekday"},
	}

	loc, err := time.LoadLocation("Local")
	require.NoError(t, err)

	loader := NewLoader(Configuration{})

	for _, testItem := range testData {
		t.Run(testItem.day, func(t *testing.T) {
			day, err := time.ParseInLocation(dateFormat, testItem.day, loc)
			require.NoError(t, err)

			result, err := loader.GetResultFromRules(day, rules)
			require.NoError(t, err)
			assert.Equal(t, testItem.result, result.Calendar)
		})
	}
}
