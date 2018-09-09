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
	DefaultPort = 50051
)

type Config struct {
	Port              int
	MessageStorageUrl string
	NatsServers       string
	DiscordToken      string
}

func main() {
	logger := logging.GetLogger("services.publishing", os.Stdout)

	port := flag.Int("port", DefaultPort, "port")
	token := flag.String("token", "", "Discord token")
	natsServers := flag.String("nats", nats.DefaultURL, "NATS servers")
	storageUrl := flag.String("db", "localhost", "Database connection url")

	flag.Parse()

	c := Config{
		Port:              *port,
		MessageStorageUrl: *storageUrl,
		NatsServers:       *natsServers,
		DiscordToken:      *token,
	}

	store, err := publishing.NewMessageStorage(c.MessageStorageUrl)
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

	discordPublisher, err := discord.NewDiscordPublisher(c.DiscordToken)
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
