package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/receiving"
	"github.com/jukeizu/treediagram/pkg/listening/discord"
	"google.golang.org/grpc"
)

func StartListener(logger log.Logger, config Config) error {
	logger = log.With(logger, "component", "listener")
	conn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := pb.NewReceivingClient(conn)

	handler, err := discord.NewDiscordListener(config.DiscordToken, client, logger)

	err = handler.Open()
	if err != nil {
		return err
	}

	logger.Log("msg", "treediagram-bot has started.")

	defer handler.Close()

	cmdErrs := make(chan error, 1)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		cmdErrs <- fmt.Errorf("%s", <-c)
	}()

	return <-cmdErrs
}
