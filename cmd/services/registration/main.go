package main

import (
	"fmt"
	"net"
	"os"

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
	Port           int
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

	logger.Log("storage connected", "true")

	service, err := registration.NewService(storage)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterRegistrationServer(s, registration.GrpcBinding{Service: service})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		logger.Log("error", err.Error())
	}

	err = s.Serve(listener)
	if err != nil {
		logger.Log("error", err.Error())
	}
}
