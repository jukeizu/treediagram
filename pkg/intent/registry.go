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
	url     string
	intents []*contract.Intent
}

func NewRegistry(url string) *Registry {
	return &Registry{
		url:     url,
		intents: []*contract.Intent{},
	}
}

func (r *Registry) Load() error {
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	intents := []*contract.Intent{}

	err = yaml.NewDecoder(resp.Body).Decode(&intents)
	if err != nil {
		return err
	}

	r.intents = intents

	return nil
}

func (r *Registry) Query(q Query) []*contract.Intent {
	if r.intents == nil {
		return []*contract.Intent{}
	}

	reply := []*contract.Intent{}

	for _, intent := range r.intents {
		if q.Match(intent) {
			reply = append(reply, intent)
		}
	}

	return reply
}
