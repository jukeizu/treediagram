package processor

import (
	"context"
	"io"
	"sync"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
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
	repository Repository
	waitGroup  *sync.WaitGroup
}

func New(logger zerolog.Logger, queue *nats.EncodedConn, registry intent.IntentRegistryClient, repository Repository) Processor {
	p := Processor{
		logger:     logger.With().Str("component", "processor").Logger(),
		queue:      queue,
		registry:   registry,
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

	query := &intent.QueryIntentsRequest{ServerId: request.ServerId}

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
