package request

import (
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
)

const (
	ProcessorQueueGroup   = "processor"
	CommandMatchedSubject = "commandMatched"
)

type Matched struct {
	logger  log.Logger
	queue   *nats.EncodedConn
	client  Client
	storage Storage
}

func NewCommandMatchedProcessor(logger log.Logger, queue *nats.EncodedConn, storage Storage) Matched {
	return Matched{
		logger:  logger,
		queue:   queue,
		client:  Client{},
		storage: storage,
	}
}

func (m Matched) Subscribe() error {
	if m.queue == nil {
		return nil
	}

	_, err := m.queue.QueueSubscribe(CommandMatchedSubject, ProcessorQueueGroup, m.process)
	if err != nil {
		return err
	}

	return nil
}

func (m Matched) process(pm processing.Match) {
	m.logger.Log("msg", "received matched command", "match", pm.Id)

	match, err := m.storage.Match(pm.Id)
	if err != nil {
		m.logger.Log("error", "could not retrieve match from db: "+err.Error())
	}

	reply, err := m.client.Do(match)
	if err != nil {
		m.logger.Log("error", err)
	}

	m.logger.Log("msg", "received reply", "reply", reply)
}
