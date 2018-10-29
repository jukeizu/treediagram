package processor

import (
	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	"github.com/jukeizu/treediagram/processor/command"
	nats "github.com/nats-io/go-nats"
)

type Processor struct {
	logger   log.Logger
	received command.Received
	matched  command.Matched
}

func New(logger log.Logger, queue *nats.EncodedConn, registrationClient registration.RegistrationClient) (Processor, error) {
	logger = log.With(logger, "service", "processor")
	p := Processor{logger: logger}

	r, err := command.NewCommandReceivedProcessor(logger, queue, registrationClient)
	if err != nil {
		return p, err
	}

	p.received = r

	m, err := command.NewCommandMatchedProcessor(logger, queue)
	if err != nil {
		return p, err
	}

	p.matched = m

	return p, nil
}

func (p Processor) Stop() {
	p.received.Stop()
	p.matched.Stop()
}
