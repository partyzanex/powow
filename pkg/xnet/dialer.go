package xnet

import (
	"context"
	"net"
)

type Logger interface {
	Println(v ...any)
}

type Dialer struct {
	*net.Dialer

	logger Logger
	debug  bool
}

func NewDialer(netDialer *net.Dialer, logger Logger, debug bool) *Dialer {
	return &Dialer{
		Dialer: netDialer,
		logger: logger,
		debug:  debug,
	}
}

func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	c, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	if d.debug {
		return &conn{
			Conn:   c,
			logger: d.logger,
		}, nil
	}

	return c, nil
}
