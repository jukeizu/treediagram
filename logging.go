package main

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

func (s *loggingService) Request(treediagramRequest TreediagramRequest) (result TreediagramResponse, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Request",
			"request", treediagramRequest,
			"result", result,
			"took", time.Since(begin),
		)

	}(time.Now())

	result, err = s.Service.Request(treediagramRequest)

	return
}
