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

func NewService(rabbitmqConfig rabbitmq.Config) (Service, error) {
	client, err := rabbitmq.NewClient(rabbitmqConfig)

	return &service{client}, err
}

func (s *service) Request(treediagramRequest TreediagramRequest) (TreediagramResponse, error) {

	id := xid.New()
	treediagramResponse := TreediagramResponse{}
	treediagramResponse.Id = id.String()

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
