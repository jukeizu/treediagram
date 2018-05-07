package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api"
)

type DiscordSubscriberConfig struct {
	Token string
}

type DiscordSubscriber interface {
	Open() error
	Close()
}

type discordSubscriber struct {
	Session *discordgo.Session
	Client  api.Client
	Logger  log.Logger
}

func NewDiscordSubscriber(config DiscordSubscriberConfig, client api.Client, logger log.Logger) (DiscordSubscriber, error) {
	dh := discordSubscriber{Client: client, Logger: logger}

	session, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		return &dh, err
	}

	dh.Session = session

	session.AddHandler(dh.messageCreate)

	return &dh, nil
}

func (d *discordSubscriber) Open() error {
	return d.Session.Open()
}

func (d *discordSubscriber) Close() {
	d.Session.Close()
}

func (d *discordSubscriber) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	request := api.TreediagramRequest{}

	request.Source = "discord"
	request.CorrelationId = m.ID
	request.Bot = mapToUser(s.State.User)
	request.Author = mapToUser(m.Author)
	request.ChannelId = m.ChannelID
	request.ServerId = getServerId(s, m.ChannelID)
	request.Content = m.Content
	request.Mentions = mapToUsers(m.Mentions)

	response, err := d.Client.Treediagram().Request(request)

	if err != nil {
		d.Logger.Log("error", err.Error(), "correlationId", request.CorrelationId, "responseId", response.Id)
	}
}
