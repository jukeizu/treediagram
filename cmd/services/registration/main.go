package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/registration"
	"github.com/jukeizu/treediagram/services/registration"
	mdb "github.com/shawntoffel/GoMongoDb"
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
	HttpPort       int
	GrpcPort       int
	CommandStorage mdb.DbConfig
}

func main() {
	logger := logging.GetLogger("services.registration", os.Stdout)

	config := Config{}
	err := configreader.ReadConfig(serviceArgs.ConfigFile, &config)
	if err != nil {
		panic(err)
	}

	storage, err := registration.NewCommandStorage(config.CommandStorage)
	if err != nil {
		panic(err)
	}

	service, err := registration.NewService(storage)
	if err != nil {
		panic(err)
	}

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
			logger.Log("error", err.Error())
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

		http.Handle("/", mux)

		logger.Log("transport", "http", "address", port, "msg", "listening")
		errChannel <- http.ListenAndServe(port, nil)
	}()

	logger.Log("stopped", <-errChannel)
}
