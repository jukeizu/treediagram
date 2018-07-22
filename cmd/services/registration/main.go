package main

import (
	"net"
	"os"

	pb "github.com/jukeizu/treediagram/api/registration"
	"github.com/jukeizu/treediagram/services/registration"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

func main() {
	logger := logging.GetLogger("services.registration", os.Stdout)

	service, err := registration.NewService()

	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterRegistrationServer(s, registration.GrpcBinding{Service: service})

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Log("error", err.Error())
	}

	err = s.Serve(listener)

	if err != nil {
		logger.Log("error", err.Error())
	}
}
