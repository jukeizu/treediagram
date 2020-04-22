package startup

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/schedulingpb"
	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	"github.com/jukeizu/treediagram/pkg/builtin"
	"github.com/jukeizu/treediagram/pkg/builtin/help"
	bintent "github.com/jukeizu/treediagram/pkg/builtin/intent"
	serverselect "github.com/jukeizu/treediagram/pkg/builtin/server_select"
	"github.com/jukeizu/treediagram/pkg/builtin/stats"
	"github.com/jukeizu/treediagram/pkg/intent"
	"github.com/jukeizu/treediagram/pkg/processor"
	"github.com/jukeizu/treediagram/pkg/scheduler"
	"github.com/jukeizu/treediagram/pkg/user"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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
	Builtin     builtin.HttpServer
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

	grpcConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	intentClient := intentpb.NewIntentRegistryClient(grpcConn)
	userClient := userpb.NewUserClient(grpcConn)
	processingClient := processingpb.NewProcessingClient(grpcConn)
	schedulerClient := schedulingpb.NewSchedulingClient(grpcConn)

	processorService, err := processor.NewService(conn, storage.ProcessorRepository)
	if err != nil {
		return nil, err
	}

	processor := processor.New(logger, conn, intentClient, userClient, schedulerClient, storage.ProcessorRepository)
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

	helpHandler := help.NewHelpHandler(logger, intentClient)
	serverSelectHandler := serverselect.NewServerSelectHandler(logger, userClient)
	statsHandler := stats.NewStatsHandler(logger, processingClient)
	intentHandler := bintent.NewIntentHandler(logger, intentClient)

	builtinServer := builtin.NewHttpServer(logger, fmt.Sprintf(":%d", config.HttpPort))
	builtinServer.RegisterHandlers(
		helpHandler,
		serverSelectHandler,
		statsHandler,
		intentHandler,
	)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    5 * time.Minute,
				Timeout: 10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			},
		))

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
		Builtin:     builtinServer,
	}

	return serverRunner, nil
}

func (r *ServerRunner) Start() error {
	r.Logger.Info().Msg("starting")

	errC := make(chan error)

	go func() {
		err := r.Builtin.Start()
		if err != nil {
			errC <- err
			return
		}
	}()

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

	r.Builtin.Stop()
	r.EncodedConn.Drain()
	r.Conn.Drain()

	r.Processor.Stop()

	r.GrpcServer.GracefulStop()

	r.WaitGroup.Wait()
}
