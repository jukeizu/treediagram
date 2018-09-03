package discord

import (
	"github.com/bwmarrin/discordgo"
	pb "github.com/jukeizu/treediagram/api/publishing"
)

var DiscordHandlerSubject = "discord"

type DiscordConfig struct {
	Token string
}

type DiscordHandler interface {
	Handle(*pb.Message) error
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

func (h *discordHandler) Handle(message *pb.Message) error {
	channelId := message.ChannelId

	if message.IsRedirect {
		id, err := h.getUserChannelId(message.User.Id)
		if err != nil {
			return err
		}

		if id == channelId {
			return nil
		}
	}

	if message.PrivateMessage {
		id, err := h.getUserChannelId(message.User.Id)
		if err != nil {
			return err
		}

		channelId = id
	}

	messageSend := MapToMessageSend(message)

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
