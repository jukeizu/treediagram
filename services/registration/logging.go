package registration

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Add(command Command) (Command, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Add",
			"command", command,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.Add(command)
}

func (s *loggingService) Disable(id string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Disable",
			"id", id,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.Disable(id)
}

func (s *loggingService) Query(query CommandQuery) (CommandQueryResult, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Query",
			"query", query,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.Query(query)
}
