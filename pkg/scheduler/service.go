package scheduler

import (
	"context"
	"strconv"
	"time"

	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

var (
	JobsSubject          = "jobs"
	SchedulerTickSubject = "scheduler.tick"
	SchedulerQueueGroup  = "scheduler"
)

type service struct {
	logger     zerolog.Logger
	Repository Repository
	Queue      *nats.EncodedConn
}

func NewService(logger zerolog.Logger, repository Repository, queue *nats.EncodedConn) (schedulingpb.SchedulingServer, error) {
	s := &service{logger, repository, queue}

	_, err := s.Queue.QueueSubscribe(SchedulerTickSubject, SchedulerQueueGroup, func(req *schedulingpb.RunJobsRequest) {
		s.Run(context.Background(), req)
	})

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s service) Create(ctx context.Context, req *schedulingpb.CreateJobRequest) (*schedulingpb.CreateJobReply, error) {
	job := &schedulingpb.Job{
		UserId:      req.UserId,
		InstanceId:  req.InstanceId,
		Source:      req.Source,
		Content:     req.Content,
		Endpoint:    req.Endpoint,
		Destination: req.Destination,
		Schedule:    req.Schedule,
		Enabled:     true,
	}

	err := s.Repository.Create(job)
	if err != nil {
		return nil, err
	}

	return &schedulingpb.CreateJobReply{Job: job}, nil
}

func (s service) Jobs(ctx context.Context, req *schedulingpb.JobsRequest) (*schedulingpb.JobsReply, error) {
	if req.Time == 0 {
		jobs, err := s.Repository.Jobs(nil)

		return &schedulingpb.JobsReply{Jobs: jobs}, err
	}

	t := time.Unix(req.Time, 0).UTC()

	schedule := &schedulingpb.Schedule{
		Minute:     strconv.Itoa(t.Minute()),
		Hour:       strconv.Itoa(t.Hour()),
		DayOfMonth: strconv.Itoa(t.Day()),
		Year:       strconv.Itoa(t.Year()),
		Month:      t.Month().String(),
		DayOfWeek:  t.Weekday().String(),
	}

	jobs, err := s.Repository.Jobs(schedule)

	return &schedulingpb.JobsReply{Jobs: jobs}, err
}

func (s service) Run(ctx context.Context, req *schedulingpb.RunJobsRequest) (*schedulingpb.RunJobsReply, error) {
	reply, err := s.Jobs(ctx, &schedulingpb.JobsRequest{Time: req.Time})
	if err != nil {
		return nil, err
	}

	for _, job := range reply.Jobs {
		err := s.Queue.Publish(JobsSubject, job)
		if err != nil {
			s.logger.Error().Err(err).Caller().
				Str("jobId", job.GetId()).
				Msg("error publishing job")
		}

		s.logger.Info().
			Str("jobId", job.GetId()).
			Msg("published job run to queue")
	}

	return &schedulingpb.RunJobsReply{Jobs: reply.Jobs}, nil
}

func (s service) Disable(ctx context.Context, req *schedulingpb.DisableJobRequest) (*schedulingpb.DisableJobReply, error) {
	err := s.Repository.Disable(req.Id)

	return &schedulingpb.DisableJobReply{Id: req.Id, Enabled: false}, err
}
