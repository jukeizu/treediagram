package discord

import (
	"github.com/bwmarrin/discordgo"
	pb "github.com/jukeizu/treediagram/api/publishing"
)

var DiscordPublisherQueueGroup = "discord"

type DiscordConfig struct {
	Token string
}

type DiscordPublisher interface {
	Publish(*pb.Message) error
}

type discordPublisher struct {
	Session *discordgo.Session
}

func NewDiscordPublisher(config DiscordConfig) (DiscordPublisher, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	return &discordPublisher{Session: session}, nil
}

func (h *discordPublisher) Publish(message *pb.Message) error {
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

func (h *discordPublisher) getUserChannelId(userId string) (string, error) {
	dmChannel, err := h.Session.UserChannelCreate(userId)

	if err != nil {
		return "", err
	}

	return dmChannel.ID, nil
}
