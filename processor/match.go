package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Match struct {
	Message Message `json:"message"`
	Command Command `json:"command"`
}

func (p Match) ExecuteCommand() (processing.Reply, error) {
	reply := processing.Reply{}

	reply.Results = []*processing.Result{&processing.Result{Content: p.Command.Response}}

	return reply, nil
}
