package intent

import (
	"context"
	"time"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service pb.IntentRegistryServer
}

func NewLoggingService(logger zerolog.Logger, s pb.IntentRegistryServer) pb.IntentRegistryServer {
	logger = logger.With().Str("service", "intent").Logger()
	return &loggingService{logger, s}
}

func (s loggingService) AddIntent(ctx context.Context, req *pb.AddIntentRequest) (reply *pb.AddIntentReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "addIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.AddIntent(ctx, req)

	return
}

func (s loggingService) DisableIntent(ctx context.Context, req *pb.DisableIntentRequest) (reply *pb.DisableIntentReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "disableIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.DisableIntent(ctx, req)

	return
}

func (s loggingService) QueryIntents(ctx context.Context, req *pb.QueryIntentsRequest) (reply *pb.QueryIntentsReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "queryIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("")
	}(time.Now())

	reply, err = s.Service.QueryIntents(ctx, req)

	return
}
