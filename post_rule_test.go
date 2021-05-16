package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchPostRules(t *testing.T) {
	testData := []struct {
		value   string
		matcher PostRuleMatcher
		match   bool
	}{
		{"ONE", PostRuleMatcher{Is: "ONE"}, true},
		{"ONE", PostRuleMatcher{Is: "TWO"}, false},
		{"ONE", PostRuleMatcher{Not: "ONE"}, false},
		{"ONE", PostRuleMatcher{Not: "TWO"}, true},
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testItem.match, MatchPostRule(testItem.value, testItem.matcher))
		})
	}
}

func TestGetResultFromPostRules(t *testing.T) {
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
	postRules := []PostRuleConfiguration{
		{
			When: &PostRuleMatcher{
				Is: "Weekend",
			},
			Previous: &PostRuleMatcher{
				Not: "Weekend",
			},
			Next: &PostRuleMatcher{
				Is: "Weekend",
			},
			Result: "FirstOff",
		},
	}
	testData := []struct {
		day    string
		result string
	}{
		{"20210514", "Weekday"},
		{"20210515", "FirstOff"},
		{"20210516", "Weekend"},
		{"20210517", "Weekday"},
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
			result, err = loader.PostRules(result, day, rules, postRules)
			require.NoError(t, err)
			assert.Equal(t, testItem.result, result)
		})
	}
}
