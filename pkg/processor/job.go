package processor

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	"github.com/rs/zerolog"
)

type Job struct {
	SchedulingJob *schedulingpb.Job
}

func (j Job) ShouldExecute() (bool, error) {
	return true, nil
}

func (j Job) Execute() (*processingpb.Response, error) {
	reply := &processingpb.Response{}

	if j.SchedulingJob.Endpoint != "" {
		client := Client{}

		r, err := client.Do(j.SchedulingJob, j.SchedulingJob.Endpoint)
		if err != nil {
			return reply, err
		}

		if r != "" {
			reply.Messages = append(reply.Messages, r)
		}
	}

	if j.SchedulingJob.Content != "" {
		reply.Messages = append(reply.Messages, j.SchedulingJob.Content)
	}

	return reply, nil
}

func (j Job) ProcessingRequest() *processingpb.ProcessingRequest {
	processingRequest := &processingpb.ProcessingRequest{
		Type:      "job",
		Source:    j.SchedulingJob.Source,
		ChannelId: j.SchedulingJob.Destination,
		UserId:    j.SchedulingJob.UserId,
	}

	return processingRequest
}

func (j Job) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", "job").
		Str("jobId", j.SchedulingJob.Id).
		Str("source", j.SchedulingJob.Source).
		Str("channelId", j.SchedulingJob.Destination).
		Str("userId", j.SchedulingJob.UserId)
}
