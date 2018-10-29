package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
)

type Bot interface {
	Open() error
	Close()
}

type bot struct {
	Session *discordgo.Session
	Client  pb.ProcessingClient
	Logger  log.Logger
}

func NewBot(token string, client pb.ProcessingClient, logger log.Logger) (Bot, error) {
	dh := bot{Client: client, Logger: logger}

	discordgo.Logger = dh.discordLogger

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return &dh, err
	}

	session.LogLevel = discordgo.LogInformational

	session.AddHandler(dh.messageCreate)

	dh.Session = session

	return &dh, nil
}

func (d *bot) Open() error {
	d.Logger.Log("session", "opening")

	return d.Session.Open()
}

func (d *bot) Close() {
	d.Logger.Log("session", "closing")

	d.Session.Close()
}

func (d *bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	request := &pb.TreediagramRequest{}

	request.Source = "discord"
	request.CorrelationId = m.ID
	request.Bot = mapToUser(s.State.User)
	request.Author = mapToUser(m.Author)
	request.ChannelId = m.ChannelID
	request.ServerId = getServerId(s, m.ChannelID)
	request.Content = m.Content
	request.Mentions = mapToUsers(m.Mentions)

	_, err := d.Client.Request(context.Background(), request)

	if err != nil {
		d.Logger.Log("error", err.Error(), "correlationId", request.CorrelationId)
	}
}

func (d *bot) discordLogger(level int, caller int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	d.Logger.Log("component", "discordgo", "level", level, "msg", message, "version", discordgo.VERSION)
}