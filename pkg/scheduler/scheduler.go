package scheduler

import (
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	nats "github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Scheduler interface {
	Start()
	Stop()
}

type scheduler struct {
	Logger zerolog.Logger
	Queue  *nats.EncodedConn
	Cron   *cron.Cron
}

func NewScheduler(logger zerolog.Logger, queue *nats.EncodedConn) Scheduler {
	s := scheduler{
		Logger: logger,
		Queue:  queue,
		Cron:   cron.New(),
	}

	s.Cron.AddFunc("* * * * *", s.run)

	return &s
}

func (s *scheduler) Start() {
	s.Cron.Start()

	s.Logger.Info().Msg("started")
}

func (s *scheduler) Stop() {
	s.Cron.Stop()

	s.Logger.Info().Msg("stopped")
}

func (s *scheduler) run() {
	request := &schedulingpb.RunJobsRequest{
		Time: time.Now().Unix(),
	}

	err := s.Queue.Publish(SchedulerTickSubject, request)
	if err != nil {
		s.Logger.Error().Err(err).Caller().Msg("error publishing scheduler tick to queue")
	}
}
