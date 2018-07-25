package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/jukeizu/treediagram/services/receiving"
	"github.com/shawntoffel/rabbitmq"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

type TreediagramConfig struct {
	config.ServiceConfig
	RabbitMqConfig rabbitmq.Config
}

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

func main() {
	logger := logging.GetLogger("services.receiving", os.Stdout)

	treediagramConfig := TreediagramConfig{}

	err := config.ReadConfig(serviceArgs.ConfigFile, &treediagramConfig)

	if err != nil {
		panic(err)
	}

	service, err := receiving.NewService(treediagramConfig.RabbitMqConfig)

	if err != nil {
		panic(err)
	}

	service = receiving.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)

	}()

	go func() {
		port := fmt.Sprintf(":%d", treediagramConfig.Port)

		listener, err := net.Listen("tcp", port)

		if err != nil {
			errChannel <- err
		}

		s := grpc.NewServer()
		pb.RegisterReceivingServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
