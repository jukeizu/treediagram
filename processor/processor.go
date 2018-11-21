package processor

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

const (
	ProcessorQueueGroup        = "processor"
	ProcessorCommandQueueGroup = "processor.command"
	CommandReceivedSubject     = "processor.command.received"
	ReplyReceivedSubject       = "processor.reply.received"
)

type Processor struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry registration.RegistrationClient
	storage  Storage
	wg       *sync.WaitGroup
}

func New(logger log.Logger, queue *nats.EncodedConn, registry registration.RegistrationClient, storage Storage) Processor {
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

	_, err = p.queue.QueueSubscribe(CommandReceivedSubject, ProcessorQueueGroup, p.processCommand)
	if err != nil {
		return err
	}

	/*
		_, err = p.queue.QueueSubscribe(CommandReceivedSubject, ProcessorCommandQueueGroup, p.saveCommand)
		if err != nil {
			return err
		}

	*/

	p.logger.Log("msg", "started")

	return nil
}

func (p Processor) Stop() {
	p.logger.Log("msg", "stopping")
	p.wg.Wait()
}

func (p Processor) processRequest(r Request) {
	p.wg.Add(1)
	go func(r Request) {
		defer p.wg.Done()
		p.logger.Log("request received", r)

		query := &registration.QueryIntentsRequest{Server: r.ServerId}

		reply, err := p.registry.QueryIntents(context.Background(), query)
		if err != nil {
			p.logger.Log("error", "could not query for intents: "+err.Error())
			return
		}

		p.publishCommands(r, reply.Intents)
	}(r)
}

func (p Processor) publishCommands(r Request, intents []*registration.Intent) {
	for _, intent := range intents {
		i := NewIntent(*intent)

		isMatch, err := i.IsMatch(r)
		if err != nil {
			p.logger.Log("error", err.Error())
		}

		if !isMatch {
			continue
		}

		command := Command{
			Id:      xid.New().String(),
			Request: r,
			Intent:  i,
		}

		p.logger.Log("request has intent", command)

		err = p.queue.Publish(CommandReceivedSubject, command)
		if err != nil {
			p.logger.Log("error", err.Error())
		}
	}
}

func (p Processor) processCommand(command Command) {
	p.wg.Add(1)
	go func(command Command) {
		defer p.wg.Done()

		err := p.storage.SaveCommand(command)
		if err != nil {
			p.logger.Log("error", err.Error())
		}

		p.logger.Log("executing command", command)

		response, err := command.Execute()
		if err != nil {
			p.saveErrorEvent(command.Id, err)
			return
		}

		p.saveResponseMessages(command, response)
	}(command)
}

func (p Processor) saveResponseMessages(command Command, response *processing.Response) {
	for _, message := range response.Messages {
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
}

func (p Processor) saveCommand(command Command) {
	err := p.storage.SaveCommand(command)
	if err != nil {
		p.logger.Log("error", err.Error())
	}
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
