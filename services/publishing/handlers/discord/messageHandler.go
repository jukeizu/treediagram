package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/treediagram/services/publishing/handlers"
)

type DiscordConfig struct {
	Token string
}

type DiscordHandler interface {
	handlers.MessageHandler
}

type discordHandler struct {
	Session *discordgo.Session
}

func NewDiscordHandler(config DiscordConfig) (DiscordHandler, error) {
	handler := discordHandler{}

	session, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		return &handler, err
	}

	handler.Session = session

	return &handler, nil
}

func (h *discordHandler) Handle(params handlers.HandlerParams) error {

	channelId := params.MessageRequest.ChannelId

	if params.MessageRequest.IsRedirect {

		id, err := h.getUserChannelId(params.MessageRequest.User.Id)

		if err != nil {
			return err
		}

		if id == channelId {
			return nil
		}
	}

	if params.MessageRequest.PrivateMessage {
		id, err := h.getUserChannelId(params.MessageRequest.User.Id)

		if err != nil {
			return err
		}

		channelId = id
	}

	messageSend := MapToMessageSend(params.Message)

	_, err := h.Session.ChannelMessageSendComplex(channelId, messageSend)

	return err
}

func (h *discordHandler) getUserChannelId(userId string) (string, error) {
	dmChannel, err := h.Session.UserChannelCreate(userId)

	if err != nil {
		return "", err
	}

	return dmChannel.ID, nil
}
