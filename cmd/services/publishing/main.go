package main

import (
	"flag"
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
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

const (
	DefaultPort                     = 50051
	DiscordTokenEnvironmentVariable = "TREEDIAGRAM_DISCORD_TOKEN"
)

type Config struct {
	Port              int
	MessageStorageUrl string
	NatsServers       string
	DiscordToken      string
}

func parseConfig() Config {
	c := Config{}

	flag.IntVar(&c.Port, "p", DefaultPort, "port")
	flag.StringVar(&c.DiscordToken, "discord-token", "", "Discord token. This can also be specified via the "+DiscordTokenEnvironmentVariable+" environment variable.")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.MessageStorageUrl, "db", "localhost", "Database connection url")

	flag.Parse()

	if c.DiscordToken == "" {
		c.DiscordToken = os.Getenv(DiscordTokenEnvironmentVariable)
	}

	return c
}

func main() {
	logger := logging.GetLogger("services.publishing", os.Stdout)

	config := parseConfig()

	store, err := publishing.NewMessageStorage(config.MessageStorageUrl)
	if err != nil {
		logger.Log("db error", err)
		os.Exit(1)
	}

	defer store.Close()

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer conn.Close()

	discordPublisher, err := discord.NewDiscordPublisher(config.DiscordToken)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	publisher := publishing.NewPublisher(store, conn)

	sub, err := publisher.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer sub.Unsubscribe()

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)

	}()

	service := publishing.NewService(conn, store)
	service = publishing.NewLoggingService(logger, service)

	go func() {
		port := fmt.Sprintf(":%d", config.Port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			errChannel <- err
		}

		s := grpc.NewServer()
		pb.RegisterPublishingServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
