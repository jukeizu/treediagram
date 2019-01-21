package processor

import (
	"errors"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
)

type Command struct {
	Id      string                    `json:"id"`
	Request processing.MessageRequest `json:"request"`
	Intent  intent.Intent             `json:"intent"`
}

type ProcessingRequest struct {
	Id        string
	IntentId  string
	Source    string
	ChannelId string
	ServerId  string
	BotId     string
	UserId    string
}

type ProcessingEvent struct {
	Id          string
	RequestId   string
	Description string `json:"description"`
	Type        string `json:"type"`
	Timestamp   int64  `json:"timestamp"`
}

type CommandEvent struct {
	CommandId   string `json:"commandId"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Timestamp   int64  `json:"timestamp"`
}

func (c Command) IsMatch() (bool, error) {
	if c.Intent.Mention && !c.isBotMentioned() {
		return false, nil
	}

	match, err := regexp.MatchString(c.Intent.Regex, c.Request.Content)
	if err != nil {
		return match, errors.New("regexp: " + err.Error())
	}

	return match, nil
}

func (c Command) Execute() (*processing.Response, error) {
	reply := &processing.Response{}

	if len(c.Intent.Endpoint) > 0 {
		client := Client{}

		r, err := client.Do(c)
		if err != nil {
			return reply, err
		}

		reply.Messages = r.Messages
	}

	if len(c.Intent.Response) > 0 {
		m := &processing.Message{
			Content: c.Intent.Response,
		}
		reply.Messages = append(reply.Messages, m)
	}

	return reply, nil
}

func (c Command) isBotMentioned() bool {
	for _, mention := range c.Request.Mentions {
		if mention.Id == c.Request.Bot.Id {
			return true
		}
	}

	return false
}
