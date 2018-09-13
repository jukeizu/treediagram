package receiving

import (
	"context"
	"encoding/json"

	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/rs/xid"
	"github.com/shawntoffel/rabbitmq"
)

type service struct {
	Queue rabbitmq.Client
}

func NewService(rabbitmqUrl string) (pb.ReceivingServer, error) {
	rabbitmqConfig := rabbitmq.Config{
		Durable:      true,
		QueueName:    "treediagram",
		Exchange:     "treediagram-exchange",
		ExchangeType: "fanout",
		Url:          rabbitmqUrl,
	}

	client, err := rabbitmq.NewPublisher(rabbitmqConfig)

	return &service{client}, err
}

func (s service) Request(ctx context.Context, req *pb.TreediagramRequest) (*pb.TreediagramReply, error) {
	id := xid.New().String()

	treediagramReply := &pb.TreediagramReply{
		Id: id,
	}

	marshalled, err := json.Marshal(req)

	if err != nil {
		return treediagramReply, err
	}

	err = s.Queue.Publish(marshalled)

	if err != nil {
		return treediagramReply, err
	}

	return treediagramReply, nil
}
