package publishing

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) SendMessage(request SendMessageRequest) (result Response, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendMessage",
			"request.ChannelId", request.ChannelId,
			"request.User", request.User,
			"request.PrivateMessage", request.PrivateMessage,
			"request.IsRedirect", request.IsRedirect,
			"request.CorrelationId", request.CorrelationId,
			"result.Id", result.Id,
			"took", time.Since(begin),
		)

	}(time.Now())

	result, err = s.Service.SendMessage(request)

	return
}
