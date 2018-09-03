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
	return &loggingService{logger, s}
}

func (s loggingService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (reply *pb.SendMessageReply, err error) {
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

	reply, err = s.Service.SendMessage(ctx, req)

	return
}
