package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/publishing"
	"github.com/jukeizu/treediagram/services/publishing"
	"github.com/jukeizu/treediagram/services/publishing/discord"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
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
	NatsServers    string
	DiscordConfig  discord.DiscordConfig
}

func main() {
	logger := logging.GetLogger("services.publishing", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	store, err := publishing.NewMessageStorage(c.MessageStorage)
	if err != nil {
		panic(err)
	}

	defer store.Close()

	nc, err := nats.Connect(c.NatsServers)
	if err != nil {
		panic(err)
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	discordPublisher, err := discord.NewDiscordPublisher(c.DiscordConfig)
	if err != nil {
		panic(err)
	}

	publisher := publishing.NewPublisher(store, conn)

	sub, err := publisher.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		panic(err)
	}

	defer sub.Unsubscribe()

	service := publishing.NewService(conn, store)
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
