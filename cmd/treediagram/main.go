package main

import (
	"flag"
	"os"

	nats "github.com/nats-io/go-nats"
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

type Config struct {
	GrpcPort          int
	HttpPort          int
	NatsServers       string
	DbUrl             string
	DiscordToken      string
	ReceivingEndpoint string
}

func parseConfig() Config {
	c := Config{}

	flag.IntVar(&c.GrpcPort, "grpc-port", DefaultGrpcPort, "grpc port")
	flag.IntVar(&c.HttpPort, "http-port", DefaultHttpPort, "http port")
	flag.StringVar(&c.NatsServers, "nats", nats.DefaultURL, "NATS servers")
	flag.StringVar(&c.DbUrl, "db", "localhost", "Database connection url")
	flag.StringVar(&c.DiscordToken, "discord-token", "", "Discord token. This can also be specified via the "+DiscordTokenEnvironmentVariable+" environment variable.")
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

	errChannel := make(chan error)

	if !startServer && !startListener && !startScheduler {
		startServer = true
		startListener = true
		startScheduler = true
	}

	if startServer {
		go func() {
			errChannel <- StartServer(logger, config)
		}()
	}

	if startListener {
		go func() {
			errChannel <- StartListener(logger, config)
		}()
	}

	if startScheduler {
		go func() {
			errChannel <- StartScheduler(logger, config)
		}()
	}

	logger.Log("stopped", <-errChannel)
}
