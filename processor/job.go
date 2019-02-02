package processor

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
)

type Job struct {
	SchedulingJob schedulingpb.Job
}

func (j Job) ShouldExecute() (bool, error) {
	return true, nil
}

func (j Job) Execute() (*processingpb.Response, error) {
	reply := &processingpb.Response{
		Messages: []string{j.SchedulingJob.Content},
	}

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
