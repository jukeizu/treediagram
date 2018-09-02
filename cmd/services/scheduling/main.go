package main

import (
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
	mdb "github.com/shawntoffel/GoMongoDb"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	GrpcPort    int
	HttpPort    int
	NatsServers string
	JobStorage  mdb.DbConfig
}

func main() {
	logger := logging.GetLogger("services.scheduling", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	storage, err := scheduling.NewJobStorage(c.JobStorage)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	nc, err := nats.Connect(c.NatsServers)
	if err != nil {
		panic(err)
	}
	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		panic(err)
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
		port := fmt.Sprintf(":%d", c.GrpcPort)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Log("error", err.Error())
		}

		s := grpc.NewServer()
		pb.RegisterSchedulingServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	go func() {
		port := fmt.Sprintf(":%d", c.HttpPort)

		httpBinding := scheduling.NewHttpBinding(logger, service)

		http.Handle("/", httpBinding.NewServeMux())

		logger.Log("transport", "http", "address", port, "msg", "listening")

		errChannel <- http.ListenAndServe(port, nil)
	}()

	logger.Log("stopped", <-errChannel)
}
