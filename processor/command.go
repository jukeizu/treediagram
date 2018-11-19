package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Command struct {
	Id      string  `json:"id"`
	Request Request `json:"request"`
	Intent  Intent  `json:"intent"`
}

type CommandEvent struct {
	CommandId   string `json:"commandId"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Timestamp   int64  `json:"timestamp"`
}

func (c Command) Execute() (*processing.Reply, error) {
	reply := &processing.Reply{}

	if len(c.Intent.Endpoint) > 0 {
		client := Client{}

		r, err := client.Do(c)
		if err != nil {
			return reply, err
		}

		reply.Messages = r.Messages
	}

	if len(c.Intent.Response) > 0 {
		m := &processing.Message{Content: c.Intent.Response}
		reply.Messages = append(reply.Messages, m)
	}

	return reply, nil
}
