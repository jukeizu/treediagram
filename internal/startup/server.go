package startup

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-kit/kit/log"
	processingpb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	publishingpb "github.com/jukeizu/treediagram/api/protobuf-spec/publishing"
	receivingpb "github.com/jukeizu/treediagram/api/protobuf-spec/receiving"
	registrationpb "github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	schedulingpb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	userpb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/publisher"
	"github.com/jukeizu/treediagram/publisher/discord"
	"github.com/jukeizu/treediagram/receiver"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	nats "github.com/nats-io/go-nats"
	"google.golang.org/grpc"
)

type ServerRunner struct {
	Logger      log.Logger
	Storage     *Storage
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
	GrpcServer  *grpc.Server
	HttpServer  *http.Server
	GrpcPort    int
	HttpPort    int
	WaitGroup   *sync.WaitGroup
	Processor   processor.Processor
}

func NewServerRunner(logger log.Logger, config Config) (*ServerRunner, error) {
	logger = log.With(logger, "component", "server")

	storage, err := NewStorage(config.DbUrl)
	if err != nil {
		return nil, errors.New("db: " + err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	nc, err := nats.Connect(config.NatsServers,
		nats.ClosedHandler(func(_ *nats.Conn) {
			wg.Done()
		}))

	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	discordPublisher, err := discord.NewDiscordPublisher(config.DiscordToken)
	if err != nil {
		return nil, err
	}

	p := publisher.NewPublisher(storage.MessageStorage, conn)
	_, err = p.Subscribe(discord.DiscordPublisherSubject, discordPublisher.Publish)
	if err != nil {
		return nil, err
	}

	publisherService := publisher.NewService(conn, storage.MessageStorage)
	publisherService = publisher.NewLoggingService(logger, publisherService)

	registryConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	registryClient := registrationpb.NewRegistrationClient(registryConn)

	processorService, err := processor.NewService(logger, conn, registryClient)
	if err != nil {
		return nil, err
	}
	processorService = processor.NewLoggingService(logger, processorService)

	processor := processor.New(logger, conn, registryClient)
	err = processor.Start()
	if err != nil {
		return nil, err
	}

	registryService := registry.NewService(storage.IntentStorage)
	registryService = registry.NewLoggingService(logger, registryService)

	schedulerService := scheduler.NewService(logger, storage.JobStorage, conn)
	schedulerService = scheduler.NewLoggingService(logger, schedulerService)

	userService := user.NewService(storage.UserStorage)
	userService = user.NewLoggingService(logger, userService)

	receiverService := receiver.NewService(conn)

	grpcServer := grpc.NewServer()

	publishingpb.RegisterPublishingServer(grpcServer, publisherService)
	processingpb.RegisterProcessingServer(grpcServer, processorService)
	schedulingpb.RegisterSchedulingServer(grpcServer, schedulerService)
	registrationpb.RegisterRegistrationServer(grpcServer, registryService)
	userpb.RegisterUserServer(grpcServer, userService)
	receivingpb.RegisterReceivingServer(grpcServer, receiverService)

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
		Logger:      logger,
		Storage:     storage,
		Conn:        nc,
		EncodedConn: conn,
		GrpcServer:  grpcServer,
		HttpServer:  httpServer,
		GrpcPort:    config.GrpcPort,
		HttpPort:    config.HttpPort,
		WaitGroup:   &wg,
		Processor:   processor,
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

	r.EncodedConn.Drain()
	r.Conn.Drain()

	r.Processor.Stop()

	r.GrpcServer.GracefulStop()
	r.HttpServer.Shutdown(nil)
	r.Storage.Close()

	r.WaitGroup.Wait()
}
