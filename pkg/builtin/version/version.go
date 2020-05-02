package version

import (
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/internal"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/rs/zerolog"
)

type VersionHandler struct {
	logger zerolog.Logger
}

func NewVersionHandler(logger zerolog.Logger) VersionHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.version").Logger()

	return VersionHandler{logger}
}

func (h VersionHandler) Registrations() []builtin.HandlerRegistration {
	return []builtin.HandlerRegistration{
		builtin.HandlerRegistration{Name: "version", Handler: h.Version},
	}
}

func (h VersionHandler) Version(request contract.Request) (*contract.Response, error) {
	return contract.StringResponse(internal.Version), nil
}
