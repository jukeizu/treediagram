package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/services/receiving"
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

	var service receiving.Service
	service, err = receiving.NewService(treediagramConfig.RabbitMqConfig)

	if err != nil {
		panic(err)
	}

	service = receiving.NewLoggingService(log.With(logger, "component", "treediagram"), service)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	var handler = receiving.MakeHandler(service, httpLogger)
	mux.Handle("/treediagram", handler)
	mux.Handle("/treediagram/", handler)

	serviceConfig := config.ServiceConfig{Port: treediagramConfig.Port}

	runner.StartService(mux, logger, serviceConfig)
}
