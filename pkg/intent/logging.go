package intent

import (
	"context"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service intentpb.IntentRegistryServer
}

func NewLoggingService(logger zerolog.Logger, s intentpb.IntentRegistryServer) intentpb.IntentRegistryServer {
	logger = logger.With().Str("service", "intent").Logger()
	return &loggingService{logger, s}
}

func (s loggingService) AddIntent(ctx context.Context, req *intentpb.AddIntentRequest) (reply *intentpb.AddIntentReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "addIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	reply, err = s.Service.AddIntent(ctx, req)

	return
}

func (s loggingService) DisableIntent(ctx context.Context, req *intentpb.DisableIntentRequest) (reply *intentpb.DisableIntentReply, err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "disableIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	reply, err = s.Service.DisableIntent(ctx, req)

	return
}

func (s loggingService) QueryIntents(req *intentpb.QueryIntentsRequest, stream intentpb.IntentRegistry_QueryIntentsServer) (err error) {
	defer func(begin time.Time) {
		l := s.logger.With().
			Str("method", "queryIntent").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Debug().Msg("called")
	}(time.Now())

	err = s.Service.QueryIntents(req, stream)

	return
}
