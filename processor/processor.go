package processor

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

const (
	ProcessorQueueGroup        = "processor"
	ProcessorCommandQueueGroup = "processor.command"
	MessageReceivedSubject     = "message.received"
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
	_, err := p.queue.QueueSubscribe(MessageReceivedSubject, ProcessorQueueGroup, p.processMessage)
	if err != nil {
		return err
	}

	_, err = p.queue.QueueSubscribe(CommandReceivedSubject, ProcessorQueueGroup, p.processCommand)
	if err != nil {
		return err
	}

	_, err = p.queue.QueueSubscribe(CommandReceivedSubject, ProcessorCommandQueueGroup, p.saveCommand)
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

func (p Processor) processMessage(m Message) {
	p.wg.Add(1)
	go func(m Message) {
		defer p.wg.Done()
		p.logger.Log("message received", m)

		query := &registration.QueryIntentsRequest{Server: m.ServerId}

		reply, err := p.registry.QueryIntents(context.Background(), query)
		if err != nil {
			p.logger.Log("error", "could not query for intents: "+err.Error())
			return
		}

		p.publishCommands(m, reply.Intents)
	}(m)
}

func (p Processor) publishCommands(m Message, intents []*registration.Intent) {
	for _, intent := range intents {
		i := NewIntent(*intent)

		isMatch, err := i.IsMatch(m)
		if err != nil {
			p.logger.Log("error", err.Error())
		}

		if !isMatch {
			continue
		}

		command := Command{
			Id:      xid.New().String(),
			Message: m,
			Intent:  i,
		}

		p.logger.Log("message has intent", command)

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

		p.logger.Log("executing command", command)

		reply, err := command.Execute()
		if err != nil {
			p.saveErrorEvent(command.Id, err)
			return
		}

		reply.Id = xid.New().String()

		p.logger.Log("received reply", reply.Id, "command", command)

		err = p.storage.SaveReply(reply)
		if err != nil {
			p.saveErrorEvent(command.Id, err)
			return
		}

		err = p.queue.Publish(ReplyReceivedSubject, reply.Id)
		if err != nil {
			p.saveErrorEvent(command.Id, err)
			return
		}
	}(command)
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
