package builtin

import (
	"bytes"
	"context"
	"flag"

	"encoding/json"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/rs/zerolog"
)

type IntentHandler struct {
	logger       zerolog.Logger
	intentClient intentpb.IntentRegistryClient
}

func NewIntentHandler(logger zerolog.Logger, intentClient intentpb.IntentRegistryClient) IntentHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.intent").Logger()

	return IntentHandler{logger, intentClient}
}

func (h IntentHandler) AddIntent(request contract.Request) (*contract.Response, error) {
	if !authorIsOwner(request) {
		return nil, nil
	}

	addIntentRequest, err := parseAddIntentRequest(request)
	if err != nil {
		return formatErrorResponse(err)
	}

	reply, err := h.intentClient.AddIntent(context.Background(), addIntentRequest)
	if err != nil {
		return nil, err
	}

	content, err := formatIntentReply(reply.Intent)
	if err != nil {
		return nil, err
	}

	return contract.StringResponse(content), nil
}

func (h IntentHandler) DisableIntent(request contract.Request) (*contract.Response, error) {
	if !authorIsOwner(request) {
		return nil, nil
	}

	intentID := parseIntentID(request)

	reply, err := h.intentClient.DisableIntent(context.Background(), &intentpb.DisableIntentRequest{Id: intentID})
	if err != nil {
		return nil, err
	}

	return contract.StringResponse("disabled intent: `" + reply.Id + "`"), nil
}

func authorIsOwner(request contract.Request) bool {
	server := request.Server()
	return request.Author.Id == server.OwnerId
}

func parseAddIntentRequest(request contract.Request) (*intentpb.AddIntentRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("addintent", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	name := parser.String("name", "", "The intent name.")
	t := parser.String("type", "command", "The intent type.")
	regex := parser.String("regex", "", "The regex used to match the intent.")
	help := parser.String("help", "", "The help text.")
	response := parser.String("response", "", "The intent response.")
	endpoint := parser.String("endpoint", "", "The intent endpoint.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	intent := &intentpb.Intent{
		ServerId: request.ServerId,
		Name:     *name,
		Type:     *t,
		Regex:    *regex,
		Help:     *help,
		Response: *response,
		Endpoint: *endpoint,
		Enabled:  true,
	}

	return &intentpb.AddIntentRequest{Intent: intent}, nil
}

func formatIntentReply(intent *intentpb.Intent) (string, error) {
	buffer := &bytes.Buffer{}
	buffer.WriteString("```json\n")

	j, err := json.MarshalIndent(intent, "", "  ")
	if err != nil {
		return "", err
	}

	buffer.Write(j)
	buffer.WriteString("```")

	return buffer.String(), nil
}
