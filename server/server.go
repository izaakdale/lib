package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	configOptions struct {
		host              *string
		port              *string
		readHeaderTimeout *time.Duration
		readTimeout       *time.Duration
		timeoutHandler    *timeoutHandlerOptions
	}
	option                func(opt *configOptions) error
	timeoutHandlerOptions struct {
		timeout time.Duration
		msg     string
	}
)

type server struct {
	Name string
	*http.Server
}

// Run prints the name of the server and the address and continues to ListenAndServe.
func (s *server) Run() error {
	log.Printf("%s running on %s\n", s.Name, s.Addr)
	return s.ListenAndServe()
}

// NewServer returns a http.Server with the specified options.
// If WithPort and WithHost are not used server address defaults to ":http"
func New(name string, handler http.Handler, optFuncs ...option) (*server, error) {
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

	if options.timeoutHandler != nil {
		handler = http.TimeoutHandler(handler, options.timeoutHandler.timeout, options.timeoutHandler.msg)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: hTimeout,
		ReadTimeout:       timeout,
	}

	return &server{name, srv}, nil
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
	if p == "" {
		p = "80"
	}
	return func(opt *configOptions) error {
		opt.port = &p
		return nil
	}
}

// WithTimeouts adds header and request read timeouts to the server.
// Defaults to 1s for both.
func WithTimeouts(header, total time.Duration) option {
	return func(opt *configOptions) error {
		opt.readHeaderTimeout = &header
		opt.readTimeout = &total
		return nil
	}
}

// WithTimeoutHanlder cancels the context after duration dt
// and prints msg to the caller with a 503 code.
func WithTimeoutHandler(dt time.Duration, msg string) option {
	return func(opt *configOptions) error {
		opt.timeoutHandler = &timeoutHandlerOptions{
			timeout: dt,
			msg:     msg,
		}
		return nil
	}
}
