package main

import (
	"errors"
	"net"
	"net/http"
	"strconv"
)

const (
	defaultHTTPPort = 8080
	defaultHost     = "localhost"
)

type options struct {
	port *int
	host *string
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

// WithHost sets the host the server listens on.
func WithHost(host string) Option {
	return func(options *options) error {
		if host == "" {
			return errors.New("host should not be empty")
		}
		options.host = &host
		return nil
	}
}

func NewServer(opts ...Option) (*http.Server, error) {
	var options options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	// Resolve the port, falling back to the default when unset.
	port := defaultHTTPPort
	if options.port != nil {
		if *options.port == 0 {
			port = randomPort()
		} else {
			port = *options.port
		}
	}
	// Resolve the host, falling back to the default when unset.
	host := defaultHost
	if options.host != nil {
		host = *options.host
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	return &http.Server{Addr: addr}, nil
}

func client() {
	_, _ = NewServer(WithPort(8080), WithHost("example.com"))
}

func randomPort() int {
	return 4 // Chosen by fair dice roll, guaranteed to be random.
}

func main() {
	c := client
	_ = c
}
