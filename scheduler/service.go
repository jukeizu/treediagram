package scheduler

import (
	"context"
	"strconv"
	"time"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

var (
	JobsSubject          = "jobs"
	SchedulerTickSubject = "scheduler.tick"
)

type service struct {
	logger     zerolog.Logger
	JobStorage JobStorage
	Queue      *nats.EncodedConn
}

func NewService(logger zerolog.Logger, storage JobStorage, queue *nats.EncodedConn) pb.SchedulingServer {
	s := &service{logger, storage, queue}

	s.Queue.Subscribe(SchedulerTickSubject, func(req *pb.RunJobsRequest) {
		s.Run(context.Background(), req)
	})

	return s
}

func (s service) Create(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobReply, error) {
	job := &pb.Job{
		Id:          xid.New().String(),
		Type:        req.Type,
		Content:     req.Content,
		User:        req.User,
		Destination: req.Destination,
		Schedule:    req.Schedule,
		Enabled:     true,
	}

	err := s.JobStorage.Create(job)

	return &pb.CreateJobReply{Job: job}, err
}

func (s service) Jobs(ctx context.Context, req *pb.JobsRequest) (*pb.JobsReply, error) {
	if req.Time == 0 {
		jobs, err := s.JobStorage.Jobs(nil)

		return &pb.JobsReply{Jobs: jobs}, err
	}

	t := time.Unix(req.Time, 0).UTC()

	schedule := &pb.Schedule{
		Minute:     strconv.Itoa(t.Minute()),
		Hour:       strconv.Itoa(t.Hour()),
		DayOfMonth: strconv.Itoa(t.Day()),
		Month:      t.Month().String(),
		DayOfWeek:  t.Weekday().String(),
	}

	jobs, err := s.JobStorage.Jobs(schedule)

	return &pb.JobsReply{Jobs: jobs}, err
}

func (s service) Run(ctx context.Context, req *pb.RunJobsRequest) (*pb.RunJobsReply, error) {
	reply, err := s.Jobs(ctx, &pb.JobsRequest{Time: req.Time})
	if err != nil {
		return nil, err
	}

	for _, job := range reply.Jobs {
		err := s.Queue.Publish(JobsSubject, job)
		if err != nil {
			s.logger.Error().Err(err).Caller().
				Str("job.id", job.GetId()).
				Msg("error publishing job")
		}
	}

	return &pb.RunJobsReply{Jobs: reply.Jobs}, nil
}

func (s service) Disable(ctx context.Context, req *pb.DisableJobRequest) (*pb.DisableJobReply, error) {
	err := s.JobStorage.Disable(req.Id)

	return &pb.DisableJobReply{Id: req.Id, Enabled: false}, err
}
