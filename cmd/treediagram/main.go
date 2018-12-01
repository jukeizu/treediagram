package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jukeizu/treediagram/internal/startup"
	nats "github.com/nats-io/go-nats"
	"github.com/oklog/run"
	"github.com/rs/zerolog"
)

var Version = ""

const (
	DefaultGrpcPort                 = 50051
	DefaultReceivingEndpoint        = "localhost:50051"
	DiscordTokenEnvironmentVariable = "TREEDIAGRAM_DISCORD_TOKEN"
)

var (
	flagServer    = false
	flagBot       = false
	flagScheduler = false
	flagVersion   = false
	flagDebug     = false
)

func parseConfig() startup.Config {
	c := startup.Config{}

	flag.IntVar(&c.GrpcPort, "grpc.port", DefaultGrpcPort, "grpc port")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.DbUrl, "db", "localhost", "Database connection url")
	flag.StringVar(&c.DiscordToken, "discord.token", "", "Discord token. This can also be specified via the "+DiscordTokenEnvironmentVariable+" environment variable.")
	flag.StringVar(&c.ReceivingEndpoint, "endpoint", DefaultReceivingEndpoint, "Url of the Receiving service")
	flag.BoolVar(&flagServer, "server", false, "Start as server")
	flag.BoolVar(&flagBot, "bot", false, "Start as bot")
	flag.BoolVar(&flagScheduler, "scheduler", false, "Start as scheduler")
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")

	flag.Parse()

	if c.DiscordToken == "" {
		c.DiscordToken = os.Getenv(DiscordTokenEnvironmentVariable)
	}

	return c
}

func main() {
	config := parseConfig()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("version", Version).
		Logger()

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
