package main

import (
	"encoding/json"
	"io"
	"os"
)

type Configuration struct {
	Calendars []CalendarConfiguration `json:"calendars"`
	Default   DefaultConfiguration    `json:"default"`
}

type CalendarConfiguration struct {
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Result   string `json:"result"`
}

type DefaultConfiguration struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

func LoadFileConfiguration(filename string) (Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	return LoadConfiguration(file)
}

func LoadConfiguration(reader io.Reader) (Configuration, error) {
	decoder := json.NewDecoder(reader)
	var config Configuration
	err := decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
