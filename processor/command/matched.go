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
	sub    *nats.Subscription
}

type matchedCommand struct {
	Request processing.TreediagramRequest `json:"request"`
	Command registration.Command          `json:"command"`
}

func NewCommandMatchedProcessor(logger log.Logger, queue *nats.EncodedConn) (Matched, error) {
	m := Matched{logger: logger, queue: queue}

	sub, err := m.queue.QueueSubscribe(CommandMatchedSubject, ProcessorQueueGroup, m.process)
	if err != nil {
		return m, err
	}

	m.sub = sub

	return m, nil
}

func (m Matched) Stop() {
	m.sub.Unsubscribe()
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
