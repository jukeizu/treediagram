package request

import (
	"context"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
)

const (
	RequestReceivedSubject = "requestReceived"
)

type Received struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry registration.RegistrationClient
	storage  Storage
}

func NewCommandReceivedProcessor(logger log.Logger, queue *nats.EncodedConn, registrationClient registration.RegistrationClient, storage Storage) Received {
	return Received{logger: logger, queue: queue, registry: registrationClient, storage: storage}
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
	match, err := regexp.MatchString(command.Regex, req.Content)
	if err != nil {
		r.logger.Log("error", "regexp: "+err.Error())
		return
	}

	if !match {
		return
	}

	m := Match{
		Request: toRequest(req),
		Command: Command{
			Id:       command.Id,
			Endpoint: command.Endpoint,
		},
	}

	id, err := r.storage.SaveMatch(m)
	if err != nil {
		r.logger.Log("error", "db: error saving match: "+err.Error())
	}

	pm := &processing.Match{
		Id: id,
	}

	r.logger.Log("msg", "publishing a command match", "match", pm.Id, "command", command.Id)

	err = r.queue.Publish(CommandMatchedSubject, pm)
	if err != nil {
		r.logger.Log("error", "error publishing commandMatched: "+err.Error())
	}
}
