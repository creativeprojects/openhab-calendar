package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidListenAddress(t *testing.T) {
	testData := []struct {
		address string
		listen  string
	}{
		{"http://:8080", ":8080"},
		{"http://0.0.0.0:8080", "0.0.0.0:8080"},
		{"http://", ":80"},
		{"https://", ":443"},
	}
	for _, testItem := range testData {
		t.Run(testItem.address, func(t *testing.T) {
			address, err := url.Parse(testItem.address)
			assert.NoError(t, err)
			assert.Equal(t, testItem.listen, getListenAddress(address))
		})
	}
}

func TestInvalidListenAddress(t *testing.T) {
	testData := []struct {
		address string
	}{
		{""},
		{"0.0.0.0"},
		{":8080"},
		{"tcp://"},
	}
	for _, testItem := range testData {
		t.Run(testItem.address, func(t *testing.T) {
			address, err := url.Parse(testItem.address)
			if err == nil {
				assert.Equal(t, "", getListenAddress(address))
			}
		})
	}
}
