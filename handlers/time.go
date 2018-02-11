package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jukeizu/treediagram/command"
	"github.com/shawntoffel/rabbitmq"
)

type Time interface {
	Handler
}

type time struct {
	Receiver rabbitmq.Client
	Command  command.Command
}

func NewTimeHandler(rabbitmqConfig rabbitmq.Config) (Time, error) {
	client, err := rabbitmq.NewClient(rabbitmqConfig)

	command := command.NewTimeCommand()

	return &time{client, command}, err
}

func (t *time) Start() error {

	errors := t.Receiver.Listen(t.handle)

	go func() {
		for err := range errors {
			fmt.Println(err.Error())

		}
	}()

	forever := make(chan bool)
	<-forever

	return nil
}

func (t *time) handle(message []byte) error {
	request := command.Request{}

	err := json.Unmarshal(message, &request)

	if err != nil {
		return err
	}

	if !t.Command.IsCommand(request) {
		return nil
	}

	response := t.Command.Handle(request)

	fmt.Println(response)

	return nil
}
