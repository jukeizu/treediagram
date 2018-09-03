package publishing

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/publishing"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

type service struct {
	Queue          *nats.EncodedConn
	MessageStorage MessageStorage
}

func NewService(q *nats.EncodedConn, store MessageStorage) pb.PublishingServer {
	return &service{q, store}
}

func (s service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	id := xid.New().String()
	req.Message.Id = id

	err := s.MessageStorage.Save(req.Message)
	if err != nil {
		return nil, err
	}

	err = s.Queue.Publish(req.Message.Source, pb.QueueMessage{Id: id})

	return &pb.SendMessageReply{Id: id}, err
}
