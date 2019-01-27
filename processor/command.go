package processor

import (
	"errors"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
)

type Command struct {
	Request processingpb.MessageRequest `json:"request"`
	Intent  intentpb.Intent             `json:"intent"`
}

func (c Command) ShouldExecute() (bool, error) {
	if c.Intent.Mention && !c.isBotMentioned() {
		return false, nil
	}

	match, err := regexp.MatchString(c.Intent.Regex, c.Request.Content)
	if err != nil {
		return match, errors.New("regexp: " + err.Error())
	}

	return match, nil
}

func (c Command) Execute() (*processingpb.Response, error) {
	reply := &processingpb.Response{}

	if c.Intent.Endpoint != "" {
		client := Client{}

		r, err := client.Do(c.Request, c.Intent.Endpoint)
		if err != nil {
			return reply, err
		}

		reply.Messages = r.Messages
	}

	if c.Intent.Response != "" {
		m := &processingpb.Message{
			Content: c.Intent.Response,
		}
		reply.Messages = append(reply.Messages, m)
	}

	return reply, nil
}

func (c Command) ProcessingRequest() *processingpb.ProcessingRequest {
	processingRequest := &processingpb.ProcessingRequest{
		Type:      "command",
		IntentId:  c.Intent.Id,
		Source:    c.Request.Source,
		ChannelId: c.Request.ChannelId,
		ServerId:  c.Request.ServerId,
		BotId:     c.Request.Bot.Id,
		UserId:    c.Request.Author.Id,
	}

	return processingRequest
}

func (c Command) isBotMentioned() bool {
	for _, mention := range c.Request.Mentions {
		if mention.Id == c.Request.Bot.Id {
			return true
		}
	}

	return false
}
