package quote

import (
	"context"
	"github.com/partyzanex/powow/pkg/proto"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (*proto.Quote, error)
	Count(ctx context.Context) (int, error)
}
