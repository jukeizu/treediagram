package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Command struct {
	Message Message `json:"message"`
	Intent  Intent  `json:"intent"`
}

func (c Command) Execute() (processing.Reply, error) {
	reply := processing.Reply{}

	reply.Results = []*processing.Result{&processing.Result{Content: c.Intent.Response}}

	return reply, nil
}
