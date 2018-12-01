package startup

import (
	"github.com/jukeizu/treediagram/scheduler"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/rs/zerolog"
)

type SchedulerRunner struct {
	Logger       zerolog.Logger
	EncodedConn  *nats.EncodedConn
	Subscription *nats.Subscription
	Scheduler    scheduler.Scheduler
	quit         chan struct{}
}

func NewSchedulerRunner(logger zerolog.Logger, config Config) (*SchedulerRunner, error) {
	logger = logger.With().Str("component", "scheduler").Logger()

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, err
	}

	scheduler := scheduler.NewScheduler(logger, conn)

	scheduleRunner := &SchedulerRunner{
		Logger:      logger,
		EncodedConn: conn,
		Scheduler:   scheduler,
		quit:        make(chan struct{}),
	}

	return scheduleRunner, nil
}

func (r *SchedulerRunner) Start() error {
	r.Logger.Info().Msg("starting")

	r.Scheduler.Start()

	<-r.quit

	return nil
}

func (r *SchedulerRunner) Stop() {
	r.Logger.Info().Msg("stopping")

	close(r.quit)
	r.Scheduler.Stop()
	r.EncodedConn.Close()
}
