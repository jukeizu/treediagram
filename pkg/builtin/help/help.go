package help

import (
	"sort"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/jukeizu/treediagram/pkg/intent"
	"github.com/rs/zerolog"
)

type HelpHandler struct {
	logger   zerolog.Logger
	registry *intent.Registry
}

func NewHelpHandler(logger zerolog.Logger, registry *intent.Registry) HelpHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.help").Logger()

	return HelpHandler{logger, registry}
}

func (h HelpHandler) Registrations() []builtin.HandlerRegistration {
	return []builtin.HandlerRegistration{
		{Name: "help", Handler: h.Help},
	}
}

func (h HelpHandler) Help(request contract.Request) (*contract.Response, error) {
	query := intent.Query{
		ServerId: request.ServerId,
		Type:     "command",
	}

	intents := h.registry.Query(query)

	embed := &contract.Embed{
		Title:       "treediagram",
		Description: "help",
		Color:       6139372,
	}

	for _, intent := range intents {
		if intent.Name == "" || intent.Help == "" {
			continue
		}

		field := &contract.EmbedField{
			Name:  intent.Name,
			Value: intent.Help,
		}

		embed.Fields = append(embed.Fields, field)
	}

	sort.Slice(embed.Fields, func(i, j int) bool {
		return embed.Fields[i].Name < embed.Fields[j].Name
	})

	return &contract.Response{Messages: []*contract.Message{{Embed: embed}}}, nil
}
