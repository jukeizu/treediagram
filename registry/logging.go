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

func (s loggingService) AddCommand(ctx context.Context, req *pb.AddCommandRequest) (reply *pb.AddCommandReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddCommand",
			"command", *req.Command,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.AddCommand(ctx, req)

	return
}

func (s loggingService) DisableCommand(ctx context.Context, req *pb.DisableCommandRequest) (reply *pb.DisableCommandReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "DisableCommand",
			"id", req.Id,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.DisableCommand(ctx, req)

	return
}

func (s loggingService) QueryCommands(ctx context.Context, req *pb.QueryCommandsRequest) (reply *pb.QueryCommandsReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "QueryCommands",
			"query", *req,
			"took", time.Since(begin),
			"error", err,
		)

	}(time.Now())

	reply, err = s.Service.QueryCommands(ctx, req)

	return
}
