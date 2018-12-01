package user

import (
	"context"
	"time"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service pb.UserServer
}

func NewLoggingService(logger zerolog.Logger, s pb.UserServer) pb.UserServer {
	logger = logger.With().Str("service", "user").Logger()
	return &loggingService{logger, s}
}

func (s *loggingService) Preference(ctx context.Context, req *pb.PreferenceRequest) (reply *pb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "preference").
			Str("took", time.Since(begin).String()).
			Str("userId", req.GetUserId()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.Preference(ctx, req)

	return
}

func (s *loggingService) SetServer(ctx context.Context, req *pb.SetServerRequest) (reply *pb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "setServer").
			Str("took", time.Since(begin).String()).
			Str("userId", req.GetUserId()).
			Str("serverId", req.GetServerId()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.SetServer(ctx, req)

	return
}
