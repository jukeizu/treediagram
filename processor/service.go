package processor

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	"github.com/jukeizu/treediagram/processor/command"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

type service struct {
	Queue    *nats.EncodedConn
	Registry registration.RegistrationClient
}

func NewService(queue *nats.EncodedConn) processing.ProcessingServer {
	return &service{Queue: queue}
}

func (s service) Request(ctx context.Context, req *processing.TreediagramRequest) (*processing.TreediagramReply, error) {
	err := s.Queue.Publish(command.RequestReceivedSubject, req)
	if err != nil {
		return nil, err
	}

	id := xid.New().String()

	return &processing.TreediagramReply{Id: id}, nil
}
