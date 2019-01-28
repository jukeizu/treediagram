package startup

import (
	"errors"
	"fmt"
	"net"
	"sync"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type ServerRunner struct {
	Logger      zerolog.Logger
	Storage     *Storage
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
	GrpcServer  *grpc.Server
	GrpcPort    int
	WaitGroup   *sync.WaitGroup
	Processor   processor.Processor
}

func NewServerRunner(logger zerolog.Logger, config Config) (*ServerRunner, error) {
	logger = logger.With().Str("component", "server").Logger()

	storage, err := NewStorage(logger, config.DbUrl)
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

	grpcConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	intentClient := intentpb.NewIntentRegistryClient(grpcConn)
	userClient := userpb.NewUserClient(grpcConn)

	processorService, err := processor.NewService(conn, storage.ProcessorRepository)
	if err != nil {
		return nil, err
	}

	processor := processor.New(logger, conn, intentClient, userClient, storage.ProcessorRepository)
	err = processor.Start()
	if err != nil {
		return nil, err
	}

	intentService := intent.NewService(storage.IntentRepository)
	intentService = intent.NewLoggingService(logger, intentService)

	schedulerService, err := scheduler.NewService(logger, storage.SchedulerRepository, conn)
	if err != nil {
		return nil, err
	}
	schedulerService = scheduler.NewLoggingService(logger, schedulerService)

	userService := user.NewService(storage.UserRepository)
	userService = user.NewLoggingService(logger, userService)

	grpcServer := grpc.NewServer()

	processingpb.RegisterProcessingServer(grpcServer, processorService)
	schedulingpb.RegisterSchedulingServer(grpcServer, schedulerService)
	intentpb.RegisterIntentRegistryServer(grpcServer, intentService)
	userpb.RegisterUserServer(grpcServer, userService)

	serverRunner := &ServerRunner{
		Logger:      logger,
		Storage:     storage,
		Conn:        nc,
		EncodedConn: conn,
		GrpcServer:  grpcServer,
		GrpcPort:    config.GrpcPort,
		WaitGroup:   &wg,
		Processor:   processor,
	}

	return serverRunner, nil
}

func (r *ServerRunner) Start() error {
	r.Logger.Info().Msg("starting")

	errC := make(chan error)

	go func() {
		port := fmt.Sprintf(":%d", r.GrpcPort)
		listener, err := net.Listen("tcp", port)
		if err != nil {
			errC <- err
			return
		}

		r.Logger.Info().
			Str("transport", "grpc").
			Str("address", port).
			Msg("listening")
		errC <- r.GrpcServer.Serve(listener)
	}()

	return <-errC
}

func (r *ServerRunner) Stop() {
	r.Logger.Info().Msg("stopping")

	r.EncodedConn.Drain()
	r.Conn.Drain()

	r.Processor.Stop()

	r.GrpcServer.GracefulStop()

	r.WaitGroup.Wait()
}
