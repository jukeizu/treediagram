package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/treediagram/services/publishing"
	"github.com/jukeizu/treediagram/services/publishing/handlers"
	"github.com/jukeizu/treediagram/services/publishing/handlers/discord"
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/runner"
)

type TreediagramPublisherConfig struct {
	Port          int
	StorageConfig storage.StorageConfig
	QueueConfig   queue.QueueConfig
	DiscordConfig discord.DiscordConfig
}

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	publisherConfig := TreediagramPublisherConfig{}

	err := config.ReadConfig(serviceArgs.ConfigFile, &publisherConfig)

	if err != nil {
		panic(err)
	}

	store, err := storage.NewStorage(publisherConfig.StorageConfig)

	if err != nil {
		panic(err)
	}

	defer store.Close()

	publisherQueue, err := queue.NewQueue(publisherConfig.QueueConfig)
	publisherQueue = queue.NewQueueLogger(log.With(logger, "component", "treediagram-publisher-queue"), publisherQueue)

	if err != nil {
		panic(err)
	}

	defer publisherQueue.Close()

	discordMessageHandler, err := discord.NewDiscordHandler(publisherConfig.DiscordConfig)

	if err != nil {
		panic(err)
	}

	queueHandler := handlers.NewQueueHandler(store, discordMessageHandler)

	listener, err := queue.NewQueue(publisherConfig.QueueConfig)
	listener = queue.NewQueueLogger(log.With(logger, "component", "treediagram-publisher-queue-listener"), listener)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	listener.Listen(queueHandler)

	sendService := publishing.NewService(publisherQueue, store)
	sendService = publishing.NewLoggingService(log.With(logger, "component", "treediagram-publisher-service"), sendService)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	var requestHandler = publishing.MakeHandler(sendService, httpLogger)
	mux.Handle("/message", requestHandler)

	serviceConfig := config.ServiceConfig{Port: publisherConfig.Port}

	runner.StartService(mux, logger, serviceConfig)
}
