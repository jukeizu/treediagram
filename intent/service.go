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

func (s *service) QueryIntents(req *pb.QueryIntentsRequest, stream pb.IntentRegistry_QueryIntentsServer) error {
	intents, err := s.IntentStorage.Query(*req)
	if err != nil {
		return err
	}

	for _, intent := range intents {
		err := stream.Send(intent)
		if err != nil {
			return err
		}
	}

	return nil
}
