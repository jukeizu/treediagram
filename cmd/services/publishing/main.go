package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/publishing"
	"github.com/jukeizu/treediagram/services/publishing"
	"github.com/jukeizu/treediagram/services/publishing/handlers"
	"github.com/jukeizu/treediagram/services/publishing/handlers/discord"
	"github.com/jukeizu/treediagram/services/publishing/queue"
	"github.com/jukeizu/treediagram/services/publishing/storage"
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
	Port           int
	MessageStorage mdb.DbConfig
	QueueConfig    queue.QueueConfig
	DiscordConfig  discord.DiscordConfig
}

func main() {
	logger := logging.GetLogger("services.publishing", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	store, err := storage.NewMessageStorage(c.MessageStorage)
	if err != nil {
		panic(err)
	}

	defer store.Close()

	publisherQueue, err := queue.NewQueue(c.QueueConfig)
	if err != nil {
		panic(err)
	}
	publisherQueue = queue.NewQueueLogger(logger, publisherQueue)

	defer publisherQueue.Close()

	discordMessageHandler, err := discord.NewDiscordHandler(c.DiscordConfig)
	if err != nil {
		panic(err)
	}

	queueHandler := handlers.NewQueueHandler(store, discordMessageHandler)

	listener, err := queue.NewQueue(c.QueueConfig)
	if err != nil {
		panic(err)
	}
	listener = queue.NewQueueLogger(logger, listener)
	defer listener.Close()

	listener.Listen(queueHandler)

	service := publishing.NewService(publisherQueue, store)
	service = publishing.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)

	}()

	go func() {
		port := fmt.Sprintf(":%d", c.Port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Log("error", err.Error())

		}

		s := grpc.NewServer()
		pb.RegisterPublishingServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
