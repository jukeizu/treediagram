package startup

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/go-kit/kit/log"
	publishingpb "github.com/jukeizu/treediagram/api/publishing"
	receivingpb "github.com/jukeizu/treediagram/api/receiving"
	registrationpb "github.com/jukeizu/treediagram/api/registration"
	schedulingpb "github.com/jukeizu/treediagram/api/scheduling"
	userpb "github.com/jukeizu/treediagram/api/user"
	"github.com/jukeizu/treediagram/services/publishing"
	"github.com/jukeizu/treediagram/services/publishing/discord"
	"github.com/jukeizu/treediagram/services/receiving"
	"github.com/jukeizu/treediagram/services/registration"
	"github.com/jukeizu/treediagram/services/scheduling"
	"github.com/jukeizu/treediagram/services/user"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"google.golang.org/grpc"
)

type ServerRunner struct {
	Logger       log.Logger
	Storage      *Storage
	EncodedConn  *nats.EncodedConn
	Subscription *nats.Subscription
	GrpcServer   *grpc.Server
	HttpServer   *http.Server
	GrpcPort     int
	HttpPort     int
}

func NewServerRunner(logger log.Logger, config Config) (*ServerRunner, error) {
	logger = log.With(logger, "component", "server")

	storage, err := NewStorage(config.DbUrl)
	if err != nil {
		return nil, errors.New("db: " + err.Error())
	}

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, err
	}

	discordPublisher, err := discord.NewDiscordPublisher(config.DiscordToken)
	if err != nil {
		return nil, err
	}

	publisher := publishing.NewPublisher(storage.MessageStorage, conn)
	sub, err := publisher.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		return nil, err
	}

	publishingService := publishing.NewService(conn, storage.MessageStorage)
	publishingService = publishing.NewLoggingService(logger, publishingService)

	receivingService := receiving.NewService(conn)
	receivingService = receiving.NewLoggingService(logger, receivingService)

	registrationService := registration.NewService(storage.CommandStorage)
	registrationService = registration.NewLoggingService(logger, registrationService)

	schedulingService := scheduling.NewService(logger, storage.JobStorage, conn)
	schedulingService = scheduling.NewLoggingService(logger, schedulingService)

	userService := user.NewService(storage.UserStorage)
	userService = user.NewLoggingService(logger, userService)

	grpcServer := grpc.NewServer()

	publishingpb.RegisterPublishingServer(grpcServer, publishingService)
	receivingpb.RegisterReceivingServer(grpcServer, receivingService)
	schedulingpb.RegisterSchedulingServer(grpcServer, schedulingService)
	registrationpb.RegisterRegistrationServer(grpcServer, registrationService)
	userpb.RegisterUserServer(grpcServer, userService)

	schedulingBinding := scheduling.NewHttpBinding(logger, schedulingService)
	registrationBinding := registration.NewHttpBinding(logger, registrationService)

	mux := http.NewServeMux()
	mux.Handle("/scheduling/", schedulingBinding.MakeHandler())
	mux.Handle("/registration/", registrationBinding.MakeHandler())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mux,
	}

	serverRunner := &ServerRunner{
		Logger:       logger,
		Storage:      storage,
		EncodedConn:  conn,
		Subscription: sub,
		GrpcServer:   grpcServer,
		HttpServer:   httpServer,
		GrpcPort:     config.GrpcPort,
		HttpPort:     config.HttpPort,
	}

	return serverRunner, nil
}

func (r *ServerRunner) Start() error {
	r.Logger.Log("msg", "starting")

	errC := make(chan error)

	go func() {
		port := fmt.Sprintf(":%d", r.GrpcPort)
		listener, err := net.Listen("tcp", port)
		if err != nil {
			errC <- err
			return
		}

		r.Logger.Log("transport", "grpc", "address", port, "msg", "listening")
		errC <- r.GrpcServer.Serve(listener)
	}()

	go func() {
		r.Logger.Log("transport", "http", "address", r.HttpServer.Addr, "msg", "listening")
		errC <- r.HttpServer.ListenAndServe()
	}()

	return <-errC
}

func (r *ServerRunner) Stop() {
	r.Logger.Log("msg", "stopping")

	r.Subscription.Unsubscribe()
	r.GrpcServer.GracefulStop()
	r.HttpServer.Shutdown(nil)
	r.EncodedConn.Close()
	r.Storage.Close()
}
