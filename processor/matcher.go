package processor

import (
	"context"
	"errors"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
)

const (
	CommandMatchedSubject = "matched"
)

type Matcher struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry registration.RegistrationClient
}

func NewMatcher(logger log.Logger, registry registration.RegistrationClient) Matcher {
	return Matcher{logger: logger, registry: registry}
}

func (m Matcher) PublishMatches(request Request) error {
	query := &registration.QueryCommandsRequest{Server: request.ServerId}

	for {
		reply, err := m.registry.QueryCommands(context.Background(), query)
		if err != nil {
			return errors.New("error querying commands: " + err.Error())
		}

		for _, command := range reply.Commands {
			go m.checkMatch(request, command)
		}

		if !reply.HasMore {
			break
		}

		query.LastId = reply.LastId
	}

	return nil
}

func (m Matcher) checkMatch(request Request, rc *registration.Command) error {
	match, err := regexp.MatchString(rc.Regex, request.Content)
	if err != nil {
		return errors.New("regexp: " + err.Error())
	}

	if !match {
		return nil
	}

	err = m.queue.Publish(CommandMatchedSubject, request)
	if err != nil {
		errors.New("error publishing command match: " + err.Error())
	}

	return nil
}
