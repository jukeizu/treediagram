package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
)

const (
	DiscordQueueGroup                  = "discord"
	DiscordMessageReplyReceivedSubject = "processor.reply.received.discord"
)

type Bot interface {
	Open() error
	Close()
}

type bot struct {
	Session *discordgo.Session
	Client  pb.ProcessingClient
	Queue   *nats.EncodedConn
	Logger  zerolog.Logger
}

func NewBot(token string, client pb.ProcessingClient, queue *nats.EncodedConn, logger zerolog.Logger) (Bot, error) {
	dh := bot{
		Client: client,
		Logger: logger,
		Queue:  queue,
	}

	discordgo.Logger = dh.discordLogger

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return &dh, err
	}

	session.LogLevel = discordgo.LogInformational

	session.AddHandler(dh.messageCreate)

	dh.Session = session

	_, err = dh.Queue.QueueSubscribe(DiscordMessageReplyReceivedSubject, DiscordQueueGroup, dh.messageReplyReceived)
	if err != nil {
		return &dh, err
	}

	return &dh, nil
}

func (d *bot) Open() error {
	d.Logger.Info().Msg("session opening")

	return d.Session.Open()
}

func (d *bot) Close() {
	d.Logger.Info().Msg("session closing")

	d.Session.Close()
}

func (d *bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	request := &pb.MessageRequest{
		Id:        m.ID,
		Source:    "discord",
		Bot:       mapToPbUser(s.State.User),
		Author:    mapToPbUser(m.Author),
		ChannelId: m.ChannelID,
		ServerId:  m.GuildID,
		Content:   m.Content,
		Mentions:  mapToPbUsers(m.Mentions),
	}

	reply, err := d.Client.SendMessageRequest(context.Background(), request)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", request.Id).
			Msg("error sending message request")
		return
	}

	d.Logger.Debug().Str("reply.id", reply.Id).Msg("request sent")
}

func (d *bot) messageReplyReceived(r *pb.MessageReplyReceived) {
	d.Logger.Debug().Str("reply", r.Id).Msg("reply received")
	message, err := d.Client.GetMessageReply(context.Background(), &pb.MessageReplyRequest{Id: r.Id})
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", r.Id).
			Msg("error getting message reply")
		return
	}

	err = d.publishMessage(message)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", r.Id).
			Msg("error publishing message")
		return
	}
}

func (d *bot) publishMessage(message *pb.MessageReply) error {
	d.Logger.Debug().Str("message.id", message.Id).Msg("received publish request")

	channelId := message.ChannelId

	if message.IsRedirect {
		id, err := d.getUserChannelId(message.UserId)
		if err != nil {
			return err
		}

		if id == channelId {
			return nil
		}
	}

	if message.IsPrivateMessage {
		id, err := d.getUserChannelId(message.UserId)
		if err != nil {
			return err
		}

		channelId = id
	}

	messageSend, err := mapToMessageSend(message)
	if err != nil {
		return err
	}

	_, err = d.Session.ChannelMessageSendComplex(channelId, messageSend)

	return err
}

func (d *bot) getUserChannelId(userId string) (string, error) {
	dmChannel, err := d.Session.UserChannelCreate(userId)

	if err != nil {
		return "", err
	}

	return dmChannel.ID, nil
}

func (d *bot) discordLogger(dgoLevel int, caller int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	d.Logger.WithLevel(mapToLevel(dgoLevel)).
		Str("component", "discordgo").
		Str("version", discordgo.VERSION).
		Msg(message)
}
