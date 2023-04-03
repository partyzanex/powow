package transport

import (
	"context"
	"encoding/json"
	"github.com/partyzanex/powow/pkg/xnet"
	"log"
	"net"
	"sync"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/pkg/errors"
)

type Service struct {
	provider ChallengeProvider
	quotes   QuoteService

	wg    sync.WaitGroup
	mu    sync.Mutex
	done  chan struct{}
	debug bool
}

func NewService(provider ChallengeProvider, quotes QuoteService, debug bool) *Service {
	return &Service{
		provider: provider,
		quotes:   quotes,
		wg:       sync.WaitGroup{},
		mu:       sync.Mutex{},
		done:     make(chan struct{}),
		debug:    debug,
	}
}

func (s *Service) HandleConn(conn Conn) {
	if s.debug {
		s.handleConn(xnet.WrapConn(conn, log.Default()))
	} else {
		s.handleConn(conn)
	}
}

func (s *Service) handleConn(conn Conn) {
	select {
	case <-s.done:
		return
	default:
	}

	s.take()
	defer s.release()

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("cannot close connection:", err)
		}
	}()

	if err := s.challenge(conn); err != nil {
		log.Println("challenge failed:", err)
		return
	}

	if err := s.sendData(conn); err != nil {
		log.Println("handling failed:", err)
		return
	}
}

func (s *Service) sendData(conn Conn) error {
	deadline := s.provider.GetDeadline()

	if err := conn.SetDeadline(deadline); err != nil {
		return errors.Wrap(err, "cannot set deadline")
	}

	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	quote, err := s.quotes.GetRandom(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot get random quote")
	}

	err = json.NewEncoder(conn).Encode(proto.Message{
		Kind: proto.KindQuote,
		Data: quote,
	})
	if err != nil {
		return errors.Wrap(err, "cannot encode message")
	}

	return nil
}

func (s *Service) challenge(conn net.Conn) error {
Start:
	if err := conn.SetDeadline(s.provider.GetDeadline()); err != nil {
		return errors.Wrap(err, "cannot set deadline")
	}

	taskRequest, err := s.provider.CreateTaskRequest()
	if err != nil {
		return errors.Wrap(err, "cannot create task request")
	}

	encoder := json.NewEncoder(conn)

	err = encoder.Encode(&proto.Message{
		Kind: proto.KindTaskRequest,
		Data: taskRequest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot encode task request")
	}

	msg := new(proto.Message)
	decoder := json.NewDecoder(conn)

	if err = decoder.Decode(msg); err != nil {
		if encodeErr := encoder.Encode(proto.NewTaskError(err)); encodeErr != nil {
			err = errors.Wrapf(encodeErr, "cannot encode task error: %s", err)
		}

		return err
	}

	taskResult := msg.GetTaskResult()
	if taskResult == nil {
		return errors.New("invalid task result")
	}

	if err = s.provider.VerifyTaskResult(taskRequest.ID, taskResult); err != nil {
		if encodeErr := encoder.Encode(proto.NewTaskError(err)); encodeErr != nil {
			return errors.Wrapf(encodeErr, "cannot write task error: %s", err)
		}

		goto Start
	}

	return nil
}

func (s *Service) take() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wg.Add(1)
}

func (s *Service) release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wg.Done()
}

func (s *Service) Close() error {
	close(s.done)
	s.wg.Wait()

	return nil
}
