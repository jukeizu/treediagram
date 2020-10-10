package intent

import (
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/jukeizu/treediagram/pkg/intent"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type IntentHandler struct {
	logger   zerolog.Logger
	registry *intent.Registry
	queue    *nats.EncodedConn
}

func NewIntentHandler(logger zerolog.Logger, registry *intent.Registry, queue *nats.EncodedConn) IntentHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.intent").Logger()

	return IntentHandler{logger, registry, queue}
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

	err := h.queue.Publish(intent.LoadRegistrySubject, nil)
	if err != nil {
		return nil, err
	}

	return contract.StringResponse("load request has been sent!"), nil
}

func authorIsOwner(request contract.Request) bool {
	if request.Application.Owner.Id == request.Author.Id {
		return true
	}

	server := request.Server()

	return request.Author.Id == server.OwnerId
}
