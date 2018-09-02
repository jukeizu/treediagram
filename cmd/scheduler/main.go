package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jukeizu/treediagram/scheduler"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	NatsServers string
}

func main() {
	logger := logging.GetLogger("scheduler", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	nc, err := nats.Connect(c.NatsServers)
	if err != nil {
		panic(err)
	}
	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		panic(err)
	}

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
