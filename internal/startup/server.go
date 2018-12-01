package startup

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	intentpb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	processingpb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	schedulingpb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	userpb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
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

	intentConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	intentClient := intentpb.NewIntentRegistryClient(intentConn)

	processorService, err := processor.NewService(conn, storage.ProcessorStorage)
	if err != nil {
		return nil, err
	}

	processor := processor.New(logger, conn, intentClient, storage.ProcessorStorage)
	err = processor.Start()
	if err != nil {
		return nil, err
	}

	intentService := intent.NewService(storage.IntentStorage)
	intentService = intent.NewLoggingService(logger, intentService)

	schedulerService := scheduler.NewService(logger, storage.JobStorage, conn)
	schedulerService = scheduler.NewLoggingService(zerolog.New(os.Stdout).With().Timestamp().Logger(), schedulerService)

	userService := user.NewService(storage.UserStorage)
	userService = user.NewLoggingService(logger, userService)

	grpcServer := grpc.NewServer()

	processingpb.RegisterProcessingServer(grpcServer, processorService)
	schedulingpb.RegisterSchedulingServer(grpcServer, schedulerService)
	intentpb.RegisterIntentRegistryServer(grpcServer, intentService)
	userpb.RegisterUserServer(grpcServer, userService)

	schedulerHttpBinding := scheduler.NewHttpBinding(logger, schedulerService)
	intentHttpBinding := intent.NewHttpBinding(logger, intentService)

	mux := http.NewServeMux()
	mux.Handle("/scheduling/", schedulerHttpBinding.MakeHandler())
	mux.Handle("/intent/", intentHttpBinding.MakeHandler())

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
