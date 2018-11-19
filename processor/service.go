package processor

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
)

type service struct {
	storage Storage
}

func NewService(logger log.Logger, storage Storage) (processing.ProcessingServer, error) {
	return &service{storage: storage}, nil
}

func (s service) Messages(ctx context.Context, req *processing.MessagesRequest) (*processing.Reply, error) {
	return s.storage.Reply(req.Id)
}
