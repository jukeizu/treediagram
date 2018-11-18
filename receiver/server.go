package receiver

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/receiving"
	nats "github.com/nats-io/go-nats"
)

const (
	MessageReceivedSubject = "message.received"
)

type service struct {
	queue *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn) receiving.ReceivingServer {
	return &service{queue: queue}
}

func (s service) Send(ctx context.Context, req *receiving.Message) (*receiving.Reply, error) {
	message := NewMessage(req)

	err := s.queue.Publish(MessageReceivedSubject, message)
	if err != nil {
		return nil, err
	}

	return &receiving.Reply{Id: message.Id}, nil
}
