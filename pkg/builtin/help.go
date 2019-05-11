package builtin

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/rs/zerolog"
)

type HelpHandler struct {
	logger       zerolog.Logger
	intentClient intentpb.IntentRegistryClient
}

func NewHelpHandler(logger zerolog.Logger, intentClient intentpb.IntentRegistryClient) HelpHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.help").Logger()

	return HelpHandler{logger, intentClient}
}

func (h HelpHandler) Help(request contract.Request) (*contract.Response, error) {
	query := &intentpb.QueryIntentsRequest{ServerId: request.ServerId}

	intentStream, err := h.intentClient.QueryIntents(context.Background(), query)
	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("```")
	buffer.WriteString("\ntreediagram help\n\n")

	for {
		intent, err := intentStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			h.logger.Error().Err(err).Caller().Msg("error receiving intent from stream")
			break
		}

		buffer.WriteString(fmt.Sprintf("%s\n    %s\n", intent.Name, intent.Help))
	}

	buffer.WriteString("```")

	return contract.StringResponse(buffer.String()), nil
}
