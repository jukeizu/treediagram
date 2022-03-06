package processor

import (
	"errors"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/rs/zerolog"
)

type Interaction struct {
	Request *processingpb.Interaction `json:"request"`
	Intent  *intentpb.Intent          `json:"intent"`
}

func (c Interaction) ShouldExecute() (bool, error) {
	if c.Intent.Type != "interaction" {
		return false, nil
	}

	match, err := regexp.MatchString(c.Intent.Regex, c.Request.Identifier)
	if err != nil {
		return match, errors.New("regexp: " + err.Error())
	}

	return match, nil
}

func (c Interaction) Execute() (*processingpb.Response, error) {
	reply := &processingpb.Response{}

	if c.Intent.Endpoint != "" {
		client := Client{}

		r, err := client.Do(c.Request, c.Intent.Endpoint)
		if err != nil {
			return reply, err
		}

		if r != "" {
			reply.Messages = append(reply.Messages, r)
		}

	}

	if c.Intent.Response != "" {
		reply.Messages = append(reply.Messages, c.Intent.Response)
	}

	return reply, nil
}

func (c Interaction) ProcessingRequest() *processingpb.ProcessingRequest {
	processingRequest := &processingpb.ProcessingRequest{
		Type:      "interaction",
		IntentId:  c.Intent.Id,
		Source:    c.Request.Source,
		ChannelId: c.Request.ChannelId,
		ServerId:  c.Request.ServerId,
		BotId:     c.Request.Bot.Id,
		UserId:    c.Request.User.Id,
	}

	return processingRequest
}

func (c Interaction) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", "interaction").
		Str("intentId", c.Intent.Id).
		Str("intentName", c.Intent.Name).
		Str("source", c.Request.Source).
		Str("channelId", c.Request.ChannelId).
		Str("serverId", c.Request.ServerId).
		Str("botId", c.Request.Bot.Id).
		Str("userId", c.Request.User.Id)
}
