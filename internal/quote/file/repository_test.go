package file

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
)

func TestRepositorySuite(t *testing.T) {
	repo, err := NewRepository("../../../assets/quotes.txt")
	require.NoError(t, err)
	require.NotNil(t, repo)

	suite.Run(t, &RepositorySuite{
		repository: repo,
	})
}

type RepositorySuite struct {
	suite.Suite

	repository *Repository
}

func (s *RepositorySuite) TestCount() {
	count, err := s.repository.Count(context.Background())
	s.NoError(err)
	s.Equal(109, count)
}

func (s *RepositorySuite) TestGetByID() {
	id := rand.Intn(109)
	ctx := context.Background()

	got, err := s.repository.GetByID(ctx, id)
	s.NoError(err)
	s.NotNil(got)

	got, err = s.repository.GetByID(ctx, 1234)
	s.Error(err)
	s.Nil(got)
}
