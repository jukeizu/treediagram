package startup

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"

	publishingpb "github.com/jukeizu/treediagram/api/publishing"
	receivingpb "github.com/jukeizu/treediagram/api/receiving"
	registrationpb "github.com/jukeizu/treediagram/api/registration"
	schedulingpb "github.com/jukeizu/treediagram/api/scheduling"
	userpb "github.com/jukeizu/treediagram/api/user"
	"github.com/jukeizu/treediagram/pkg/publishing"
	"github.com/jukeizu/treediagram/pkg/publishing/discord"
	"github.com/jukeizu/treediagram/pkg/receiving"
	"github.com/jukeizu/treediagram/pkg/registration"
	"github.com/jukeizu/treediagram/pkg/scheduling"
	"github.com/jukeizu/treediagram/pkg/user"
	"google.golang.org/grpc"
)

func StartServer(logger log.Logger, config Config) error {
	logger = log.With(logger, "component", "server")

	storage, err := NewStorage(config.DbUrl)
	if err != nil {
		return errors.New("db: " + err.Error())
	}

	defer storage.Close()

	natsConnection, err := natsConnection(config.NatsServers)
	if err != nil {
		return err
	}

	defer natsConnection.Close()

	publishingService := publishing.NewService(natsConnection, storage.MessageStorage)
	publishingService = publishing.NewLoggingService(logger, publishingService)

	discordPublisher, err := discord.NewDiscordPublisher(config.DiscordToken)
	if err != nil {
		return err
	}

	publisher := publishing.NewPublisher(storage.MessageStorage, natsConnection)

	sub, err := publisher.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	receivingService := receiving.NewService(natsConnection)
	receivingService = receiving.NewLoggingService(logger, receivingService)

	registrationService := registration.NewService(storage.CommandStorage)
	registrationService = registration.NewLoggingService(logger, registrationService)

	schedulingService := scheduling.NewService(logger, storage.JobStorage, natsConnection)
	schedulingService = scheduling.NewLoggingService(logger, schedulingService)

	userService := user.NewService(storage.UserStorage)
	userService = user.NewLoggingService(logger, userService)

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
			errChannel <- err
		}

		s := grpc.NewServer()

		publishingpb.RegisterPublishingServer(s, publishingService)
		receivingpb.RegisterReceivingServer(s, receivingService)
		schedulingpb.RegisterSchedulingServer(s, schedulingService)
		registrationpb.RegisterRegistrationServer(s, registration.GrpcBinding{Service: registrationService})
		userpb.RegisterUserServer(s, userService)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	go func() {
		port := fmt.Sprintf(":%d", config.HttpPort)

		schedulingBinding := scheduling.NewHttpBinding(logger, schedulingService)

		logger.Log("transport", "http", "address", port, "msg", "listening")

		http.Handle("/scheduling/", schedulingBinding.MakeHandler())
		http.Handle("/registration/", registration.MakeHandler(registrationService, logger))

		errChannel <- http.ListenAndServe(port, nil)
	}()

	return <-errChannel
}

func natsConnection(servers string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(servers)
	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
