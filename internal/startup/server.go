package startup

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/go-kit/kit/log"
	processingpb "github.com/jukeizu/treediagram/api/processing"
	publishingpb "github.com/jukeizu/treediagram/api/publishing"
	registrationpb "github.com/jukeizu/treediagram/api/registration"
	schedulingpb "github.com/jukeizu/treediagram/api/scheduling"
	userpb "github.com/jukeizu/treediagram/api/user"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/processor/user"
	"github.com/jukeizu/treediagram/publisher"
	"github.com/jukeizu/treediagram/publisher/discord"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
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

	p := publisher.NewPublisher(storage.MessageStorage, conn)
	sub, err := p.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		return nil, err
	}

	publisherService := publisher.NewService(conn, storage.MessageStorage)
	publisherService = publisher.NewLoggingService(logger, publisherService)

	processorService := processor.NewService(conn)
	processorService = processor.NewLoggingService(logger, processorService)

	registryService := registry.NewService(storage.CommandStorage)
	registryService = registry.NewLoggingService(logger, registryService)

	schedulerService := scheduler.NewService(logger, storage.JobStorage, conn)
	schedulerService = scheduler.NewLoggingService(logger, schedulerService)

	userService := user.NewService(storage.UserStorage)
	userService = user.NewLoggingService(logger, userService)

	grpcServer := grpc.NewServer()

	publishingpb.RegisterPublishingServer(grpcServer, publisherService)
	processingpb.RegisterProcessingServer(grpcServer, processorService)
	schedulingpb.RegisterSchedulingServer(grpcServer, schedulerService)
	registrationpb.RegisterRegistrationServer(grpcServer, registryService)
	userpb.RegisterUserServer(grpcServer, userService)

	schedulerHttpBinding := scheduler.NewHttpBinding(logger, schedulerService)
	registryHttpBinding := registry.NewHttpBinding(logger, registryService)

	mux := http.NewServeMux()
	mux.Handle("/scheduling/", schedulerHttpBinding.MakeHandler())
	mux.Handle("/registration/", registryHttpBinding.MakeHandler())

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
