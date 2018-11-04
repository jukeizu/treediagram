package command

import (
	"context"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
)

const (
	ProcessorQueueGroup    = "processor"
	RequestReceivedSubject = "requestReceived"
)

type Received struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry registration.RegistrationClient
}

func NewCommandReceivedProcessor(logger log.Logger, queue *nats.EncodedConn, registrationClient registration.RegistrationClient) Received {
	return Received{logger: logger, queue: queue, registry: registrationClient}
}

func (r Received) Subscribe() error {
	if r.queue == nil {
		return nil
	}

	_, err := r.queue.QueueSubscribe(RequestReceivedSubject, ProcessorQueueGroup, r.process)
	if err != nil {
		return err
	}

	return nil
}

func (r Received) process(req *processing.TreediagramRequest) {
	query := &registration.QueryCommandsRequest{Server: req.ServerId}

	for {
		reply, err := r.registry.QueryCommands(context.Background(), query)
		if err != nil {
			r.logger.Log("error", "error querying commands: "+err.Error())
			return
		}

		for _, command := range reply.Commands {
			go r.checkCommand(*req, *command)
		}

		if !reply.HasMore {
			break
		}

		query.LastId = reply.LastId
	}
}

func (r Received) checkCommand(req processing.TreediagramRequest, command registration.Command) {
	match, _ := regexp.MatchString(command.Regex, req.Content)

	if match {
		r.logger.Log("msg", "publishing a command match", "command", command.Id)

		err := r.queue.Publish(CommandMatchedSubject, matchedCommand{Request: req, Command: command})
		if err != nil {
			r.logger.Log("error", "error publishing commandMatched: "+err.Error())
		}
	}
}
