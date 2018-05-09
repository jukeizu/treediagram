package receiving

import (
	"encoding/json"
	"github.com/rs/xid"
	"github.com/shawntoffel/rabbitmq"
)

type Service interface {
	Request(TreediagramRequest) (TreediagramResponse, error)
}

type service struct {
	Queue rabbitmq.Client
}

func NewService(rabbitmqConfig rabbitmq.Config) (Service, error) {
	client, err := rabbitmq.NewPublisher(rabbitmqConfig)

	return &service{client}, err
}

func (s *service) Request(treediagramRequest TreediagramRequest) (TreediagramResponse, error) {
	correlationId := xid.New().String()

	treediagramResponse := TreediagramResponse{
		Id: correlationId,
	}

	marshalled, err := json.Marshal(treediagramRequest)

	if err != nil {
		return treediagramResponse, err
	}

	err = s.Queue.Publish(marshalled)

	if err != nil {
		return treediagramResponse, err
	}

	return treediagramResponse, nil
}
