package publishing

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/publishing"
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
	"github.com/rs/xid"
)

type service struct {
	Queue          queue.Queue
	MessageStorage storage.MessageStorage
}

func NewService(q queue.Queue, store storage.MessageStorage) pb.PublishingServer {
	return &service{q, store}
}

func (s service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	req.Message.Id = xid.New().String()

	err := s.MessageStorage.Save(req.Message)
	if err != nil {
		return nil, err
	}

	queueMessage := queue.QueueMessage{Id: req.Message.Id}
	err = s.Queue.PublishMessageRequest(queueMessage)

	return &pb.SendMessageReply{Id: req.Message.Id}, err
}
