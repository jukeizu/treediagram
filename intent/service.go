package intent

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
)

type service struct {
	Repository Repository
}

func NewService(repository Repository) intentpb.IntentRegistryServer {
	return &service{Repository: repository}
}

func (s *service) AddIntent(ctx context.Context, req *intentpb.AddIntentRequest) (*intentpb.AddIntentReply, error) {
	err := s.Repository.Save(req.Intent)
	if err != nil {
		return nil, err
	}

	return &intentpb.AddIntentReply{Intent: req.Intent}, nil
}

func (s *service) DisableIntent(ctx context.Context, req *intentpb.DisableIntentRequest) (*intentpb.DisableIntentReply, error) {
	err := s.Repository.Disable(req.Id)
	if err != nil {
		return nil, err
	}

	return &intentpb.DisableIntentReply{Id: req.Id}, nil
}

func (s *service) QueryIntents(req *intentpb.QueryIntentsRequest, stream intentpb.IntentRegistry_QueryIntentsServer) error {
	intents, err := s.Repository.Query(*req)
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
