package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/services/treediagram"
	"github.com/shawntoffel/rabbitmq"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/runner"
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
	logger := log.NewJSONLogger(os.Stdout)

	treediagramConfig := TreediagramConfig{}

	err := config.ReadConfig(serviceArgs.ConfigFile, &treediagramConfig)

	if err != nil {
		panic(err)
	}

	var service treediagram.Service
	service, err = treediagram.NewService(treediagramConfig.RabbitMqConfig)

	if err != nil {
		panic(err)
	}

	service = treediagram.NewLoggingService(log.With(logger, "component", "treediagram"), service)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	var handler = treediagram.MakeHandler(service, httpLogger)
	mux.Handle("/treediagram", handler)
	mux.Handle("/treediagram/", handler)

	serviceConfig := config.ServiceConfig{Port: treediagramConfig.Port}

	runner.StartService(mux, logger, serviceConfig)
}
