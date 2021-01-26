package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/creativeprojects/clog"
	"github.com/icholy/digest"
)

type Client struct {
	URL      string
	Username string
	Password string
	client   *http.Client
}

type HTTPClient struct {
	clients []Client
	other   *http.Client
}

func NewHTTPClient(config Configuration) *HTTPClient {
	clients := make([]Client, len(config.Auth))
	for index, clientConfig := range config.Auth {
		clientConfig.URL = strings.ToLower(clientConfig.URL)
		clog.Debugf("add authentication for url %q", clientConfig.URL)
		clients[index] = Client{
			URL:      clientConfig.URL,
			Username: clientConfig.Username,
			Password: clientConfig.Password,
			client: &http.Client{
				Transport: &digest.Transport{
					Username: clientConfig.Username,
					Password: clientConfig.Password,
					Transport: &http.Transport{
						DialContext: (&net.Dialer{
							Timeout:   1 * time.Second,
							KeepAlive: 1 * time.Second,
							DualStack: true,
						}).DialContext,
						ForceAttemptHTTP2:     true,
						MaxIdleConns:          1,
						IdleConnTimeout:       1 * time.Second,
						TLSHandshakeTimeout:   10 * time.Second,
						ExpectContinueTimeout: 1 * time.Second,
					},
				},
			},
		}
	}
	return &HTTPClient{
		clients: clients,
		other:   &http.Client{},
	}
}

// Get a url contents.
// Please note it is the responsability of the caller to close the io.ReadCloser
func (c *HTTPClient) Get(url string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.getClient(url).Do(request)
	if err != nil {
		response.Body.Close()
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("server returned %s", response.Status)
	}
	return response.Body, nil
}

func (c *HTTPClient) getClient(url string) *http.Client {
	if len(c.clients) == 0 {
		return c.other
	}
	url = strings.ToLower(url)
	for _, client := range c.clients {
		if strings.HasPrefix(url, client.URL) {
			return client.client
		}
	}
	return c.other
}
