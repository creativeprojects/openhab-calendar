package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %s", response.Status)
	}
	return ics.ParseCalendar(response.Body)
}

func (l *Loader) SaveRemoteCalendar(url, to string) error {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	response, err := l.getClient(url).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", response.Status)
	}
	// read all remote content
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// open the destination file
	file, err := os.Create(to)
	if err != nil {
		return err
	}
	defer file.Close()

	// copy the calendar to the file
	_, err = file.Write(content)
	return err
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
