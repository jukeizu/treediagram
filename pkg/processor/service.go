package processor

import (
	"context"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	nats "github.com/nats-io/nats.go"
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

func (s service) SendReaction(ctx context.Context, req *processingpb.Reaction) (*processingpb.SendReactionReply, error) {
	err := s.queue.Publish(ReactionReceivedSubject, req)

	return &processingpb.SendReactionReply{}, err
}

func (s service) SendInteraction(ctx context.Context, req *processingpb.Interaction) (*processingpb.SendInteractionReply, error) {
	err := s.queue.Publish(InteractionReceivedSubject, req)

	return &processingpb.SendInteractionReply{}, err
}

func (s service) SendProcessingEvent(ctx context.Context, req *processingpb.ProcessingEvent) (*processingpb.SendProcessingEventReply, error) {
	err := s.queue.Publish(EventReceivedSubject, req)

	return &processingpb.SendProcessingEventReply{}, err
}

func (s service) ProcessingRequestIntentStatistics(ctx context.Context, req *processingpb.ProcessingRequestIntentStatisticsRequest) (*processingpb.ProcessingRequestIntentStatisticsReply, error) {
	if req.CreatedLessThanOrEqualTo == 0 {
		req.CreatedLessThanOrEqualTo = time.Now().UTC().Unix()
	}

	userStatistics, err := s.repository.CountProcessingRequestsForIntentByUser(req)
	if err != nil {
		return nil, err
	}

	reply := &processingpb.ProcessingRequestIntentStatisticsReply{
		IntentId:       req.IntentId,
		UserStatistics: userStatistics,
	}

	return reply, nil
}
