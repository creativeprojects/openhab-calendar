package main

import (
	"encoding/json"
	"io"
	"os"
)

type Configuration struct {
	Rules   []RuleConfiguration            `json:"rules"`
	Default DefaultConfiguration           `json:"default"`
	Auth    []AuthConfiguration            `json:"authentication"`
	Servers map[string]ServerConfiguration `json:"servers"`
}

type RuleConfiguration struct {
	Priority int                   `json:"priority"`
	Name     string                `json:"name"`
	Weekdays []Weekday             `json:"weekdays"`
	Calendar CalendarConfiguration `json:"calendar"`
	Result   string                `json:"result"`
}

type CalendarConfiguration struct {
	File string `json:"file"`
	URL  string `json:"url"`
}

type DefaultConfiguration struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

type AuthConfiguration struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfiguration struct {
	Listen      string `json:"listen"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
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
