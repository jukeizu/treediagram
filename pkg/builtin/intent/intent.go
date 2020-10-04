package intent

import (
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/jukeizu/treediagram/pkg/intent"
	"github.com/rs/zerolog"
)

type IntentHandler struct {
	logger   zerolog.Logger
	registry *intent.Registry
}

func NewIntentHandler(logger zerolog.Logger, registry *intent.Registry) IntentHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.intent").Logger()

	return IntentHandler{logger, registry}
}

func (h IntentHandler) Registrations() []builtin.HandlerRegistration {
	return []builtin.HandlerRegistration{
		{Name: "intentload", Handler: h.Load},
	}
}

func (h IntentHandler) Load(request contract.Request) (*contract.Response, error) {
	if !authorIsOwner(request) {
		return nil, nil
	}

	err := h.registry.Load()
	if err != nil {
		return contract.StringResponse("```" + err.Error() + "```"), nil
	}

	return contract.StringResponse("loaded registry!"), nil
}

func authorIsOwner(request contract.Request) bool {
	if request.Application.Owner.Id == request.Author.Id {
		return true
	}

	server := request.Server()

	return request.Author.Id == server.OwnerId
}
