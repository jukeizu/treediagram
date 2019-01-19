package processor

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

const (
	ProcessorQueueGroup  = "processor"
	ReplyReceivedSubject = "processor.reply.received"
)

type Processor struct {
	logger   zerolog.Logger
	queue    *nats.EncodedConn
	registry intent.IntentRegistryClient
	storage  Storage
	wg       *sync.WaitGroup
}

func New(logger zerolog.Logger, queue *nats.EncodedConn, registry intent.IntentRegistryClient, storage Storage) Processor {
	p := Processor{
		logger:   logger.With().Str("component", "processor").Logger(),
		queue:    queue,
		registry: registry,
		storage:  storage,
		wg:       &sync.WaitGroup{},
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
	p.wg.Wait()
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
	p.wg.Add(1)
	go func(command Command) {
		defer p.wg.Done()

		isMatch, err := command.IsMatch()
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("could not determine if command is match")
		}

		if !isMatch {
			return
		}

		command.Id = xid.New().String()

		p.saveCommand(command)

		p.logger.Debug().Msg("executing command")
		response, err := command.Execute()
		if err != nil {
			p.saveErrorEvent(command.Id, err)
			return
		}

		for _, message := range response.Messages {
			p.saveResponseMessage(command, *message)
		}
	}(command)
}

func (p Processor) saveResponseMessage(command Command, message processing.Message) {
	messageReply := processing.MessageReply{
		Id:               xid.New().String(),
		CommandId:        command.Id,
		ChannelId:        command.Request.ChannelId,
		UserId:           command.Request.Author.Id,
		IsPrivateMessage: message.IsPrivateMessage,
		IsRedirect:       message.IsRedirect,
		Content:          message.Content,
		Embed:            message.Embed,
		Tts:              message.Tts,
		Files:            message.Files,
	}

	err := p.storage.SaveMessageReply(messageReply)
	if err != nil {
		p.saveErrorEvent(command.Id, err)
	}

	messageReplyReceived := processing.MessageReplyReceived{Id: messageReply.Id}

	err = p.queue.Publish(ReplyReceivedSubject+"."+command.Request.Source, messageReplyReceived)
	if err != nil {
		p.saveErrorEvent(command.Id, err)
	}
}

func (p Processor) saveCommand(command Command) {
	p.wg.Add(1)
	go func(command Command) {
		p.wg.Done()

		err := p.storage.SaveCommand(command)
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("error saving command")
		}
	}(command)
}

func (p Processor) saveCommandEvent(commandId string, t string, d string) {
	e := CommandEvent{
		CommandId:   commandId,
		Type:        t,
		Description: d,
		Timestamp:   time.Now().Unix(),
	}

	err := p.storage.SaveCommandEvent(e)
	if err != nil {
		p.logger.Error().Err(err).Caller().Msg("error saving command event")
	}
}

func (p Processor) saveErrorEvent(commandId string, commandError error) {
	p.saveCommandEvent(commandId, "error", commandError.Error())
}
