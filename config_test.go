package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleWeekConfiguration(t *testing.T) {
	json := `
{
    "rules": [
        {
            "priority": 10,
            "name": "Weekday",
            "weekdays": [
                "Mon",
                "Tue",
                "Wed",
                "Thu",
                "Fri"
            ],
            "result": "Weekday"
        },
        {
            "priority": 20,
            "name": "Weekend",
            "weekdays": [
                "Sat",
                "Sun"
            ],
            "result": "Weekend"
        }
    ],
    "default": {
        "name": "Unknown",
        "result": "ERROR"
    }
}`
	config, err := LoadConfiguration(bytes.NewBufferString(json))
	assert.NoError(t, err)

	assert.Len(t, config.Rules, 2)
	assert.ElementsMatch(t, config.Rules[0].Weekdays, []Weekday{1, 2, 3, 4, 5})
	assert.ElementsMatch(t, config.Rules[1].Weekdays, []Weekday{0, 6})
}
