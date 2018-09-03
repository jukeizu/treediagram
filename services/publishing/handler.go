package publishing

import (
	pb "github.com/jukeizu/treediagram/api/publishing"
	nats "github.com/nats-io/go-nats"
)

type HandlerFunc func(*pb.Message) error

type Handler interface {
	Handle(subject string, handlerFunc HandlerFunc) (*nats.Subscription, error)
}

type handler struct {
	MessageStorage MessageStorage
	Queue          *nats.EncodedConn
}

func NewQueueHandler(s MessageStorage, q *nats.EncodedConn) Handler {
	return &handler{s, q}
}

func (h *handler) Handle(subject string, handlerFunc HandlerFunc) (*nats.Subscription, error) {
	sub, err := h.Queue.Subscribe(subject, func(queueMessage pb.QueueMessage) error {
		return h.process(queueMessage, handlerFunc)
	})

	return sub, err
}

func (h *handler) process(queueMessage pb.QueueMessage, handlerFunc HandlerFunc) error {
	message, err := h.MessageStorage.Message(queueMessage.Id)
	if err != nil {
		return err
	}

	return handlerFunc(message)
}
