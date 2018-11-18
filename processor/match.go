package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Match struct {
	Message Message `json:"message"`
	Intent  Intent  `json:"intent"`
}

func (p Match) ExecuteIntent() (processing.Reply, error) {
	reply := processing.Reply{}

	reply.Results = []*processing.Result{&processing.Result{Content: p.Intent.Response}}

	return reply, nil
}
