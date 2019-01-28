package processor

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	nats "github.com/nats-io/go-nats"
)

type service struct {
	repository Repository
	queue      *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn, repository Repository) (processingpb.ProcessingServer, error) {
	return &service{queue: queue, repository: repository}, nil
}

func (s service) SendMessageRequest(ctx context.Context, req *processingpb.MessageRequest) (*processingpb.SendMessageRequestReply, error) {
	err := s.queue.Publish(MessageRequestReceivedSubject, req)

	return &processingpb.SendMessageRequestReply{Id: req.Id}, err
}

func (s service) GetMessageReply(ctx context.Context, req *processingpb.MessageReplyRequest) (*processingpb.MessageReply, error) {
	return s.repository.MessageReply(req.Id)
}
