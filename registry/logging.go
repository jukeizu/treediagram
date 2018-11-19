package registry

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/registration"
)

type loggingService struct {
	logger  log.Logger
	Service pb.RegistrationServer
}

func NewLoggingService(logger log.Logger, s pb.RegistrationServer) pb.RegistrationServer {
	logger = log.With(logger, "service", "registration")
	return &loggingService{logger, s}
}

func (s loggingService) AddIntent(ctx context.Context, req *pb.AddIntentRequest) (reply *pb.AddIntentReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddIntent",
			"intent", *req.Intent,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.AddIntent(ctx, req)

	return
}

func (s loggingService) DisableIntent(ctx context.Context, req *pb.DisableIntentRequest) (reply *pb.DisableIntentReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "DisableIntent",
			"id", req.Id,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.DisableIntent(ctx, req)

	return
}

func (s loggingService) QueryIntents(ctx context.Context, req *pb.QueryIntentsRequest) (reply *pb.QueryIntentsReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "QueryIntents",
			"query", *req,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.QueryIntents(ctx, req)

	return
}
