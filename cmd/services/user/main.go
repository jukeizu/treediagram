package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/user"
	"github.com/jukeizu/treediagram/services/user"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

const (
	DefaultPort = 50055
)

type Config struct {
	Port           int
	UserStorageUrl string
}

func parseConfig() Config {
	c := Config{}

	flag.IntVar(&c.Port, "p", DefaultPort, "port")
	flag.StringVar(&c.UserStorageUrl, "db", "localhost", "Database connection url")

	flag.Parse()

	return c
}

func main() {
	logger := logging.GetLogger("services.user", os.Stdout)

	config := parseConfig()

	storage, err := user.NewUserStorage(config.UserStorageUrl)
	if err != nil {
		logger.Log("db error", err)
		os.Exit(1)
	}

	defer storage.Close()

	service := user.NewService(storage)
	service = user.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		port := fmt.Sprintf(":%d", config.Port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			errChannel <- err
		}

		s := grpc.NewServer()

		pb.RegisterUserServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")
		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
