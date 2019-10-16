package builtin

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/rs/zerolog"
)

type StatsHandler struct {
	logger           zerolog.Logger
	processingClient processingpb.ProcessingClient
}

func NewStatsHandler(logger zerolog.Logger, processingClient processingpb.ProcessingClient) StatsHandler {
	logger = logger.With().Str("component", "intent.endpoint.builtin.stats").Logger()

	return StatsHandler{logger, processingClient}
}

func (h StatsHandler) Stats(request contract.Request) (*contract.Response, error) {
	intentID := parseIntentID(request)

	statsRequest := &processingpb.ProcessingRequestIntentStatisticsRequest{
		ServerId:  request.ServerId,
		IntentId:  intentID,
		UserLimit: 20,
		Type:      "command",
	}

	statsReply, err := h.processingClient.ProcessingRequestIntentStatistics(context.Background(), statsRequest)
	if err != nil {
		return nil, fmt.Errorf("error receiving stats for intent: %s", err.Error())
	}

	reply := formatUserStatisticsReply(intentID, statsReply.UserStatistics)
	if reply == "" {
		return nil, nil
	}

	return contract.StringResponse(reply), nil
}

func parseIntentID(request contract.Request) string {
	words := strings.Fields(request.Content)

	return words[1]
}

func formatUserStatisticsReply(intentID string, userStatistics []*processingpb.UserStatistic) string {
	if len(userStatistics) < 1 {
		return "No stats are available for a command with id: `" + intentID + "`"
	}

	buffer := bytes.Buffer{}

	for _, userStat := range userStatistics {
		buffer.WriteString(fmt.Sprintf("%s: %d\n", userStat.UserId, userStat.Count))
	}

	return buffer.String()
}
