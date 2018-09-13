package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/jukeizu/treediagram/listeners/discord"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

const (
	DefaultReceivingEndpoint        = "localhost:50052"
	DiscordTokenEnvironmentVariable = "TREEDIAGRAM_DISCORD_TOKEN"
)

type Config struct {
	DiscordToken      string
	ReceivingEndpoint string
}

func parseConfig() Config {
	c := Config{}

	flag.StringVar(&c.DiscordToken, "discord-token", "", "Discord token. This can also be specified via the "+DiscordTokenEnvironmentVariable+" environment variable.")
	flag.StringVar(&c.ReceivingEndpoint, "endpoint", DefaultReceivingEndpoint, "Url of the Receiving service")
	flag.Parse()

	if c.DiscordToken == "" {
		c.DiscordToken = os.Getenv(DiscordTokenEnvironmentVariable)
	}

	return c
}

func main() {
	logger := logging.GetLogger("treediagram-bot", os.Stdout)

	c := parseConfig()

	conn, err := grpc.Dial(c.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer conn.Close()

	client := pb.NewReceivingClient(conn)

	handler, err := discord.NewDiscordListener(c.DiscordToken, client, logger)

	err = handler.Open()
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	logger.Log("msg", "treediagram-bot has started.")

	defer handler.Close()

	cmdErrs := make(chan error, 1)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		cmdErrs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("stopped", <-cmdErrs)
}
