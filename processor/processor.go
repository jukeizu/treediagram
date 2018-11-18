package processor

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
)

const (
	ProcessorQueueGroup    = "processor"
	MessageReceivedSubject = "message.received"
	IntentReceivedSubject  = "processor.intent.received"
	CommandReceivedSubject = "processor.command.received"
)

type Processor struct {
	logger   log.Logger
	queue    *nats.EncodedConn
	registry registration.RegistrationClient
	wg       *sync.WaitGroup
}

func New(logger log.Logger, queue *nats.EncodedConn, registry registration.RegistrationClient) Processor {
	p := Processor{
		logger:   log.With(logger, "component", "processor"),
		queue:    queue,
		registry: registry,
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
			Message: m,
			Intent:  i,
		}

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

		_, err := command.Execute()
		if err != nil {
			p.logger.Log("error", err.Error())
		}
	}(command)
}
