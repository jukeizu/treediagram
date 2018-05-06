package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api"
)

type DiscordHandlerConfig struct {
	Token string
}

type DiscordHandler interface {
	Open() error
	Close()
}

type discordHandler struct {
	Session *discordgo.Session
	Client  api.Client
	Logger  log.Logger
}

func NewDiscordHandler(config DiscordHandlerConfig, client api.Client, logger log.Logger) (DiscordHandler, error) {
	dh := discordHandler{Client: client, Logger: logger}

	session, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		return &dh, err
	}

	dh.Session = session

	session.AddHandler(dh.messageCreate)

	return &dh, nil
}

func (d *discordHandler) Open() error {
	return d.Session.Open()
}

func (d *discordHandler) Close() {
	d.Session.Close()
}

func (d *discordHandler) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
