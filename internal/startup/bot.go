package startup

import (
	"sync"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/jukeizu/treediagram/bot/discord"
	nats "github.com/nats-io/go-nats"
	"google.golang.org/grpc"
)

type BotRunner struct {
	Logger      log.Logger
	ClientConn  *grpc.ClientConn
	DiscordBot  discord.Bot
	quit        chan struct{}
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
	WaitGroup   *sync.WaitGroup
}

func NewBotRunner(logger log.Logger, config Config) (*BotRunner, error) {
	logger = log.With(logger, "component", "bot")

	wg := sync.WaitGroup{}
	wg.Add(1)

	nc, err := nats.Connect(config.NatsServers,
		nats.ClosedHandler(func(_ *nats.Conn) {
			wg.Done()
		}))

	queue, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	clientConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewProcessingClient(clientConn)

	handler, err := discord.NewBot(config.DiscordToken, client, queue, logger)
	if err != nil {
		return nil, err
	}

	listenerRunner := &BotRunner{
		Logger:      logger,
		ClientConn:  clientConn,
		Conn:        nc,
		EncodedConn: queue,
		DiscordBot:  handler,
		quit:        make(chan struct{}),
		WaitGroup:   &wg,
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

	r.EncodedConn.Drain()
	r.Conn.Drain()
	r.WaitGroup.Wait()

	close(r.quit)
	r.DiscordBot.Close()
	r.ClientConn.Close()
}
