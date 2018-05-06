package queue

import (
	"github.com/go-kit/kit/log"
	"time"
)

type queueLogger struct {
	Logger log.Logger
	Queue  Queue
}

func NewQueueLogger(logger log.Logger, q Queue) Queue {
	return &queueLogger{logger, q}
}

func (q *queueLogger) PublishMessageRequest(queueMessage QueueMessage) error {

	defer func(begin time.Time) {
		q.Logger.Log(
			"method", "PublishMessageRequest",
			"queueMessage", queueMessage,
			"took", time.Since(begin),
		)
	}(time.Now())

	return q.Queue.PublishMessageRequest(queueMessage)
}

func (q *queueLogger) Listen(queueHandler QueueHandler) <-chan error {

	defer func(begin time.Time) {
		q.Logger.Log(
			"method", "Starting queue listeners.",
			"took", time.Since(begin),
		)
	}(time.Now())

	errors := q.Queue.Listen(queueHandler)

	go func() {
		for err := range errors {
			q.Logger.Log(
				"error", err.Error(),
			)
		}
	}()

	return errors
}

func (q *queueLogger) Close() {
	defer func(begin time.Time) {
		q.Logger.Log(
			"method", "Close",
			"took", time.Since(begin),
		)
	}(time.Now())

	q.Queue.Close()
}
