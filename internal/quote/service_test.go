package quote

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/partyzanex/powow/internal/quote/mock"
	"github.com/partyzanex/powow/pkg/proto"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &ServiceSuite{})
}

type ServiceSuite struct {
	suite.Suite

	ctrl       *gomock.Controller
	repository *mock.MockRepository
	service    *Service
}

func (s *ServiceSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.repository = mock.NewMockRepository(s.ctrl)
	s.service = NewService(s.repository)
}

func (s *ServiceSuite) AfterTest(_, _ string) {
	s.ctrl.Finish()
}

func (s *ServiceSuite) TestGetRandom() {
	const samples = 1000

	ctx := context.Background()
	want := &proto.Quote{
		Content: "test content",
		Author:  "test author",
	}

	s.repository.EXPECT().
		Count(gomock.Any()).
		Return(samples, nil).
		Times(samples)
	s.repository.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, id int) (*proto.Quote, error) {
				s.True(id >= 0)
				s.True(id < samples)

				res := *want

				return &res, nil
			},
		).
		Times(samples)

	for i := 0; i < samples; i++ {
		got, err := s.service.GetRandom(ctx)
		s.NoError(err)
		s.NotNil(got)
		s.Equal(want, got)
	}
}
