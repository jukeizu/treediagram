package scheduling

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/scheduling"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/xid"
)

var (
	JobsSubject          = "treediagram.jobs"
	SchedulerTickSubject = "treediagram.scheduler.tick"
)

type service struct {
	logger     log.Logger
	JobStorage JobStorage
	Queue      *nats.EncodedConn
}

func NewService(logger log.Logger, storage JobStorage, queue *nats.EncodedConn) pb.SchedulingServer {
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
	schedule := &pb.Schedule{}

	if req.Time > 0 {
		t := time.Unix(req.Time, 0).UTC()

		schedule.Minute = strconv.Itoa(t.Minute())
		schedule.Hour = strconv.Itoa(t.Hour())
		schedule.DayOfMonth = strconv.Itoa(t.Day())
		schedule.Month = t.Month().String()
		schedule.DayOfWeek = t.Weekday().String()
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
			s.logger.Log(
				"method", "Run",
				"job", *job,
				"error", "Error publishing job: "+err.Error(),
			)
		}
	}

	return &pb.RunJobsReply{Jobs: reply.Jobs}, nil
}

func (s service) Disable(ctx context.Context, req *pb.DisableJobRequest) (*pb.DisableJobReply, error) {
	err := s.JobStorage.Disable(req.Id)

	return &pb.DisableJobReply{Id: req.Id, Enabled: false}, err
}
