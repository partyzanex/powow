package transport

import (
	"context"
	"net"
	"time"

	"github.com/partyzanex/powow/pkg/proto"
)

type Conn interface {
	net.Conn
}

type ChallengeProvider interface {
	CreateTaskRequest() (*proto.TaskRequest, error)
	VerifyTaskResult(id proto.ID, result *proto.TaskResult) error
	GetDeadline() time.Time
}

type QuoteService interface {
	GetRandom(ctx context.Context) (*proto.Quote, error)
}
