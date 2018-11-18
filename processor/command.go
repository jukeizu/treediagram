package processor

import (
	"errors"
	"regexp"

	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
)

type Command struct {
	Id       string `json:"id"`
	Server   string `json:"server"`
	Name     string `json:"name"`
	Regex    string `json:"regex"`
	Response string `json:"response"`
	Endpoint string `json:"endpoint"`
	Help     string `json:"help"`
	Enabled  bool   `json:"enabled"`
}

func NewCommand(rc registration.Command) Command {
	c := Command{
		Id:       rc.Id,
		Server:   rc.Server,
		Name:     rc.Name,
		Regex:    rc.Regex,
		Response: rc.Response,
		Endpoint: rc.Endpoint,
		Help:     rc.Help,
		Enabled:  rc.Enabled,
	}

	return c
}

func (c Command) IsMatch(m Message) (bool, error) {
	match, err := regexp.MatchString(c.Regex, m.Content)
	if err != nil {
		return match, errors.New("regexp: " + err.Error())
	}

	return match, nil
}
