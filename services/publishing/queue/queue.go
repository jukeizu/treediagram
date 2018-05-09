package queue

import (
	"encoding/json"
	"github.com/shawntoffel/rabbitmq"
)

type Queue interface {
	PublishMessageRequest(QueueMessage) error
	Listen(QueueHandler) <-chan error
	Close()
}

type QueueHandler interface {
	Handle(QueueMessage) error
}

type queue struct {
	MessageQueue rabbitmq.Client
	QueueHandler QueueHandler
}

type QueueConfig struct {
	MessageQueueConfig rabbitmq.Config
}

type QueueMessage struct {
	Id string `json:"id"`
}

func NewQueue(queueConfig QueueConfig) (Queue, error) {
	q := queue{}

	messageQueue, err := rabbitmq.NewClient(queueConfig.MessageQueueConfig)
	if err != nil {
		return &q, err
	}

	q.MessageQueue = messageQueue
	return &q, nil
}

func (q *queue) PublishMessageRequest(queueMessage QueueMessage) error {
	return publish(q.MessageQueue, queueMessage)
}

func (q *queue) Listen(queueHandler QueueHandler) <-chan error {
	q.QueueHandler = queueHandler

	return q.MessageQueue.Listen(q.runAction)
}

func (q *queue) Close() {
	q.MessageQueue.Close()
}

func publish(client rabbitmq.Client, queueMessage QueueMessage) error {
	marshalled, err := json.Marshal(queueMessage)

	if err != nil {
		return err
	}

	return client.Publish(marshalled)
}

func (q *queue) runAction(bytes []byte) error {
	queueMessage := QueueMessage{}

	err := json.Unmarshal(bytes, &queueMessage)

	if err != nil {
		return err
	}

	return q.QueueHandler.Handle(queueMessage)
}
