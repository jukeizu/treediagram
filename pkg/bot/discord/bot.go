package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/internal"
	nats "github.com/nats-io/nats.go"
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
	Session     *discordgo.Session
	Application *discordgo.Application
	Client      processingpb.ProcessingClient
	Queue       *nats.EncodedConn
	Logger      zerolog.Logger
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
	session.State.MaxMessageCount = 20

	// Enable all intents to include privileged intents
	session.Identify.Intents = discordgo.IntentsAll

	application, _ := session.Application("@me")
	if err != nil {
		return &dh, err
	}
	dh.Application = application

	session.AddHandler(dh.messageCreate)
	session.AddHandler(dh.messageReactionAdd)
	session.AddHandler(dh.interactionCreate)

	dh.Session = session

	_, err = dh.Queue.QueueSubscribe(DiscordMessageReplyReceivedSubject, DiscordQueueGroup, dh.messageReplyReceived)
	if err != nil {
		return &dh, err
	}

	return &dh, nil
}

func (d *bot) Open() error {
	d.Logger.Info().Msg("session opening")

	err := d.Session.Open()
	if err != nil {
		return err
	}

	return d.registerApplicationCommands()
}

func (d *bot) Close() {
	d.Logger.Info().Msg("session closing")

	d.Session.Close()
}

func (d *bot) registerApplicationCommands() error {
	_, err := d.Session.ApplicationCommandBulkOverwrite(d.Session.State.User.ID, "", d.defaultGlobalCommands())
	return err
}

func (d *bot) defaultGlobalCommands() []*discordgo.ApplicationCommand {
	version := discordgo.ApplicationCommand{
		Name:        "version",
		Description: "Returns the current bot version",
		Version:     internal.Version,
	}

	return []*discordgo.ApplicationCommand{&version}
}

func (d *bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bots
	if m.Author.Bot {
		return
	}

	request := mapToPbMessageRequest(s.State, d.Application, m.Message)

	reply, err := d.Client.SendMessageRequest(context.Background(), request)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("id", request.Id).
			Msg("error sending message request")
		return
	}

	d.Logger.Debug().Str("requestId", reply.Id).Msg("request sent")
}

func (d *bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		return
	}

	user := d.getInteractionUser(i)

	if i.Type == discordgo.InteractionApplicationCommand {
		d.handleApplicationCommand(s, i)
		return
	}

	d.Logger.Info().
		Str("requestId", i.ID).
		Str("type", i.Type.String()).
		Interface("user", user).
		Msg("received interaction create")

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	pbInteraction := mapToPbInteraction(s.State, i)

	_, err := d.Client.SendInteraction(context.Background(), pbInteraction)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Msg("error sending interaction")
		return
	}

	customId := i.MessageComponentData().CustomID

	d.Logger.Info().
		Str("customId", customId).
		Interface("user", user).
		Msg("interaction request sent")

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Msg("error responding to interaction")
		return
	}

	d.Logger.Info().
		Str("customId", customId).
		Interface("user", user).
		Msg("responded to discord interaction")
}

func (d *bot) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	user := d.getInteractionUser(i)

	d.Logger.Info().
		Str("requestId", i.ID).
		Interface("user", user).
		Interface("data", data).
		Str("serverId", i.GuildID).
		Str("channelId", i.ChannelID).
		Msg("received application command")

	switch data.Name {
	case "version":
		d.sendStringCommandResponse(internal.Version, s, i)
	}
}

func (d *bot) messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := d.Session.State.Message(r.ChannelID, r.MessageID)
	if err == discordgo.ErrStateNotFound {
		d.Logger.Debug().
			Str("messageID", r.MessageID).
			Str("channelID", r.ChannelID).
			Msg("looking up reaction message from discord api")

		message, err = d.Session.ChannelMessage(r.ChannelID, r.MessageID)
	}
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Str("messageID", r.MessageID).
			Str("channelID", r.ChannelID).
			Msg("error looking up message")

		return
	}

	reaction := mapToPbReaction(s.State, r.MessageReaction, d.Application, message)

	_, err = d.Client.SendReaction(context.Background(), reaction)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Msg("error sending reaction")
		return
	}

	d.Logger.Debug().
		Str("emojiID", reaction.Emoji.Id).
		Str("emojiName", reaction.Emoji.Name).
		Str("messageId", reaction.MessageRequest.Id).
		Msg("reaction sent")
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

	err = d.processMessageReply(message)
	if err != nil {
		d.Logger.Error().
			Err(err).
			Str("messageReplyId", message.Id).
			Msg("failed to process message reply")

		d.Logger.Info().
			Str("messageReplyId", message.Id).
			Msg("sending error event")

		errorEvent := &processingpb.ProcessingEvent{
			Type:                "error",
			Description:         err.Error(),
			ProcessingRequestId: message.GetProcessingRequestId(),
		}

		_, err = d.Client.SendProcessingEvent(context.Background(), errorEvent)
		if err != nil {
			d.Logger.Error().Caller().Err(err).
				Str("id", r.Id).
				Msg("failed to send event")

			return
		}

		d.Logger.Info().
			Str("messageReplyId", message.Id).
			Msg("finished sending error event")

		return
	}
}

func (d *bot) processMessageReply(messageReply *processingpb.MessageReply) error {
	d.Logger.Info().
		Str("messageReplyId", messageReply.Id).
		Str("processingRequestId", messageReply.ProcessingRequestId).
		Str("channelId", messageReply.ChannelId).
		Str("userId", messageReply.UserId).
		Msg("starting processing for message reply")

	response := contract.Response{}

	err := json.Unmarshal([]byte(messageReply.Content), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %s", err.Error())
	}

	err = d.publishReactions(messageReply, response.Reactions)
	if err != nil {
		return fmt.Errorf("failed to publish reactions: %s", err.Error())
	}

	err = d.publishMessages(messageReply, response.Messages)
	if err != nil {
		return fmt.Errorf("failed to publish messages: %s", err.Error())
	}

	d.Logger.Info().
		Str("messageReplyId", messageReply.Id).
		Str("processingRequestId", messageReply.ProcessingRequestId).
		Str("channelId", messageReply.ChannelId).
		Str("userId", messageReply.UserId).
		Msg("finished processing for message reply")

	return nil
}

func (d *bot) publishMessages(messageReply *processingpb.MessageReply, messages []*contract.Message) error {
	for _, message := range messages {
		if message == nil {
			continue
		}

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

		if message.EditMessageId != "" {
			messageEdit, err := mapToMessageEdit(message, channelId)
			if err != nil {
				return err
			}

			d.Logger.Info().
				Str("messageReplyId", messageReply.Id).
				Str("channelId", channelId).
				Str("editMessageId", message.EditMessageId).
				Msg("sending message edit to discord")

			_, err = d.Session.ChannelMessageEditComplex(messageEdit)
			if err != nil {
				return err
			}

			continue
		}

		messageSend, err := mapToMessageSend(message)
		if err != nil {
			return err
		}

		d.Logger.Info().
			Str("messageReplyId", messageReply.Id).
			Str("channelId", channelId).
			Msg("sending message to discord")

		publishedMessage, err := d.Session.ChannelMessageSendComplex(channelId, messageSend)
		if err != nil {
			d.handleRestError(messageReply, err)
			return err
		}

		reactions := []*contract.Reaction{}

		for _, emojiId := range message.Reactions {
			reaction := &contract.Reaction{
				EmojiId:   emojiId,
				ChannelId: publishedMessage.ChannelID,
				MessageId: publishedMessage.ID,
			}

			reactions = append(reactions, reaction)
		}

		d.publishReactions(messageReply, reactions)
	}

	return nil
}

func (d *bot) publishReactions(messageReply *processingpb.MessageReply, reactions []*contract.Reaction) error {
	for _, reaction := range reactions {
		if reaction == nil {
			continue
		}

		d.Logger.Info().
			Str("messageReplyId", messageReply.Id).
			Str("channelId", reaction.ChannelId).
			Str("messageId", reaction.MessageId).
			Str("emojiId", reaction.EmojiId).
			Msg("sending message reaction to discord")

		err := d.Session.MessageReactionAdd(reaction.ChannelId, reaction.MessageId, reaction.EmojiId)
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

func (d *bot) handleRestError(messageReply *processingpb.MessageReply, responseError error) {
	restError := &discordgo.RESTError{}
	if !errors.As(responseError, &restError) {
		return
	}

	d.Logger.Info().
		Err(responseError).
		Str("processingRequestId", messageReply.ProcessingRequestId).
		Str("messageReplyId", messageReply.Id).
		Str("channelId", messageReply.ChannelId).
		Str("userId", messageReply.UserId).
		Msg("handling error received from discord")

	helpMessage := mapRestErrorToHelpMessage(restError, messageReply)
	if helpMessage == "" {
		d.Logger.Info().
			Err(responseError).
			Str("processingRequestId", messageReply.ProcessingRequestId).
			Str("messageReplyId", messageReply.Id).
			Str("channelId", messageReply.ChannelId).
			Str("userId", messageReply.UserId).
			Msg("no help message required")
		return
	}

	_, err := d.Session.ChannelMessageSend(messageReply.ChannelId, helpMessage)
	if err != nil {
		d.Logger.Error().Caller().Err(err).
			Err(responseError).
			Str("processingRequestId", messageReply.ProcessingRequestId).
			Str("messageReplyId", messageReply.Id).
			Str("channelId", messageReply.ChannelId).
			Str("userId", messageReply.UserId).
			Msg("error sending discord error help message")
		return
	}

	d.Logger.Info().
		Err(responseError).
		Str("processingRequestId", messageReply.ProcessingRequestId).
		Str("messageReplyId", messageReply.Id).
		Str("channelId", messageReply.ChannelId).
		Str("userId", messageReply.UserId).
		Msg("help message sent to user")
}

func (d *bot) discordLogger(dgoLevel int, caller int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	d.Logger.WithLevel(mapToLevel(dgoLevel)).
		Str("component", "discordgo").
		Str("version", discordgo.VERSION).
		Msg(message)
}

func (d *bot) sendStringCommandResponse(content string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	d.sendCommandResponse(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}, s, i)
}

func (d *bot) sendCommandResponse(response *discordgo.InteractionResponse, s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	user := d.getInteractionUser(i)

	d.Logger.Debug().
		Interface("data", data).
		Str("requestId", i.ID).
		Interface("user", user).
		Str("serverId", i.GuildID).
		Str("channelId", i.ChannelID).
		Msg("sending interaction response")

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		d.Logger.Error().Caller().
			Err(err).
			Interface("data", data).
			Str("requestId", i.ID).
			Interface("user", user).
			Str("serverId", i.GuildID).
			Str("channelId", i.ChannelID).
			Msg("failed to send interaction response")
		return
	}

	d.Logger.Info().
		Interface("data", data).
		Str("requestId", i.ID).
		Interface("user", user).
		Str("serverId", i.GuildID).
		Str("channelId", i.ChannelID).
		Msg("interaction response sent")
}

func (d *bot) getInteractionUser(i *discordgo.InteractionCreate) *discordgo.User {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User
	}

	return i.User
}
