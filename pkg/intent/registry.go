package intent

import (
	"net/http"

	"github.com/jukeizu/contract"
	"gopkg.in/yaml.v2"
)

const (
	LoadRegistrySubject = "intent.registry.load.received"
)

type Registry struct {
	Url     string             `json:"url,omitempty" yaml:"url,omitempty"`
	Intents []*contract.Intent `json:"intents,omitempty" yaml:"intents,omitempty"`
}

func NewRegistry(url string) *Registry {
	return &Registry{
		Url:     url,
		Intents: []*contract.Intent{},
	}
}

func (r *Registry) Load() error {
	resp, err := http.Get(r.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	intents := []*contract.Intent{}

	err = yaml.NewDecoder(resp.Body).Decode(&intents)
	if err != nil {
		return err
	}

	r.Intents = intents

	return nil
}

func (r *Registry) Query(q Query) []*contract.Intent {
	if r.Intents == nil {
		return []*contract.Intent{}
	}

	reply := []*contract.Intent{}

	for _, intent := range r.Intents {
		if q.Match(intent) {
			reply = append(reply, intent)
		}
	}

	return reply
}
