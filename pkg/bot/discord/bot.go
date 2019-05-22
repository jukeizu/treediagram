package discord

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
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
	Client  processingpb.ProcessingClient
	Queue   *nats.EncodedConn
	Logger  zerolog.Logger
}

func NewBot(token string, client processingpb.ProcessingClient, queue *nats.EncodedConn, logger zerolog.Logger) (Bot, error) {
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
	// ignore messages from bots
	if m.Author.Bot {
		return
	}

	request := &processingpb.MessageRequest{
		Id:        m.ID,
		Source:    "discord",
		Bot:       mapToPbUser(s.State.User),
		Author:    mapToPbUser(m.Author),
		ChannelId: m.ChannelID,
		ServerId:  m.GuildID,
		Servers:   mapToPbServers(m.Author.ID, d.Session.State.Guilds),
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

	d.Logger.Debug().Str("requestId", reply.Id).Msg("request sent")
}

func (d *bot) messageReplyReceived(r *processingpb.MessageReplyReceived) {
	d.Logger.Info().
		Str("messageReplyId", r.Id).
		Msg("message reply received")

	message, err := d.Client.GetMessageReply(context.Background(), &processingpb.MessageReplyRequest{Id: r.Id})
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", r.Id).
			Msg("error getting message reply")
		return
	}

	d.Logger.Info().
		Str("messageReplyId", message.Id).
		Str("processingRequestId", message.ProcessingRequestId).
		Str("channelId", message.ChannelId).
		Str("userId", message.UserId).
		Msg("starting processing for message reply")

	err = d.publishMessage(message)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", r.Id).
			Msg("error publishing message")
		return
	}

	d.Logger.Info().
		Str("messageReplyId", message.Id).
		Str("processingRequestId", message.ProcessingRequestId).
		Str("channelId", message.ChannelId).
		Str("userId", message.UserId).
		Msg("finished processing for message reply")
}

func (d *bot) publishMessage(messageReply *processingpb.MessageReply) error {

	response := contract.Response{}

	err := json.Unmarshal([]byte(messageReply.Content), &response)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %s", err.Error())
	}

	if len(response.Messages) < 1 {
		d.Logger.Debug().
			Str("messageReplyId", messageReply.Id).
			Msg("message reply contains no messages")

		return nil
	}

	for _, message := range response.Messages {
		channelId := messageReply.ChannelId

		if message.IsRedirect {
			id, err := d.getUserChannelId(messageReply.UserId)
			if err != nil {
				return err
			}

			if id == channelId {
				continue
			}
		}

		if message.IsPrivateMessage {
			id, err := d.getUserChannelId(messageReply.UserId)
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
		if err != nil {
			return err
		}
	}

	return nil
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
