package processor

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

const (
	ProcessorQueueGroup  = "processor"
	ReplyReceivedSubject = "processor.reply.received"
)

type Processor struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry intent.IntentRegistryClient
	storage  Storage
	wg       *sync.WaitGroup
}

func New(logger log.Logger, queue *nats.EncodedConn, registry intent.IntentRegistryClient, storage Storage) Processor {
	p := Processor{
		logger:   log.With(logger, "component", "processor"),
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

	p.logger.Log("msg", "started")

	return nil
}

func (p Processor) Stop() {
	p.logger.Log("msg", "stopping")
	p.wg.Wait()
}

func (p Processor) processRequest(request *processing.MessageRequest) {
	p.logger.Log("request received", request)

	query := &intent.QueryIntentsRequest{Server: request.ServerId}

	reply, err := p.registry.QueryIntents(context.Background(), query)
	if err != nil {
		p.logger.Log("error", "error querying for intents: "+err.Error())
		return
	}

	for _, intent := range reply.Intents {
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
			p.logger.Log("error", err.Error())
		}

		if !isMatch {
			return
		}

		command.Id = xid.New().String()

		p.saveCommand(command)

		p.logger.Log("executing command", command)
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
			p.logger.Log("error", err.Error())
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

	p.logger.Log("command.event", e)

	err := p.storage.SaveCommandEvent(e)
	if err != nil {
		p.logger.Log("error", "error saving command event: "+err.Error())
	}
}

func (p Processor) saveErrorEvent(commandId string, commandError error) {
	p.saveCommandEvent(commandId, "error", commandError.Error())
}
