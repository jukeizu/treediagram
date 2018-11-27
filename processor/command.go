package processor

import (
	"errors"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
)

type Command struct {
	Id      string                    `json:"id"`
	Request processing.MessageRequest `json:"request"`
	Intent  registration.Intent       `json:"intent"`
}

type CommandEvent struct {
	CommandId   string `json:"commandId"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Timestamp   int64  `json:"timestamp"`
}

func (c Command) IsMatch() (bool, error) {
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
