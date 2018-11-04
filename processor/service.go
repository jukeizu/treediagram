package processor

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	"github.com/jukeizu/treediagram/processor/command"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

type service struct {
	Queue *nats.EncodedConn
}

func NewService(logger log.Logger, queue *nats.EncodedConn, registrationClient registration.RegistrationClient) (processing.ProcessingServer, error) {
	s := &service{Queue: queue}

	received := command.NewCommandReceivedProcessor(logger, queue, registrationClient)
	err := received.Subscribe()
	if err != nil {
		return s, err
	}

	matched := command.NewCommandMatchedProcessor(logger, queue)
	err = matched.Subscribe()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s service) Request(ctx context.Context, req *processing.TreediagramRequest) (*processing.TreediagramReply, error) {
	err := s.Queue.Publish(command.RequestReceivedSubject, req)
	if err != nil {
		return nil, err
	}

	id := xid.New().String()

	return &processing.TreediagramReply{Id: id}, nil
}
