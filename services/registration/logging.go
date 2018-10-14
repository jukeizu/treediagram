package registration

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/registration"
)

type loggingService struct {
	logger  log.Logger
	Service pb.RegistrationServer
}

func NewLoggingService(logger log.Logger, s pb.RegistrationServer) pb.RegistrationServer {
	logger = log.With(logger, "service", "registration")
	return &loggingService{logger, s}
}

func (s loggingService) AddCommand(ctx context.Context, req *pb.AddCommandRequest) (*pb.AddCommandReply, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddCommand",
			"command", *req.Command,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.AddCommand(ctx, req)
}

func (s loggingService) DisableCommand(ctx context.Context, req *pb.DisableCommandRequest) (*pb.DisableCommandReply, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "DisableCommand",
			"id", req.Id,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.DisableCommand(ctx, req)
}

func (s loggingService) QueryCommands(ctx context.Context, req *pb.QueryCommandsRequest) (*pb.QueryCommandsReply, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "QueryCommands",
			"query", *req,
			"took", time.Since(begin),
		)

	}(time.Now())

	return s.Service.QueryCommands(ctx, req)
}
