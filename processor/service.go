package processor

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

type service struct {
	Queue *nats.EncodedConn
}

func NewService(logger log.Logger, queue *nats.EncodedConn, registrationClient registration.RegistrationClient) (processing.ProcessingServer, error) {
	s := &service{Queue: queue}

	return s, nil
}

func (s service) Request(ctx context.Context, req *processing.TreediagramRequest) (*processing.TreediagramReply, error) {
	id := xid.New().String()

	return &processing.TreediagramReply{Id: id}, nil
}
