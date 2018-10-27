package processor

import (
	"context"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

const (
	RequestSubject = "RequestReceived"
)

type matchedCommand struct {
	Request processing.TreediagramRequest
	Command registration.Command
}

type service struct {
	Queue    *nats.EncodedConn
	Registry registration.RegistrationClient
}

func NewService(queue *nats.EncodedConn, registrationClient registration.RegistrationClient) processing.ProcessingServer {
	return &service{Queue: queue, Registry: registrationClient}
}

func (s service) Request(ctx context.Context, req *processing.TreediagramRequest) (*processing.TreediagramReply, error) {
	query := &registration.QueryCommandsRequest{
		Server: req.ServerId,
	}

	for {
		reply, err := s.Registry.QueryCommands(ctx, query)
		if err != nil {
			return nil, err
		}

		for _, command := range reply.Commands {
			go s.checkCommand(*req, *command)
		}

		if !reply.HasMore {
			break
		}

		query.LastId = reply.LastId
	}

	id := xid.New().String()

	return &processing.TreediagramReply{Id: id}, nil
}

func (s service) checkCommand(req processing.TreediagramRequest, command registration.Command) {
	match, _ := regexp.MatchString(command.Regex, req.Content)

	if match {
		s.Queue.Publish(RequestSubject, matchedCommand{Request: req, Command: command})
	}
}
