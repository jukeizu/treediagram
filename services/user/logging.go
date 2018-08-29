package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/user"
)

type loggingService struct {
	logger  log.Logger
	Service pb.UserServer
}

func NewLoggingService(logger log.Logger, s pb.UserServer) pb.UserServer {
	return &loggingService{logger, s}
}

func (s *loggingService) Preference(ctx context.Context, req *pb.PreferenceRequest) (reply *pb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Preference",
			"request", *req,
			"reply", *reply,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Preference(ctx, req)

	return
}

func (s *loggingService) SetServer(ctx context.Context, req *pb.SetServerRequest) (reply *pb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SetServer",
			"request", *req,
			"reply", *reply,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.SetServer(ctx, req)

	return
}
