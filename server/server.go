package server

import (
	"fmt"
	"net/http"
	"time"
)

type (
	configOptions struct {
		host              *string
		port              *string
		readHeaderTimeout *time.Duration
		readTimeout       *time.Duration
	}
	option func(opt *configOptions) error
)

// NewServer returns a http.Server with the specified options.
// If WithPort and WithHost are not used server address defaults to ":http"
func New(handler http.Handler, optFuncs ...option) (*http.Server, error) {
	var options configOptions
	for _, optFunc := range optFuncs {
		err := optFunc(&options)
		if err != nil {
			return nil, err
		}
	}

	var port string
	if options.port == nil {
		port = "80"
	} else {
		port = *options.port
	}

	var host string
	if options.host == nil {
		host = ""
	} else {
		host = *options.host
	}

	var hTimeout time.Duration
	if options.readHeaderTimeout == nil {
		hTimeout = time.Second
	} else {
		hTimeout = *options.readHeaderTimeout
	}

	var timeout time.Duration
	if options.readTimeout == nil {
		timeout = time.Second
	} else {
		timeout = *options.readTimeout
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: hTimeout,
		ReadTimeout:       timeout,
	}

	return srv, nil
}

// WithHost adds the hostname to the configOptions to use with NewServer.
// Defaults local connectivity
func WithHost(h string) option {
	return func(opt *configOptions) error {
		opt.host = &h
		return nil
	}
}

// WithPort adds the specified to port to the configOptions to use with NewServer.
// Defaults to 80 if no port is specified
func WithPort(p string) option {
	return func(opt *configOptions) error {
		opt.port = &p
		return nil
	}
}

func WithTimeouts(header, total time.Duration) option {
	return func(opt *configOptions) error {
		opt.readHeaderTimeout = &header
		opt.readTimeout = &total
		return nil
	}
}
