package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/creativeprojects/clog"
)

// HTTPServer encapsulates a *http.Server
type HTTPServer struct {
	config ServerConfiguration
	name   string
	listen string
	tls    bool
	server *http.Server
}

// NewHTTPServer creates a new HTTPServer object after checking that the configuration is valid
func NewHTTPServer(name string, config ServerConfiguration) (*HTTPServer, error) {
	address, err := url.Parse(config.Listen)
	if err != nil {
		return nil, fmt.Errorf("invalid listen address in configuration: %v", err)
	}
	listen := getListenAddress(address)
	if listen == "" {
		return nil, fmt.Errorf("invalid listen address in configuration, the accepted format is \"scheme://host:port\" but found %q", config.Listen)
	}
	tls := address.Scheme == "https"
	if tls {
		if _, err := os.Stat(config.Certificate); os.IsNotExist(err) {
			return nil, fmt.Errorf("TLS certificate not found: %q", config.Certificate)
		}
		if _, err := os.Stat(config.PrivateKey); os.IsNotExist(err) {
			return nil, fmt.Errorf("TLS private key not found: %q", config.PrivateKey)
		}
	}
	return &HTTPServer{
		config: config,
		name:   name,
		listen: listen,
		tls:    tls,
	}, nil
}

// Start HTTP(s) server
func (s *HTTPServer) Start() {
	clog.Infof("%v: listening on %q", s.name, s.listen)
	s.server = &http.Server{
		Addr:    s.listen,
		Handler: getServeMux(),
	}
	if s.tls {
		if err := s.server.ListenAndServeTLS(s.config.Certificate, s.config.PrivateKey); err != http.ErrServerClosed {
			clog.Error(err.Error())
		}
	} else {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			clog.Error(err.Error())
		}
	}
	clog.Warningf("%v: stopped listening", s.name)
}

func getServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return mux
}

// Shutdown HTTP(s) server gracefully
func (s *HTTPServer) Shutdown(wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	if s.server != nil {
		clog.Debugf("%v: shutting down server...", s.name)
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

func getListenAddress(listen *url.URL) string {
	host := listen.Host
	port := listen.Port()
	if port == "" {
		if listen.Scheme == "http" {
			return host + ":80"
		}
		if listen.Scheme == "https" {
			return host + ":443"
		}
	}
	return host
}
