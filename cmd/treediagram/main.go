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
	"github.com/shawntoffel/services-core/logging"
)

const (
	DefaultGrpcPort                 = 50051
	DefaultHttpPort                 = 10001
	DefaultReceivingEndpoint        = "localhost:50051"
	DiscordTokenEnvironmentVariable = "TREEDIAGRAM_DISCORD_TOKEN"
)

var (
	startServer    = false
	startListener  = false
	startScheduler = false
)

func parseConfig() startup.Config {
	c := startup.Config{}

	flag.IntVar(&c.GrpcPort, "grpc.port", DefaultGrpcPort, "grpc port")
	flag.IntVar(&c.HttpPort, "http.port", DefaultHttpPort, "http port")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.DbUrl, "db", "localhost", "Database connection url")
	flag.StringVar(&c.DiscordToken, "discord.token", "", "Discord token. This can also be specified via the "+DiscordTokenEnvironmentVariable+" environment variable.")
	flag.StringVar(&c.ReceivingEndpoint, "endpoint", DefaultReceivingEndpoint, "Url of the Receiving service")
	flag.BoolVar(&startServer, "server", false, "Start as server")
	flag.BoolVar(&startListener, "listener", false, "Start as listener")
	flag.BoolVar(&startScheduler, "scheduler", false, "Start as scheduler")

	flag.Parse()

	if c.DiscordToken == "" {
		c.DiscordToken = os.Getenv(DiscordTokenEnvironmentVariable)
	}

	return c
}

func main() {
	logger := logging.GetLogger("treediagram", os.Stdout)

	config := parseConfig()

	if !startServer && !startListener && !startScheduler {
		startServer = true
		startListener = true
		startScheduler = true
	}

	g := run.Group{}

	if startScheduler {
		s, err := startup.NewSchedulerRunner(logger, config)
		if err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}

		g.Add(func() error {
			return s.Start()
		}, func(error) {
			s.Stop()
		})
	}

	if startListener {
		l, err := startup.NewListenerRunner(logger, config)
		if err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}

		g.Add(func() error {
			return l.Start()
		}, func(error) {
			l.Stop()
		})
	}

	if startServer {
		s, err := startup.NewServerRunner(logger, config)
		if err != nil {
			logger.Log("error", err)
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

	logger.Log("stopped", g.Run())
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
