package command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	nats "github.com/nats-io/go-nats"
)

const (
	CommandMatchedSubject = "commandMatched"
)

type Matched struct {
	logger log.Logger
	queue  *nats.EncodedConn
}

type matchedCommand struct {
	Request processing.TreediagramRequest `json:"request"`
	Command registration.Command          `json:"command"`
}

func NewCommandMatchedProcessor(logger log.Logger, queue *nats.EncodedConn) Matched {
	return Matched{logger: logger, queue: queue}
}

func (m Matched) Subscribe() error {
	if m.queue == nil {
		return nil
	}

	_, err := m.queue.QueueSubscribe(CommandMatchedSubject, ProcessorQueueGroup, m.process)
	if err != nil {
		return err
	}

	return nil
}

func (m Matched) process(e matchedCommand) {
	m.logger.Log("msg", "received matched command", "command", e.Command.Id)

	body, err := json.Marshal(e.Command)
	if err != nil {
		m.logger.Log("error", "could not marshal command: "+err.Error())
	}

	resp, err := http.Post(e.Command.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		m.logger.Log("error", "error sending command request: "+err.Error())
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.logger.Log("error", "error reading response body: "+err.Error())
	}

	m.logger.Log("msg", "received reply: "+string(b))
}
