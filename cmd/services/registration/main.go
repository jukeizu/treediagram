package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/registration"
	"github.com/jukeizu/treediagram/services/registration"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

const (
	DefaultGrpcPort = 50053
	DefaultHttpPort = 10001
)

type Config struct {
	HttpPort          int
	GrpcPort          int
	CommandStorageUrl string
}

func parseConfig() Config {
	c := Config{}

	flag.IntVar(&c.GrpcPort, "p", DefaultGrpcPort, "port")
	flag.IntVar(&c.HttpPort, "http-port", DefaultHttpPort, "http-port")
	flag.StringVar(&c.CommandStorageUrl, "db", "localhost", "Database connection url")

	flag.Parse()

	return c
}

func main() {
	logger := logging.GetLogger("services.registration", os.Stdout)

	config := parseConfig()

	storage, err := registration.NewCommandStorage(config.CommandStorageUrl)
	if err != nil {
		logger.Log("db error", err)
		os.Exit(1)
	}

	defer storage.Close()

	service, err := registration.NewService(storage)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	service = registration.NewLoggingService(logger, service)

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

		pb.RegisterRegistrationServer(s, registration.GrpcBinding{Service: service})

		logger.Log("transport", "grpc", "address", port, "msg", "listening")
		errChannel <- s.Serve(listener)
	}()

	go func() {
		port := fmt.Sprintf(":%d", config.HttpPort)

		handler := registration.MakeHandler(service, logger)

		mux := http.NewServeMux()
		mux.Handle("/add", handler)
		mux.Handle("/disable", handler)
		mux.Handle("/query", handler)

		logger.Log("transport", "http", "address", port, "msg", "listening")
		errChannel <- http.ListenAndServe(port, mux)
	}()

	logger.Log("stopped", <-errChannel)
}
