package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Match struct {
	Command Command `json:"command"`
	Request Request `json:"request"`
}

type Command struct {
	Id       string `json:"id"`
	Server   string `json:"server"`
	Name     string `json:"name"`
	Regex    string `json:"regex"`
	Endpoint string `json:"endpoint"`
	Help     string `json:"help"`
	Enabled  bool   `json:"enabled"`
}

func (c Command) Execute() (processing.Reply, error) {
	result := processing.Reply{}

	return result, nil
}
