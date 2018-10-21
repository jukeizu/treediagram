package startup

import (
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/bot/discord"
	"google.golang.org/grpc"
)

type BotRunner struct {
	Logger     log.Logger
	ClientConn *grpc.ClientConn
	DiscordBot discord.Bot
	quit       chan struct{}
}

func NewBotRunner(logger log.Logger, config Config) (*BotRunner, error) {
	logger = log.With(logger, "component", "bot")

	conn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewProcessingClient(conn)

	handler, err := discord.NewBot(config.DiscordToken, client, logger)
	if err != nil {
		return nil, err
	}

	listenerRunner := &BotRunner{
		Logger:     logger,
		ClientConn: conn,
		DiscordBot: handler,
		quit:       make(chan struct{}),
	}

	return listenerRunner, nil
}

func (r *BotRunner) Start() error {
	r.Logger.Log("msg", "starting")

	err := r.DiscordBot.Open()
	if err != nil {
		return err
	}

	r.Logger.Log("msg", "treediagram-bot has started.")

	<-r.quit

	return nil
}

func (r *BotRunner) Stop() {
	r.Logger.Log("msg", "stopping")

	close(r.quit)
	r.DiscordBot.Close()
	r.ClientConn.Close()
}
