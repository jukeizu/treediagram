package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/scheduling"
	"github.com/jukeizu/treediagram/services/scheduling"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

const (
	DefaultGrpcPort = 50054
	DefaultHttpPort = 10002
)

type Config struct {
	GrpcPort      int
	HttpPort      int
	NatsServers   string
	JobStorageUrl string
}

func parseConfig() Config {
	c := Config{}

	flag.IntVar(&c.GrpcPort, "p", DefaultGrpcPort, "port")
	flag.IntVar(&c.HttpPort, "http-port", DefaultHttpPort, "http-port")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.JobStorageUrl, "db", "localhost", "Database connection url")

	flag.Parse()

	return c
}

func main() {
	logger := logging.GetLogger("services.scheduling", os.Stdout)

	config := parseConfig()

	storage, err := scheduling.NewJobStorage(config.JobStorageUrl)
	if err != nil {
		logger.Log("db error", err)
		os.Exit(1)
	}

	defer storage.Close()

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	service := scheduling.NewService(logger, storage, conn)
	service = scheduling.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)

	}()

	go func() {
		port := fmt.Sprintf(":%d", config.GrpcPort)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			errChannel <- err
		}

		s := grpc.NewServer()
		pb.RegisterSchedulingServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	go func() {
		port := fmt.Sprintf(":%d", config.HttpPort)

		httpBinding := scheduling.NewHttpBinding(logger, service)

		logger.Log("transport", "http", "address", port, "msg", "listening")

		errChannel <- http.ListenAndServe(port, httpBinding.NewServeMux())
	}()

	logger.Log("stopped", <-errChannel)
}
