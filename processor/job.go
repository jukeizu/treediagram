package processor

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
)

type Job struct {
	SchedulingJob schedulingpb.Job
}

func (j Job) Execute() (*processingpb.Response, error) {
	reply := &processingpb.Response{}

	if len(j.SchedulingJob.Endpoint) > 0 {
		client := Client{}

		r, err := client.Do(j.SchedulingJob, j.SchedulingJob.Endpoint)
		if err != nil {
			return reply, err
		}

		reply.Messages = r.Messages
	}

	if len(j.SchedulingJob.Content) > 0 {
		m := &processingpb.Message{
			Content: j.SchedulingJob.Content,
		}
		reply.Messages = append(reply.Messages, m)
	}

	return reply, nil
}
