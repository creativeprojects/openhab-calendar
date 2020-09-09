package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalWeekday(t *testing.T) {
	testData := []struct {
		day      string
		expected Weekday
	}{
		{`"sUn"`, 0}, {`"sUndAy"`, 0},
		{`"mOn"`, 1}, {`"mOndAy"`, 1},
		{`"tUe"`, 2}, {`"tUesdAy"`, 2},
		{`"wEd"`, 3}, {`"wEdnesdAy"`, 3},
		{`"tHu"`, 4}, {`"tHursdAy"`, 4},
		{`"fRi"`, 5}, {`"fRidAy"`, 5},
		{`"sAt"`, 6}, {`"sAturdAy"`, 6},
	}
	for _, testItem := range testData {
		t.Run(testItem.day, func(t *testing.T) {
			var weekday Weekday
			err := weekday.UnmarshalJSON([]byte(testItem.day))
			assert.NoError(t, err)
			assert.Equal(t, testItem.expected, weekday)
		})
	}
}

func TestErrorUnmarshalWeekday(t *testing.T) {
	testData := []string{"", "1", "su", "son", "sund", "sunto", "sunday2"}
	for _, testItem := range testData {
		t.Run(testItem, func(t *testing.T) {
			var weekday Weekday
			err := weekday.UnmarshalJSON([]byte(testItem))
			assert.Error(t, err)
		})
	}
}
