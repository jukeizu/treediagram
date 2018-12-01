package intent

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
)

type service struct {
	IntentStorage IntentStorage
}

func NewService(commandStorage IntentStorage) pb.IntentRegistryServer {
	return &service{IntentStorage: commandStorage}
}

func (s *service) AddIntent(ctx context.Context, req *pb.AddIntentRequest) (*pb.AddIntentReply, error) {
	err := s.IntentStorage.Save(*req.Intent)
	if err != nil {
		return nil, err
	}

	return &pb.AddIntentReply{Intent: req.Intent}, nil
}

func (s *service) DisableIntent(ctx context.Context, req *pb.DisableIntentRequest) (*pb.DisableIntentReply, error) {
	err := s.IntentStorage.Disable(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DisableIntentReply{Id: req.Id}, nil
}

func (s *service) QueryIntents(ctx context.Context, req *pb.QueryIntentsRequest) (*pb.QueryIntentsReply, error) {
	if req.PageSize < 1 {
		req.PageSize = 50
	}

	intents, err := s.IntentStorage.Query(*req)
	if err != nil {
		return nil, err
	}

	reply := &pb.QueryIntentsReply{
		Intents: intents,
	}

	numIntents := len(reply.Intents)

	if numIntents > 0 {
		reply.LastId = reply.Intents[numIntents-1].Id
	}

	if numIntents == int(req.PageSize) {
		reply.HasMore = true
	}

	return reply, nil
}
