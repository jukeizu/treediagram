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
	reply := &processingpb.Response{}

	if j.SchedulingJob.Endpoint != "" {
		client := Client{}

		r, err := client.Do(j.SchedulingJob, j.SchedulingJob.Endpoint)
		if err != nil {
			return reply, err
		}

		reply.Messages = r.Messages
	}

	if j.SchedulingJob.Content != "" {
		m := &processingpb.Message{
			Content: j.SchedulingJob.Content,
		}
		reply.Messages = append(reply.Messages, m)
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
