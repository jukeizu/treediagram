package receiving

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/receiving"
)

type loggingService struct {
	logger  log.Logger
	Service pb.ReceivingServer
}

func NewLoggingService(logger log.Logger, s pb.ReceivingServer) pb.ReceivingServer {
	logger = log.With(logger, "service", "receiving")
	return &loggingService{logger, s}
}

func (s *loggingService) Request(ctx context.Context, req *pb.TreediagramRequest) (reply *pb.TreediagramReply, err error) {
	defer func(begin time.Time) {
		logRequest := req

		logRequest.Content = ""

		s.logger.Log(
			"method", "Request",
			"request", *logRequest,
			"reply", *reply,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Request(ctx, req)

	return
}
