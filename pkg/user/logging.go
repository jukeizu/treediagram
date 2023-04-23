package user

import (
	"context"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service userpb.UserServer
}

func NewLoggingService(logger zerolog.Logger, s userpb.UserServer) userpb.UserServer {
	logger = logger.With().Str("service", "user").Logger()
	return &loggingService{logger, s}
}

func (s *loggingService) Preference(ctx context.Context, req *userpb.PreferenceRequest) (reply *userpb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "Preference").
			Str("took", time.Since(begin).String()).
			Str("userId", req.GetUserId()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Debug().Msg("")
	}(time.Now())

	reply, err = s.Service.Preference(ctx, req)

	return
}

func (s *loggingService) SetServer(ctx context.Context, req *userpb.SetServerRequest) (reply *userpb.PreferenceReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "SetServer").
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
