package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jukeizu/treediagram/scheduler"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/shawntoffel/services-core/logging"
)

type Config struct {
	NatsServers string
}

func parseConfig() Config {
	c := Config{}

	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.Parse()

	return c
}

func main() {
	logger := logging.GetLogger("scheduler", os.Stdout)

	c := parseConfig()

	nc, err := nats.Connect(c.NatsServers)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer conn.Close()

	scheduler := scheduler.NewScheduler(logger, conn)
	scheduler.Start()
	defer scheduler.Stop()

	errs := make(chan error, 2)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("stopped", <-errs)
}
