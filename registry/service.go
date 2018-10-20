package registry

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/registration"
)

type service struct {
	CommandStorage CommandStorage
}

func NewService(commandStorage CommandStorage) pb.RegistrationServer {
	return &service{CommandStorage: commandStorage}
}

func (s *service) AddCommand(ctx context.Context, req *pb.AddCommandRequest) (*pb.AddCommandReply, error) {
	err := s.CommandStorage.Save(*req.Command)
	if err != nil {
		return nil, err
	}

	return &pb.AddCommandReply{Command: req.Command}, nil
}

func (s *service) DisableCommand(ctx context.Context, req *pb.DisableCommandRequest) (*pb.DisableCommandReply, error) {
	err := s.CommandStorage.Disable(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DisableCommandReply{Id: req.Id}, nil
}

func (s *service) QueryCommands(ctx context.Context, req *pb.QueryCommandsRequest) (*pb.QueryCommandsReply, error) {
	if req.PageSize < 1 {
		req.PageSize = 50
	}

	commands, err := s.CommandStorage.Query(*req)
	if err != nil {
		return nil, err
	}

	reply := &pb.QueryCommandsReply{
		Commands: commands,
	}

	if len(reply.Commands) == int(req.PageSize) {
		reply.HasMore = true
	}

	return reply, nil
}
