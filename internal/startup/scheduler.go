package startup

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/pkg/scheduling"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

func StartScheduler(logger log.Logger, config Config) error {
	logger = log.With(logger, "component", "scheduler")

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		return err
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return err
	}

	defer conn.Close()

	scheduler := scheduling.NewScheduler(logger, conn)
	scheduler.Start()
	defer scheduler.Stop()

	errs := make(chan error, 2)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return <-errs
}
