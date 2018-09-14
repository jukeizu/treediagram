package publishing

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/treediagram/api/publishing"
)

type loggingService struct {
	logger  log.Logger
	Service pb.PublishingServer
}

func NewLoggingService(logger log.Logger, s pb.PublishingServer) pb.PublishingServer {
	logger = log.With(logger, "service", "publishing")
	return &loggingService{logger, s}
}

func (s loggingService) PublishMessage(ctx context.Context, req *pb.PublishMessageRequest) (reply *pb.PublishMessageReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendMessage",
			"request.ChannelId", req.Message.ChannelId,
			"request.User", *req.Message.User,
			"request.PrivateMessage", req.Message.PrivateMessage,
			"request.IsRedirect", req.Message.IsRedirect,
			"request.CorrelationId", req.Message.CorrelationId,
			"reply.Id", reply.Id,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.PublishMessage(ctx, req)

	return
}
