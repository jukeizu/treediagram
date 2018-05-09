package publishing

import (
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
	"github.com/rs/xid"
)

type Service interface {
	SendMessage(SendMessageRequest) (Response, error)
}

type service struct {
	Queue   queue.Queue
	Storage storage.Storage
}

func NewService(q queue.Queue, store storage.Storage) Service {
	return &service{q, store}
}

func (s *service) SendMessage(sendMessageRequest SendMessageRequest) (Response, error) {
	response := Response{}

	messageRequestId, err := s.saveSendMessageRequest(sendMessageRequest)

	if err != nil {
		return response, err
	}

	response.Id = messageRequestId

	queueMessage := queue.QueueMessage{Id: messageRequestId}

	err = s.Queue.PublishMessageRequest(queueMessage)

	return response, err
}

func (s *service) saveSendMessageRequest(sendMessageRequest SendMessageRequest) (string, error) {
	message := sendMessageRequest.Message
	message.Id = generateId()

	err := s.Storage.MessageStorage.SaveMessage(message)

	if err != nil {
		return "", err
	}

	messageRequest := storage.MessageRequest{}
	messageRequest.Id = generateId()
	messageRequest.CorrelationId = sendMessageRequest.CorrelationId
	messageRequest.ChannelId = sendMessageRequest.ChannelId
	messageRequest.User = sendMessageRequest.User
	messageRequest.PrivateMessage = sendMessageRequest.PrivateMessage
	messageRequest.IsRedirect = sendMessageRequest.IsRedirect
	messageRequest.MessageId = message.Id

	err = s.Storage.MessageRequestStorage.SaveMessageRequest(messageRequest)

	if err != nil {
		return "", err
	}

	return messageRequest.Id, nil
}

func generateId() string {
	return xid.New().String()
}
