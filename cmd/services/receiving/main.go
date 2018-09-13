package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/jukeizu/treediagram/services/receiving"
	"github.com/shawntoffel/rabbitmq"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

type Config struct {
	Port           int
	RabbitMqConfig rabbitmq.Config
}

const (
	DefaultPort                    = 50052
	RabbitMqUrlEnvironmentVariable = "TREEDIAGRAM_RABBITMQ_URL"
)

func parseConfig() Config {
	c := Config{}

	c.RabbitMqConfig = rabbitmq.Config{
		Durable:      true,
		QueueName:    "treediagram",
		Exchange:     "treediagram-exchange",
		ExchangeType: "fanout",
	}

	flag.IntVar(&c.Port, "p", DefaultPort, "port")
	flag.StringVar(&c.RabbitMqConfig.Url, "rmq", rabbitmq.DefaultUrl, "RabbitMQ url. This can also be specified via the "+RabbitMqUrlEnvironmentVariable+" environment variable.")

	flag.Parse()

	if c.RabbitMqConfig.Url == "" {
		c.RabbitMqConfig.Url = os.Getenv(RabbitMqUrlEnvironmentVariable)
	}

	return c
}

func main() {
	logger := logging.GetLogger("services.receiving", os.Stdout)

	treediagramConfig := parseConfig()

	service, err := receiving.NewService(treediagramConfig.RabbitMqConfig)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
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
