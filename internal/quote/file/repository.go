package file

import (
	"bytes"
	"context"
	"os"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/pkg/errors"
)

var (
	sep = []byte(" - ")
	eol = []byte("\n")
)

type Repository struct {
	quotes []*proto.Quote
	count  int
}

func NewRepository(filePath string) (*Repository, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read file %q", filePath)
	}

	lines := bytes.Split(b, eol)
	count := len(lines)
	quotes := make([]*proto.Quote, count)

	for i, line := range lines {
		quotes[i] = parseQuote(line)
	}

	return &Repository{
		quotes: quotes,
		count:  count,
	}, nil
}

func parseQuote(b []byte) *proto.Quote {
	parts := bytes.Split(b, sep)

	if n := len(parts); n >= 2 {
		return &proto.Quote{
			Content: string(bytes.Join(parts[:n-1], sep)),
			Author:  string(parts[n-1]),
		}
	}

	return &proto.Quote{
		Content: string(b),
	}
}

func (r *Repository) GetByID(_ context.Context, id int) (*proto.Quote, error) {
	if id < 0 {
		return nil, errors.New("invalid id")
	}

	if id > r.count-1 {
		return nil, errors.New("quote not found")
	}

	return r.quotes[id], nil
}

func (r *Repository) Count(_ context.Context) (int, error) {
	return r.count, nil
}
