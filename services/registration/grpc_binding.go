package registration

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/registration"
)

type GrpcBinding struct {
	Service
}

func (b GrpcBinding) AddCommand(ctx context.Context, req *pb.AddCommandRequest) (*pb.AddCommandReply, error) {
	_, err := b.Service.Add(toCommand(req.Command))

	return &pb.AddCommandReply{}, err
}

func (b GrpcBinding) DisableCommand(ctx context.Context, req *pb.DisableCommandRequest) (*pb.DisableCommandReply, error) {
	err := b.Service.Disable(req.Id)

	return &pb.DisableCommandReply{}, err
}

func (b GrpcBinding) QueryCommands(ctx context.Context, req *pb.QueryCommandsRequest) (*pb.QueryCommandsReply, error) {
	queryResult, err := b.Service.Query(toCommandQuery(req))

	if err != nil {
		return nil, err
	}

	return toPbCommandQueryReply(queryResult), err
}

func toCommand(pbCommand *pb.Command) Command {
	command := Command{
		Server:         pbCommand.Server,
		Name:           pbCommand.Name,
		Regex:          pbCommand.Regex,
		RequireMention: pbCommand.RequireMention,
		Endpoint:       pbCommand.Endpoint,
		Help:           pbCommand.Help,
	}

	return command
}

func toCommandQuery(pbCommandQuery *pb.QueryCommandsRequest) CommandQuery {
	query := CommandQuery{
		Server:   pbCommandQuery.Server,
		LastId:   pbCommandQuery.LastId,
		PageSize: int(pbCommandQuery.PageSize),
	}

	return query
}

func toPbCommandQueryReply(queryResult CommandQueryResult) *pb.QueryCommandsReply {
	commandsReply := &pb.QueryCommandsReply{
		HasMore: queryResult.HasMore,
	}

	for _, command := range queryResult.Commands {
		mapped := toPbCommand(command)
		commandsReply.Commands = append(commandsReply.Commands, mapped)
	}

	return commandsReply
}

func toPbCommand(command Command) *pb.Command {
	pbCommand := &pb.Command{
		Id:             command.Id.String(),
		Server:         command.Server,
		Name:           command.Name,
		Regex:          command.Regex,
		RequireMention: command.RequireMention,
		Endpoint:       command.Endpoint,
		Help:           command.Help,
	}

	return pbCommand
}
