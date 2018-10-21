package scheduler

import (
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	nats "github.com/nats-io/go-nats"
	"github.com/robfig/cron"
)

type Scheduler interface {
	Start()
	Stop()
}

type scheduler struct {
	Logger log.Logger
	Queue  *nats.EncodedConn
	Cron   *cron.Cron
}

func NewScheduler(logger log.Logger, queue *nats.EncodedConn) Scheduler {
	s := scheduler{
		Logger: logger,
		Queue:  queue,
		Cron:   cron.New(),
	}

	s.Cron.AddFunc("0 * * * * *", s.run)

	return &s
}

func (s *scheduler) Start() {
	s.Cron.Start()

	s.Logger.Log("msg", "started")
}

func (s *scheduler) Stop() {
	s.Cron.Stop()

	s.Logger.Log("msg", "stopped")
}

func (s *scheduler) run() {
	request := &pb.RunJobsRequest{
		Time: time.Now().Unix(),
	}

	err := s.Queue.Publish(SchedulerTickSubject, request)
	if err != nil {
		s.Logger.Log("error", err.Error())
	}
}
