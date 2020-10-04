package processor

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/rs/zerolog"
)

type Reaction struct {
	Request *processingpb.Reaction `json:"request"`
	Intent  *contract.Intent       `json:"intent"`
}

func (c Reaction) ShouldExecute() (bool, error) {
	if c.Intent.Type != "reaction" {
		return false, nil
	}

	match, err := regexp.MatchString(c.Intent.Regex, c.Request.Emoji.Name)
	if err != nil {
		return match, errors.New("regexp: " + err.Error())
	}

	return match, nil
}

func (c Reaction) Execute() (*processingpb.Response, error) {
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

	if c.Intent.Response != nil {
		content, err := json.Marshal(c.Intent.Response)
		if err != nil {
			return nil, err
		}

		reply.Messages = append(reply.Messages, string(content))
	}

	return reply, nil
}

func (c Reaction) ProcessingRequest() *processingpb.ProcessingRequest {
	processingRequest := &processingpb.ProcessingRequest{
		Type:      "reaction",
		IntentId:  c.Intent.Id,
		Source:    c.Request.MessageRequest.Source,
		ChannelId: c.Request.ChannelId,
		ServerId:  c.Request.ServerId,
		BotId:     c.Request.MessageRequest.Bot.Id,
		UserId:    c.Request.MessageRequest.Author.Id,
	}

	return processingRequest
}

func (c Reaction) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", "reaction").
		Str("intentId", c.Intent.Id).
		Str("intentName", c.Intent.Name).
		Str("source", c.Request.MessageRequest.Source).
		Str("channelId", c.Request.ChannelId).
		Str("serverId", c.Request.ServerId).
		Str("botId", c.Request.MessageRequest.Bot.Id).
		Str("userId", c.Request.MessageRequest.Author.Id)
}
