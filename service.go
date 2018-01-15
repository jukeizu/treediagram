package main

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

func NewService(rabbitmqConfig rabbitmq.Config) Service {
	client := rabbitmq.NewClient(rabbitmqConfig)

	return &service{client}
}

func (s *service) Request(treediagramRequest TreediagramRequest) (TreediagramResponse, error) {

	treediagramResponse := TreediagramResponse{}
	treediagramResponse.Id = xid.New().String()

	marshalled, err := json.Marshal(treediagramRequest)

	if err != nil {
		return treediagramResponse, err
	}

	err = s.Queue.Publish(marshalled)

	if err != nil {
		return treediagramResponse, err
	}

	treediagramResponse.Status = "queued"

	return treediagramResponse, nil
}
