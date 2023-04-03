package client

import (
	"log"
	"net"
)

type Option func(c *config)

type config struct {
	dialer Dialer
	logger Logger
}

func newConfig(options ...Option) *config {
	c := &config{
		dialer: &net.Dialer{},
		logger: log.Default(),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithDialer(dialer Dialer) Option {
	return func(c *config) {
		c.dialer = dialer
	}
}

func WithLogger(logger Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}
