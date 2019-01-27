package processor

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/user"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
)

const (
	ProcessorQueueGroup  = "processor"
	ReplyReceivedSubject = "processor.reply.received"
)

type Processor struct {
	logger     zerolog.Logger
	queue      *nats.EncodedConn
	registry   intent.IntentRegistryClient
	userClient user.UserClient
	repository Repository
	waitGroup  *sync.WaitGroup
}

func New(logger zerolog.Logger, queue *nats.EncodedConn, registry intent.IntentRegistryClient, userClient user.UserClient, repository Repository) Processor {
	p := Processor{
		logger:     logger.With().Str("component", "processor").Logger(),
		queue:      queue,
		registry:   registry,
		userClient: userClient,
		repository: repository,
		waitGroup:  &sync.WaitGroup{},
	}

	return p
}

func (p Processor) Start() error {
	_, err := p.queue.QueueSubscribe(RequestReceivedSubject, ProcessorQueueGroup, p.processRequest)
	if err != nil {
		return err
	}

	p.logger.Info().Msg("started")

	return nil
}

func (p Processor) Stop() {
	p.logger.Info().Msg("stopping")
	p.waitGroup.Wait()
}

func (p Processor) processRequest(request *processing.MessageRequest) {
	p.logger.Debug().Msgf("request received %+v", request)

	serverId, err := p.findServerId(request)
	if err != nil {
		p.logger.Error().Err(err).Caller().
			Str("userId", request.Author.Id).
			Msg("could not determine serverId for user")
		return
	}

	if serverId == "" {
		server, err := p.setDefaultServer(request)
		if err != nil {
			p.logger.Error().Err(err).Caller().
				Str("userId", request.Author.Id).
				Msg("could not set default server")
			return
		}

		serverId = server.Id
	}

	query := &intent.QueryIntentsRequest{ServerId: serverId}

	stream, err := p.registry.QueryIntents(context.Background(), query)
	if err != nil {
		p.logger.Error().Err(err).Caller().Msg("error getting QueryIntents stream")
		return
	}

	for {
		intent, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("error receiving intent from stream")
			break
		}

		command := Command{
			Request: *request,
			Intent:  *intent,
		}

		p.processCommand(command)
	}
}

func (p Processor) processCommand(command Command) {
	p.waitGroup.Add(1)
	go func(command Command) {
		defer p.waitGroup.Done()

		isMatch, err := command.IsMatch()
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("could not determine if command is match")
			return
		}

		if !isMatch {
			return
		}

		processingRequest, err := p.createProcessingRequest(command)
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("error creating ProcessingRequest")
			return
		}

		p.logger.Debug().
			Str("processingRequest", processingRequest.Id).
			Msg("executing command")

		response, err := command.Execute()
		if err != nil {
			p.createErrorEvent(processingRequest.Id, err)
			return
		}

		for _, message := range response.Messages {
			p.saveResponseMessage(processingRequest, *message)
		}
	}(command)
}

func (p Processor) findServerId(request *processing.MessageRequest) (string, error) {
	if request.ServerId != "" {
		return request.ServerId, nil
	}

	userId := request.Author.Id

	p.logger.Debug().
		Str("userId", userId).
		Msg("looking up server preference")

	preferenceReply, err := p.userClient.Preference(context.Background(), &user.PreferenceRequest{UserId: userId})
	if err != nil {
		return "", err
	}

	preference := preferenceReply.Preference
	if preference == nil {
		p.logger.Debug().
			Str("userId", userId).
			Msg("user does not have a preferred server")

		return "", nil
	}

	p.logger.Debug().
		Str("userId", userId).
		Str("preferredServerId", preference.ServerId).
		Msg("found server preference for user")

	return preference.ServerId, nil
}

func (p Processor) setDefaultServer(request *processing.MessageRequest) (*processing.Server, error) {
	if len(request.Servers) < 1 {
		return nil, errors.New("no servers are available to select default")
	}

	userId := request.Author.Id
	server := request.Servers[0]

	setServerRequest := user.SetServerRequest{
		UserId:   userId,
		ServerId: server.Id,
	}

	_, err := p.userClient.SetServer(context.Background(), &setServerRequest)
	if err != nil {
		return nil, err
	}

	p.logger.Debug().
		Str("userId", userId).
		Str("serverId", server.Id).
		Msg("user server preference has been set to the first available server")

	return server, nil
}

func (p Processor) saveResponseMessage(processingRequest *processing.ProcessingRequest, message processing.Message) {
	messageReply := processing.MessageReply{
		ProcessingRequestId: processingRequest.Id,
		ChannelId:           processingRequest.ChannelId,
		UserId:              processingRequest.UserId,
		IsPrivateMessage:    message.IsPrivateMessage,
		IsRedirect:          message.IsRedirect,
		Content:             message.Content,
	}

	err := p.repository.SaveMessageReply(&messageReply)
	if err != nil {
		p.createErrorEvent(processingRequest.Id, err)
		return
	}

	messageReplyReceived := processing.MessageReplyReceived{Id: messageReply.Id}

	err = p.queue.Publish(ReplyReceivedSubject+"."+processingRequest.Source, messageReplyReceived)
	if err != nil {
		p.createErrorEvent(processingRequest.Id, err)
		return
	}
}

func (p Processor) createProcessingRequest(command Command) (*processing.ProcessingRequest, error) {
	processingRequest := &processing.ProcessingRequest{
		IntentId:  command.Intent.Id,
		Source:    command.Request.Source,
		ChannelId: command.Request.ChannelId,
		ServerId:  command.Request.ServerId,
		BotId:     command.Request.Bot.Id,
		UserId:    command.Request.Author.Id,
	}

	err := p.repository.SaveProcessingRequest(processingRequest)
	if err != nil {
		return nil, err
	}

	return processingRequest, nil
}

func (p Processor) createProcessingEvent(processingRequestId string, eventType string, eventDescription string) {
	e := processing.ProcessingEvent{
		ProcessingRequestId: processingRequestId,
		Type:                eventType,
		Description:         eventDescription,
	}

	err := p.repository.SaveProcessingEvent(&e)
	if err != nil {
		p.logger.Error().Err(err).Caller().Msg("error saving processing event")
		return
	}

	p.logger.Debug().
		Str("processingRequestId", processingRequestId).
		Str("type", eventType).
		Str("description", eventDescription).
		Msg("saved processing event")
}

func (p Processor) createErrorEvent(processingRequestId string, processingError error) {
	p.createProcessingEvent(processingRequestId, "error", processingError.Error())
	p.logger.Error().Err(processingError).Caller().
		Str("processingRequestId", processingRequestId).
		Msg("")
}
