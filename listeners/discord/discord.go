package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api"
)

type DiscordListenerConfig struct {
	Token string
}

type DiscordListener interface {
	Open() error
	Close()
}

type discordListener struct {
	Session *discordgo.Session
	Client  api.Client
	Logger  log.Logger
}

func NewDiscordListener(config DiscordListenerConfig, client api.Client, logger log.Logger) (DiscordListener, error) {
	dh := discordListener{Client: client, Logger: logger}

	discordgo.Logger = dh.discordLogger

	session, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		return &dh, err
	}

	session.LogLevel = discordgo.LogWarning

	session.AddHandler(dh.messageCreate)

	dh.Session = session

	return &dh, nil
}

func (d *discordListener) Open() error {
	return d.Session.Open()
}

func (d *discordListener) Close() {
	d.Session.Close()
}

func (d *discordListener) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func (d *discordListener) discordLogger(level int, caller int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	d.Logger.Log("component", "discordgo", "level", level, "message", message, "version", discordgo.VERSION)
}
