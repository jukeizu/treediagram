package processor

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
)

const (
	RequestReceivedSubject = "request.received"
)

type service struct {
	storage Storage
	queue   *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn, storage Storage) (processing.ProcessingServer, error) {
	return &service{queue: queue, storage: storage}, nil
}

func (s service) SendRequest(ctx context.Context, req *processing.Request) (*processing.SendReply, error) {
	err := s.queue.Publish(RequestReceivedSubject, req)

	return &processing.SendReply{Id: req.Id}, err
}

func (s service) GetMessage(ctx context.Context, req *processing.MessageRequest) (*processing.MessageReply, error) {
	return s.storage.MessageReply(req.Id)
}
