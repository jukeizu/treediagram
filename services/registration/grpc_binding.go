package registration

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/registration"
)

type grpcBinding struct {
	Service
}

func (b grpcBinding) Add(ctx context.Context, req *pb.AddCommandRequest) (*pb.AddCommandReply, error) {
	c := Command{
		Id:       req.Command.Id,
		Server:   req.Command.Server,
		Name:     req.Command.Name,
		Regex:    req.Command.Regex,
		Endpoint: req.Command.Endpoint,
		Help:     req.Command.Help,
	}

	_, err := b.Service.Add(c)

	return &pb.AddCommandReply{}, err
}

func (b grpcBinding) Remove(ctx context.Context, req *pb.RemoveCommandRequest) (*pb.RemoveCommandReply, error) {
	err := b.Service.Remove(req.Id)

	return &pb.RemoveCommandReply{}, err
}

func (b grpcBinding) Commands(ctx context.Context, req *pb.CommandsRequest) (*pb.CommandsReply, error) {
	commands, err := b.Service.Commands()

	if err != nil {
		return nil, err
	}

	commandsReply := &pb.CommandsReply{}

	for _, command := range commands {

		mapped := &pb.Command{
			Id:       command.Id,
			Server:   command.Server,
			Name:     command.Name,
			Regex:    command.Regex,
			Endpoint: command.Endpoint,
			Help:     command.Help,
		}

		commandsReply.Commands = append(commandsReply.Commands, mapped)
	}

	return commandsReply, err
}
