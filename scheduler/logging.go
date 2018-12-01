package scheduler

import (
	"context"
	"time"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service pb.SchedulingServer
}

func NewLoggingService(logger zerolog.Logger, s pb.SchedulingServer) pb.SchedulingServer {
	logger = logger.With().Str("service", "scheduling").Logger()
	return &loggingService{logger, s}
}

func (s loggingService) Create(ctx context.Context, req *pb.CreateJobRequest) (reply *pb.CreateJobReply, err error) {
	defer func(begin time.Time) {
		s.logger.Info().
			Str("method", "create").
			Err(err).
			Str("took", time.Since(begin).String()).
			Msg("")
	}(time.Now())

	reply, err = s.Service.Create(ctx, req)

	return
}

func (s loggingService) Jobs(ctx context.Context, req *pb.JobsRequest) (reply *pb.JobsReply, err error) {
	defer func(begin time.Time) {
		s.logger.Debug().
			Str("method", "jobs").
			Str("took", time.Since(begin).String()).
			Msg("")
	}(time.Now())

	reply, err = s.Service.Jobs(ctx, req)

	return
}

func (s loggingService) Run(ctx context.Context, req *pb.RunJobsRequest) (reply *pb.RunJobsReply, err error) {
	defer func(begin time.Time) {
		s.logger.Info().
			Str("method", "run").
			Err(err).
			Str("took", time.Since(begin).String()).
			Msg("")
	}(time.Now())

	reply, err = s.Service.Run(ctx, req)

	return
}

func (s loggingService) Disable(ctx context.Context, req *pb.DisableJobRequest) (reply *pb.DisableJobReply, err error) {
	defer func(begin time.Time) {
		s.logger.Info().
			Str("method", "disable").
			Err(err).
			Str("took", time.Since(begin).String()).
			Msg("")
	}(time.Now())

	reply, err = s.Service.Disable(ctx, req)

	return
}
