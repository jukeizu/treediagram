package builtin

import (
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/rs/zerolog"
)

type SelectServerHandler struct {
	logger       zerolog.Logger
	intentClient intentpb.IntentRegistryClient
}

func NewSelectServerHandler(logger zerolog.Logger, intentClient intentpb.IntentRegistryClient) SelectServerHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.selectserver").Logger()

	return SelectServerHandler{logger, intentClient}
}

func (h SelectServerHandler) SelectServer(request contract.Request) (*contract.Response, error) {

	return contract.StringResponse(""), nil
}
