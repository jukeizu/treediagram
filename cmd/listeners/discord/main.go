package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/jukeizu/treediagram/listeners/discord"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

type Config struct {
	DiscordToken      string
	ReceivingEndpoint string
}

var commandArgs command.CommandArgs

func init() {
	commandArgs = command.ParseArgs()
}

func main() {
	logger := logging.GetLogger("treediagram-bot", os.Stdout)

	c := Config{}

	err := config.ReadConfig(commandArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(c.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewReceivingClient(conn)

	handler, err := discord.NewDiscordListener(c.DiscordToken, client, logger)
	err = handler.Open()

	if err != nil {
		panic(err)
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
