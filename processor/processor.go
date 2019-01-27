package processor

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
)

const (
	ProcessorQueueGroup           = "processor"
	MessageRequestReceivedSubject = "messagerequest.received"
	JobReceivedSubject            = "jobs"
	ReplyReceivedSubject          = "processor.reply.received"
)

type Executable interface {
	ShouldExecute() (bool, error)
	Execute() (*processingpb.Response, error)
	ProcessingRequest() *processingpb.ProcessingRequest
}

type Processor struct {
	logger     zerolog.Logger
	queue      *nats.EncodedConn
	registry   intentpb.IntentRegistryClient
	userClient userpb.UserClient
	repository Repository
	waitGroup  *sync.WaitGroup
}

func New(logger zerolog.Logger, queue *nats.EncodedConn, registry intentpb.IntentRegistryClient, userClient userpb.UserClient, repository Repository) Processor {
	p := Processor{
		logger:     logger.With().Str("component", "processor").Logger(),
		queue:      queue,
		registry:   registry,
		userClient: userClient,
		repository: repository,
		waitGroup:  &sync.WaitGroup{},
	}

	return p
}

func (p Processor) Start() error {
	_, err := p.queue.QueueSubscribe(MessageRequestReceivedSubject, ProcessorQueueGroup, p.processMessageRequest)
	if err != nil {
		return err
	}

	_, err = p.queue.QueueSubscribe(JobReceivedSubject, ProcessorQueueGroup, p.processJob)
	if err != nil {
		return err
	}

	p.logger.Info().Msg("started")

	return nil
}

func (p Processor) Stop() {
	p.logger.Info().Msg("stopping")
	p.waitGroup.Wait()
}

func (p Processor) processMessageRequest(request *processingpb.MessageRequest) {
	p.logger.Debug().Msgf("request received %+v", request)

	serverId, err := p.findServerId(request)
	if err != nil {
		p.logger.Error().Err(err).Caller().
			Str("userId", request.Author.Id).
			Msg("could not determine serverId for user")
		return
	}

	if serverId == "" {
		server, err := p.setDefaultServer(request)
		if err != nil {
			p.logger.Error().Err(err).Caller().
				Str("userId", request.Author.Id).
				Msg("could not set default server")
			return
		}

		serverId = server.Id
	}

	query := &intentpb.QueryIntentsRequest{ServerId: serverId}

	stream, err := p.registry.QueryIntents(context.Background(), query)
	if err != nil {
		p.logger.Error().Err(err).Caller().Msg("error getting QueryIntents stream")
		return
	}

	for {
		intent, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("error receiving intent from stream")
			break
		}

		command := Command{
			Request: *request,
			Intent:  *intent,
		}

		p.process(command)
	}
}

func (p Processor) processJob(schedulingJob *schedulingpb.Job) {
	p.logger.Debug().Msgf("job received %+v", schedulingJob)

	job := Job{
		SchedulingJob: *schedulingJob,
	}

	p.process(job)
}

func (p Processor) process(executable Executable) {
	p.waitGroup.Add(1)
	go func(executable Executable) {
		defer p.waitGroup.Done()
		p.logger.Debug().Msg("starting processing for executable")

		shouldExecute, err := executable.ShouldExecute()
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("could not determine if should execute")
			return
		}

		if !shouldExecute {
			p.logger.Debug().Msg("executable should not execute")
			return
		}

		processingRequest := executable.ProcessingRequest()

		err = p.repository.SaveProcessingRequest(processingRequest)
		if err != nil {
			p.logger.Error().Err(err).Caller().Msg("error saving ProcessingRequest")
			return
		}

		p.logger.Debug().
			Str("processingRequestId", processingRequest.Id).
			Str("processingRequestType", processingRequest.Type).
			Msg("executing")

		response, err := executable.Execute()
		if err != nil {
			p.createErrorEvent(processingRequest.Id, err)
			return
		}

		p.logger.Debug().
			Str("processingRequestId", processingRequest.Id).
			Str("processingRequestType", processingRequest.Type).
			Msg("saving responses from execute")

		for _, message := range response.Messages {
			p.saveResponseMessage(processingRequest, *message)
		}

		p.logger.Debug().
			Str("processingRequestId", processingRequest.Id).
			Str("processingRequestType", processingRequest.Type).
			Msg("finished processing executable")
	}(executable)
}

func (p Processor) findServerId(request *processingpb.MessageRequest) (string, error) {
	if request.ServerId != "" {
		return request.ServerId, nil
	}

	userId := request.Author.Id

	p.logger.Debug().
		Str("userId", userId).
		Msg("looking up server preference")

	preferenceReply, err := p.userClient.Preference(context.Background(), &userpb.PreferenceRequest{UserId: userId})
	if err != nil {
		return "", err
	}

	preference := preferenceReply.Preference
	if preference == nil {
		p.logger.Debug().
			Str("userId", userId).
			Msg("user does not have a preferred server")

		return "", nil
	}

	p.logger.Debug().
		Str("userId", userId).
		Str("preferredServerId", preference.ServerId).
		Msg("found server preference for user")

	return preference.ServerId, nil
}

func (p Processor) setDefaultServer(request *processingpb.MessageRequest) (*processingpb.Server, error) {
	if len(request.Servers) < 1 {
		return nil, errors.New("no servers are available to select default")
	}

	userId := request.Author.Id
	server := request.Servers[0]

	setServerRequest := userpb.SetServerRequest{
		UserId:   userId,
		ServerId: server.Id,
	}

	_, err := p.userClient.SetServer(context.Background(), &setServerRequest)
	if err != nil {
		return nil, err
	}

	p.logger.Debug().
		Str("userId", userId).
		Str("serverId", server.Id).
		Msg("user server preference has been set to the first available server")

	return server, nil
}

func (p Processor) saveResponseMessage(processingRequest *processingpb.ProcessingRequest, message processingpb.Message) {
	messageReply := processingpb.MessageReply{
		ProcessingRequestId: processingRequest.Id,
		ChannelId:           processingRequest.ChannelId,
		UserId:              processingRequest.UserId,
		IsPrivateMessage:    message.IsPrivateMessage,
		IsRedirect:          message.IsRedirect,
		Content:             message.Content,
	}

	err := p.repository.SaveMessageReply(&messageReply)
	if err != nil {
		p.createErrorEvent(processingRequest.Id, err)
		return
	}

	messageReplyReceived := processingpb.MessageReplyReceived{Id: messageReply.Id}

	err = p.queue.Publish(ReplyReceivedSubject+"."+processingRequest.Source, messageReplyReceived)
	if err != nil {
		p.createErrorEvent(processingRequest.Id, err)
		return
	}
}

func (p Processor) createProcessingEvent(processingRequestId string, eventType string, eventDescription string) {
	e := processingpb.ProcessingEvent{
		ProcessingRequestId: processingRequestId,
		Type:                eventType,
		Description:         eventDescription,
	}

	err := p.repository.SaveProcessingEvent(&e)
	if err != nil {
		p.logger.Error().Err(err).Caller().Msg("error saving processing event")
		return
	}

	p.logger.Debug().
		Str("processingRequestId", processingRequestId).
		Str("type", eventType).
		Str("description", eventDescription).
		Msg("saved processing event")
}

func (p Processor) createErrorEvent(processingRequestId string, processingError error) {
	p.createProcessingEvent(processingRequestId, "error", processingError.Error())
	p.logger.Error().Err(processingError).Caller().
		Str("processingRequestId", processingRequestId).
		Msg("")
}
