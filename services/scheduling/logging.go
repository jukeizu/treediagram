package scheduling

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/scheduling"
)

type loggingService struct {
	logger  log.Logger
	Service pb.SchedulingServer
}

func NewLoggingService() pb.SchedulingServer {
	return &service{}
}

func (s loggingService) Create(ctx context.Context, req *pb.CreateJobRequest) (reply *pb.CreateJobReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Create",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Create(ctx, req)

	return
}

func (s loggingService) Jobs(ctx context.Context, req *pb.JobsRequest) (reply *pb.JobsReply, err error) {
	reply, err = s.Service.Jobs(ctx, req)

	return
}

func (s loggingService) Disable(ctx context.Context, req *pb.DisableJobRequest) (reply *pb.DisableJobReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Disable",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Disable(ctx, req)

	return
}
