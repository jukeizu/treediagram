package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	"github.com/jukeizu/treediagram/internal"
	"github.com/jukeizu/treediagram/internal/startup"
	nats "github.com/nats-io/nats.go"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/grpclog"
)

const (
	DefaultGrpcPort                     = 50051
	DefaultHttpPort                     = 10002
	DefaultReceivingEndpoint            = "localhost:50051"
	DiscordTokenEnvironmentVariable     = "TREEDIAGRAM_DISCORD_TOKEN"
	DiscordTokenFileEnvironmentVariable = "TREEDIAGRAM_DISCORD_TOKEN_FILE"
	DefaultIntentUrl                    = "http://localhost:8080/intent.yml"
)

var (
	flagServer    = false
	flagBot       = false
	flagScheduler = false
	flagMigrate   = false
	flagVersion   = false
	flagDebug     = false
)

func parseConfig() startup.Config {
	c := startup.Config{}

	flag.IntVar(&c.GrpcPort, "grpc.port", DefaultGrpcPort, "grpc port")
	flag.IntVar(&c.HttpPort, "http.port", DefaultHttpPort, "http port")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.DbUrl, "db", "root@localhost:26257", "Database connection url")
	flag.StringVar(&c.ReceivingEndpoint, "endpoint", DefaultReceivingEndpoint, "Url of the Receiving service")
	flag.BoolVar(&flagServer, "server", false, "Start as server")
	flag.BoolVar(&flagBot, "bot", false, "Start as bot")
	flag.BoolVar(&flagScheduler, "scheduler", false, "Start as scheduler")
	flag.BoolVar(&flagMigrate, "migrate", false, "Run db migrations")
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")
	flag.StringVar(&c.IntentEndpoint, "intent.url", DefaultIntentUrl, "intent url")

	flag.Parse()

	c.DiscordTokenFile = os.Getenv(DiscordTokenFileEnvironmentVariable)
	c.DiscordToken = os.Getenv(DiscordTokenEnvironmentVariable)

	return c
}

func main() {
	config := parseConfig()

	if flagVersion {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("version", internal.Version).
		Logger()

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	if flagMigrate {
		migrationRunner, err := startup.NewMigrationRunner(logger, config.DbUrl)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("")
			os.Exit(1)
		}

		err = migrationRunner.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("migrations did not complete")
			os.Exit(1)
		}
	}

	if !flagServer && !flagBot && !flagScheduler {
		flagServer = true
		flagBot = true
		flagScheduler = true
	}

	g := run.Group{}

	if flagScheduler {
		s, err := startup.NewSchedulerRunner(logger, config)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("")
			os.Exit(1)
		}

		g.Add(func() error {
			return s.Start()
		}, func(error) {
			s.Stop()
		})
	}

	if flagBot {
		l, err := startup.NewBotRunner(logger, config)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("")
			os.Exit(1)
		}

		g.Add(func() error {
			return l.Start()
		}, func(error) {
			l.Stop()
		})
	}

	if flagServer {
		s, err := startup.NewServerRunner(logger, config)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("")
			os.Exit(1)
		}

		g.Add(func() error {
			return s.Start()
		}, func(error) {
			s.Stop()
		})
	}

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}
