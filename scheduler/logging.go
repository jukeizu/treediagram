package scheduler

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

func NewLoggingService(logger log.Logger, s pb.SchedulingServer) pb.SchedulingServer {
	logger = log.With(logger, "service", "scheduling")
	return &loggingService{logger, s}
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
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Log(
				"method", "Jobs",
				"request", *req,
				"error", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	reply, err = s.Service.Jobs(ctx, req)

	return
}

func (s loggingService) Run(ctx context.Context, req *pb.RunJobsRequest) (reply *pb.RunJobsReply, err error) {
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Log(
				"method", "Run",
				"request", *req,
				"error", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	reply, err = s.Service.Run(ctx, req)

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
