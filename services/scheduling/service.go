package scheduling

import (
	"context"
	"strconv"
	"time"

	pb "github.com/jukeizu/treediagram/api/scheduling"
	"github.com/rs/xid"
)

type service struct {
	JobStorage JobStorage
}

func NewService() pb.SchedulingServer {
	return &service{}
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

func (s service) Disable(ctx context.Context, req *pb.DisableJobRequest) (*pb.DisableJobReply, error) {
	err := s.JobStorage.Disable(req.Id)

	return &pb.DisableJobReply{Id: req.Id, Enabled: false}, err
}
