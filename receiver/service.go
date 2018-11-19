package receiver

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/receiving"
	nats "github.com/nats-io/go-nats"
)

const (
	RequestReceivedSubject = "request.received"
)

type service struct {
	queue *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn) receiving.ReceivingServer {
	return &service{queue: queue}
}

func (s service) Send(ctx context.Context, req *receiving.Request) (*receiving.Reply, error) {
	err := s.queue.Publish(RequestReceivedSubject, req)

	return &receiving.Reply{Id: req.Id}, err
}
