package client

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"net"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	network  = "tcp"
	saltSize = 24
)

type Dialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

type Logger interface {
	Println(v ...any)
}

type Client struct {
	dialer Dialer
	logger Logger

	serverAddress string
}

func NewClient(serverAddress string, options ...Option) *Client {
	cfg := newConfig(options...)

	return &Client{
		dialer:        cfg.dialer,
		logger:        cfg.logger,
		serverAddress: serverAddress,
	}
}

func (c *Client) GetRandomWisdom(ctx context.Context) (*proto.Quote, error) {
	netConn, err := c.dialer.DialContext(ctx, network, c.serverAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open tcp connection")
	}

	defer func() {
		if closeErr := netConn.Close(); closeErr != nil {
			c.logger.Println("cannot close connection:", closeErr)
		}
	}()

	decoder := json.NewDecoder(netConn)
	msg := new(proto.Message)

	err = decoder.Decode(msg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode message")
	}

	if msg.Kind != proto.KindTaskRequest {
		return nil, errors.Errorf("unexpected message kind: %s", msg.Kind)
	}

	hash, salt, err := c.executeTask(msg.Data.(*proto.TaskRequest))
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute task")
	}

	encoder := json.NewEncoder(netConn)

	err = encoder.Encode(&proto.Message{
		Kind: proto.KindTaskResult,
		Data: &proto.TaskResult{
			Salt: salt,
			Hash: hash,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot encode task result")
	}

	msg = new(proto.Message)

	err = decoder.Decode(msg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode message")
	}

	switch msg.Kind {
	case proto.KindTaskError:
		return nil, msg.GetTaskError()
	case proto.KindQuote:
		return msg.GetQuote(), nil
	default:
		return nil, errors.Errorf("unexpected message kind: %s", msg.Kind)
	}
}

func (c *Client) executeTask(taskRequest *proto.TaskRequest) (hash, salt []byte, err error) {
	salt = make([]byte, saltSize)

	_, err = rand.Read(salt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot create salt")
	}

	hash = argon2.IDKey(
		taskRequest.ID[:],
		salt,
		taskRequest.Time,
		taskRequest.Memory,
		taskRequest.Threads,
		taskRequest.KeyLength,
	)

	return hash, salt, nil
}
