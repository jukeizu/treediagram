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

func (c Command) Execute() (processing.Reply, error) {
	reply := processing.Reply{}

	reply.Messages = []*processing.Message{&processing.Message{Content: c.Intent.Response}}

	return reply, nil
}
