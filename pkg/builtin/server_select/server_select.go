package serverselect

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/rs/zerolog"
)

type ServerSelectHandler struct {
	logger     zerolog.Logger
	userClient userpb.UserClient
}

func NewServerSelectHandler(logger zerolog.Logger, userClient userpb.UserClient) ServerSelectHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.serverselect").Logger()

	return ServerSelectHandler{logger, userClient}
}

func (h ServerSelectHandler) Registrations() []builtin.HandlerRegistration {
	return []builtin.HandlerRegistration{
		builtin.HandlerRegistration{Name: "serverselect", Handler: h.ServerSelect},
	}
}

func (h ServerSelectHandler) ServerSelect(request contract.Request) (*contract.Response, error) {
	if len(request.Servers) < 1 {
		return nil, fmt.Errorf("server selection request did not contain any servers")
	}

	fields := strings.Fields(request.Content)

	if len(fields) == 1 {
		return contract.StringResponse(h.formatSelectionPrompt(request)), nil
	}

	selection, err := h.parseSelection(fields[1], len(request.Servers))
	if err != nil {
		return builtin.FormatErrorResponse(err)
	}

	server := request.Servers[selection-1]

	setServerRequest := userpb.SetServerRequest{
		UserId:   request.Author.Id,
		ServerId: server.Id,
	}

	_, err = h.userClient.SetServer(context.Background(), &setServerRequest)
	if err != nil {
		return nil, err
	}

	return contract.StringResponse(h.formatSelectionReply(server)), nil
}

func (h ServerSelectHandler) parseSelection(input string, serverCount int) (int, error) {
	selection, err := strconv.Atoi(input)
	if err != nil {
		return 0, builtin.NewParseError("Selection must be an integer")
	}

	if selection < 1 || selection > serverCount {
		return 0, builtin.NewParseError("That selection is not valid")
	}

	return selection, nil
}

func (h ServerSelectHandler) formatSelectionPrompt(request contract.Request) string {
	if len(request.Servers) == 1 {
		return fmt.Sprintf("You only share 1 server! Your preferred server is: %s", request.Servers[0].Name)
	}

	buffer := bytes.Buffer{}

	buffer.WriteString("Select a server to interact with:\n")

	for i, server := range request.Servers {
		buffer.WriteString(fmt.Sprintf("\n%d. %s `%s`", i+1, server.Name, server.Id))
		if server.Id == request.ServerId {
			buffer.WriteString(" (current server)")
		}
	}

	buffer.WriteString("\n\ne.g.`!server 1`")

	return buffer.String()
}

func (h ServerSelectHandler) formatSelectionReply(server contract.Server) string {
	return fmt.Sprintf("Your preferred server has been set to: %s `%s`", server.Name, server.Id)
}
