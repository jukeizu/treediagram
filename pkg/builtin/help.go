package builtin

import (
	"context"
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

	embed := &contract.Embed{
		Title:       "treediagram",
		Description: "help",
		Color:       6139372,
	}

	for {
		intent, err := intentStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			h.logger.Error().Err(err).Caller().Msg("error receiving intent from stream")
			break
		}

		if intent.Name == "" || intent.Help == "" {
			continue
		}

		field := &contract.EmbedField{
			Name:  intent.Name,
			Value: intent.Help,
		}

		embed.Fields = append(embed.Fields, field)
	}

	return &contract.Response{Messages: []*contract.Message{&contract.Message{Embed: embed}}}, nil
}
