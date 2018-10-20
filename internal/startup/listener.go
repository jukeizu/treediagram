package startup

import (
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/processing"
	"github.com/jukeizu/treediagram/listener/discord"
	"google.golang.org/grpc"
)

type ListenerRunner struct {
	Logger          log.Logger
	ClientConn      *grpc.ClientConn
	DiscordListener discord.DiscordListener
	quit            chan struct{}
}

func NewListenerRunner(logger log.Logger, config Config) (*ListenerRunner, error) {
	logger = log.With(logger, "component", "listener")

	conn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewProcessingClient(conn)

	handler, err := discord.NewDiscordListener(config.DiscordToken, client, logger)
	if err != nil {
		return nil, err
	}

	listenerRunner := &ListenerRunner{
		Logger:          logger,
		ClientConn:      conn,
		DiscordListener: handler,
		quit:            make(chan struct{}),
	}

	return listenerRunner, nil
}

func (r *ListenerRunner) Start() error {
	r.Logger.Log("msg", "starting")

	err := r.DiscordListener.Open()
	if err != nil {
		return err
	}

	r.Logger.Log("msg", "treediagram-bot has started.")

	<-r.quit

	return nil
}

func (r *ListenerRunner) Stop() {
	r.Logger.Log("msg", "stopping")

	close(r.quit)
	r.DiscordListener.Close()
	r.ClientConn.Close()
}
