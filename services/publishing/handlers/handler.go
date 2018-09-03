package handlers

import (
	pb "github.com/jukeizu/treediagram/api/publishing"
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
)

type HandlerParams struct {
	Message *pb.Message
}

type MessageHandler interface {
	Handle(HandlerParams) error
}

type Handler interface {
	queue.QueueHandler
}

type handler struct {
	MessageStorage storage.MessageStorage
	MessageHandler MessageHandler
}

func NewQueueHandler(s storage.MessageStorage, messageHandler MessageHandler) Handler {
	return &handler{s, messageHandler}
}

func (h *handler) Handle(queueMessage queue.QueueMessage) error {
	message, err := h.MessageStorage.Message(queueMessage.Id)

	if err != nil {
		return err
	}

	params := HandlerParams{message}

	return h.MessageHandler.Handle(params)
}
