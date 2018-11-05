package command

import (
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
)

const (
	CommandReplyReceivedSubject = "commandReplyReceived"
)

type ReplyProcessor struct {
	logger  log.Logger
	queue   *nats.EncodedConn
	storage Storage
}

func NewCommandReplyProcessor(logger log.Logger, queue *nats.EncodedConn, storage Storage) ReplyProcessor {
	return ReplyProcessor{
		logger:  logger,
		queue:   queue,
		storage: storage,
	}
}

func (p ReplyProcessor) Process(pr processing.Reply, errors []error) error {
	r := Reply{
		ProcessingReply: pr,
	}

	for _, err := range errors {
		r.Errors = append(r.Errors, err.Error())
	}

	id, err := p.storage.SaveReply(r)
	if err != nil {
		return err
	}

	replyReceived := processing.ReplyReceived{
		Id: id,
	}

	return p.queue.Publish(CommandReplyReceivedSubject, replyReceived)
}
