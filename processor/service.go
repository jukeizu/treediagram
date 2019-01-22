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
	repository Repository
	queue      *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn, repository Repository) (processing.ProcessingServer, error) {
	return &service{queue: queue, repository: repository}, nil
}

func (s service) SendMessageRequest(ctx context.Context, req *processing.MessageRequest) (*processing.SendMessageRequestReply, error) {
	err := s.queue.Publish(RequestReceivedSubject, req)

	return &processing.SendMessageRequestReply{Id: req.Id}, err
}

func (s service) GetMessageReply(ctx context.Context, req *processing.MessageReplyRequest) (*processing.MessageReply, error) {
	return s.repository.MessageReply(req.Id)
}
