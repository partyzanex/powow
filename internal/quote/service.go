package quote

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/pkg/errors"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) GetRandom(ctx context.Context) (*proto.Quote, error) {
	count, err := s.repository.Count(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get count of quotes")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(count)))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get random index")
	}

	quote, err := s.repository.GetByID(ctx, int(n.Int64()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get quote by id")
	}

	return quote, nil
}
