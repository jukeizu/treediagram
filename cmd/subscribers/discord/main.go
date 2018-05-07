package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	base "github.com/jukeizu/client-base"
	"github.com/jukeizu/treediagram/api"
	"github.com/jukeizu/treediagram/subscribers/discord"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
)

type Config struct {
	DiscordSubscriberConfig discord.DiscordSubscriberConfig
	ClientConfig            base.ClientConfig
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

	client := api.NewClient(c.ClientConfig)

	handler, err := discord.NewDiscordSubscriber(c.DiscordSubscriberConfig, client, logger)

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
