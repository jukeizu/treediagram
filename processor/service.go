package processor

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

const (
	RequestSubject = "treediagram.request"
)

type service struct {
	Queue *nats.EncodedConn
}

func NewService(queue *nats.EncodedConn) pb.ProcessingServer {
	return &service{queue}
}

func (s service) Request(ctx context.Context, req *pb.TreediagramRequest) (*pb.TreediagramReply, error) {
	id := xid.New().String()

	treediagramReply := &pb.TreediagramReply{Id: id}

	err := s.Queue.Publish(RequestSubject, req)

	return treediagramReply, err
}
