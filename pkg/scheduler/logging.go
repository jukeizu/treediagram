package scheduler

import (
	"context"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service schedulingpb.SchedulingServer
}

func NewLoggingService(logger zerolog.Logger, s schedulingpb.SchedulingServer) schedulingpb.SchedulingServer {
	logger = logger.With().Str("service", "scheduling").Logger()
	return &loggingService{logger, s}
}

func (s loggingService) Create(ctx context.Context, req *schedulingpb.CreateJobRequest) (reply *schedulingpb.CreateJobReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "create").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.Create(ctx, req)

	return
}

func (s loggingService) Jobs(ctx context.Context, req *schedulingpb.JobsRequest) (reply *schedulingpb.JobsReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "jobs").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.Jobs(ctx, req)

	return
}

func (s loggingService) Run(ctx context.Context, req *schedulingpb.RunJobsRequest) (reply *schedulingpb.RunJobsReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "run").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.Run(ctx, req)

	return
}

func (s loggingService) Disable(ctx context.Context, req *schedulingpb.DisableJobRequest) (reply *schedulingpb.DisableJobReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "disable").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.Disable(ctx, req)

	return
}
