package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/user"
	"github.com/jukeizu/treediagram/services/user"
	"github.com/shawntoffel/services-core/command"
	configreader "github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	Port           int
	UserStorageUrl string
}

func main() {
	logger := logging.GetLogger("services.user", os.Stdout)

	config := Config{}
	err := configreader.ReadConfig(serviceArgs.ConfigFile, &config)
	if err != nil {
		panic(err)
	}

	storage, err := user.NewUserStorage(config.UserStorageUrl)
	if err != nil {
		panic(err)
	}

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
			logger.Log("error", err.Error())
		}

		s := grpc.NewServer()

		pb.RegisterUserServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")
		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
