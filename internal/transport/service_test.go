package transport

import (
	"crypto/rand"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/partyzanex/powow/internal/transport/mock"
	"github.com/partyzanex/powow/pkg/proto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestService_HandleConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var id proto.ID

	_, err := rand.Read(id[:])
	require.NoError(t, err)

	taskRequest := &proto.TaskRequest{
		ID:        id,
		Time:      1,
		Memory:    1024 * 64,
		Threads:   4,
		KeyLength: 64,
	}
	taskResult := &proto.TaskResult{
		Salt: []byte("test salt"),
		Hash: []byte("test hash"),
	}
	quote := &proto.Quote{
		Content: "test",
		Author:  "test",
	}
	deadline := time.Now().Add(time.Second)

	provider := mock.NewMockChallengeProvider(ctrl)
	provider.EXPECT().GetDeadline().Return(deadline).AnyTimes()
	provider.EXPECT().CreateTaskRequest().Return(taskRequest, nil)
	provider.EXPECT().
		VerifyTaskResult(taskRequest.ID, taskResult).
		Return(nil)

	quotes := mock.NewMockQuoteService(ctrl)
	quotes.EXPECT().GetRandom(gomock.Any()).Return(quote, nil)

	conn := mock.NewMockConn(ctrl)
	conn.EXPECT().SetDeadline(deadline).Return(nil).Times(2)
	conn.EXPECT().Write(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		return len(b), nil
	}).AnyTimes()
	conn.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		msg := &proto.Message{
			Kind: proto.KindTaskResult,
			Data: taskResult,
		}

		p, _ := json.Marshal(msg)
		n := copy(b, p)

		return n, nil
	}).AnyTimes()
	conn.EXPECT().Close().Return(nil)

	service := NewService(provider, quotes, true)
	service.HandleConn(conn)
}
