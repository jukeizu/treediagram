package discord

import (
	"context"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	nats "github.com/nats-io/go-nats"
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
	Session   *discordgo.Session
	Client    pb.ProcessingClient
	Queue     *nats.EncodedConn
	Logger    log.Logger
	WaitGroup *sync.WaitGroup
}

func NewBot(token string, client pb.ProcessingClient, queue *nats.EncodedConn, logger log.Logger) (Bot, error) {
	dh := bot{
		Client:    client,
		Logger:    logger,
		Queue:     queue,
		WaitGroup: &sync.WaitGroup{},
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
	d.Logger.Log("session", "opening")

	return d.Session.Open()
}

func (d *bot) Close() {
	d.Logger.Log("session", "closing")

	d.WaitGroup.Wait()

	d.Session.Close()
}

func (d *bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	request := &pb.Request{
		Id:        m.ID,
		Source:    "discord",
		Bot:       mapToPbUser(s.State.User),
		Author:    mapToPbUser(m.Author),
		ChannelId: m.ChannelID,
		ServerId:  m.GuildID,
		Content:   m.Content,
		Mentions:  mapToPbUsers(m.Mentions),
	}

	reply, err := d.Client.SendRequest(context.Background(), request)
	if err != nil {
		d.Logger.Log("error", err.Error(), "id", request.Id)
		return
	}

	d.Logger.Log("request sent", reply.Id)
}

func (d *bot) messageReplyReceived(r *pb.MessageReplyReceived) {
	d.WaitGroup.Add(1)
	go func(r *pb.MessageReplyReceived) {
		d.WaitGroup.Done()

	}(r)
	d.Logger.Log("msg", "reply received", "reply", r.Id)
	message, err := d.Client.GetMessage(context.Background(), &pb.MessageRequest{Id: r.Id})
	if err != nil {
		d.Logger.Log("error", err.Error(), "id", r.Id)
		return
	}

	err = d.publishMessage(message)
	if err != nil {
		d.Logger.Log("error", err.Error(), "id", r.Id)
		return
	}
}

func (d *bot) publishMessage(message *pb.MessageReply) error {
	d.Logger.Log("msg", "received publish request", "message", message.Id)

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

	messageSend := mapToMessageSend(message)

	_, err := d.Session.ChannelMessageSendComplex(channelId, messageSend)

	return err
}

func (d *bot) getUserChannelId(userId string) (string, error) {
	dmChannel, err := d.Session.UserChannelCreate(userId)

	if err != nil {
		return "", err
	}

	return dmChannel.ID, nil
}

func (d *bot) discordLogger(level int, caller int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	d.Logger.Log("component", "discordgo", "level", level, "msg", message, "version", discordgo.VERSION)
}
