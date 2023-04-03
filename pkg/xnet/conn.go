package xnet

import (
	"net"
)

type conn struct {
	net.Conn

	logger Logger
}

func WrapConn(c net.Conn, logger Logger) net.Conn {
	return &conn{
		Conn:   c,
		logger: logger,
	}
}

func (c *conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)

	c.logger.Println("read:", string(b))

	return n, err
}

func (c *conn) Write(p []byte) (n int, err error) {
	n, err = c.Conn.Write(p)

	c.logger.Println("write:", string(p))

	return n, err
}
