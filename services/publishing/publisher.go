package publishing

import (
	pb "github.com/jukeizu/treediagram/api/publishing"
	nats "github.com/nats-io/go-nats"
)

var PublisherSubject = "publisher"

type PublisherFunc func(*pb.Message) error

type Publisher interface {
	Subscribe(subject string, publisherFunc PublisherFunc) (*nats.Subscription, error)
}

type publisher struct {
	messageStorage MessageStorage
	queue          *nats.EncodedConn
}

func NewPublisher(s MessageStorage, q *nats.EncodedConn) Publisher {
	return &publisher{s, q}
}

func (h *publisher) Subscribe(queue string, publisherFunc PublisherFunc) (*nats.Subscription, error) {
	sub, err := h.queue.QueueSubscribe(PublisherSubject, queue, func(queueMessage pb.PublishMessageRequestReceived) error {
		return h.process(queueMessage, publisherFunc)
	})

	return sub, err
}

func (h *publisher) process(queueMessage pb.PublishMessageRequestReceived, publisherFunc PublisherFunc) error {
	message, err := h.messageStorage.Message(queueMessage.Id)
	if err != nil {
		return err
	}

	return publisherFunc(message)
}
