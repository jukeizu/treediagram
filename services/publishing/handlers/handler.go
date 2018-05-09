package handlers

import (
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
)

type HandlerParams struct {
	MessageRequest storage.MessageRequest
	Message        storage.Message
}

type MessageHandler interface {
	Handle(HandlerParams) error
}

type Handler interface {
	queue.QueueHandler
}

type handler struct {
	Storage        storage.Storage
	MessageHandler MessageHandler
}

func NewQueueHandler(s storage.Storage, messageHandler MessageHandler) Handler {
	return &handler{s, messageHandler}
}

func (h *handler) Handle(queueMessage queue.QueueMessage) error {

	messageRequest, err := h.Storage.MessageRequestStorage.GetMessageRequest(queueMessage.Id)

	if err != nil {
		return err
	}

	message, err := h.Storage.MessageStorage.GetMessage(messageRequest.MessageId)

	if err != nil {
		return err
	}

	params := HandlerParams{messageRequest, message}

	return h.MessageHandler.Handle(params)
}
