package startup

import (
	"bufio"
	"os"
	"sync"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/bot/discord"
	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type BotRunner struct {
	Logger      zerolog.Logger
	ClientConn  *grpc.ClientConn
	DiscordBot  discord.Bot
	quit        chan struct{}
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
	WaitGroup   *sync.WaitGroup
}

func NewBotRunner(logger zerolog.Logger, config Config) (*BotRunner, error) {
	logger = logger.With().Str("component", "bot").Logger()

	wg := sync.WaitGroup{}
	wg.Add(1)

	nc, err := nats.Connect(config.NatsServers,
		nats.ClosedHandler(func(_ *nats.Conn) {
			wg.Done()
		}))
	if err != nil {
		return nil, err
	}

	queue, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	clientConn, err := grpc.Dial(config.ReceivingEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := processingpb.NewProcessingClient(clientConn)

	if config.DiscordTokenFile != "" {
		f, _ := os.Open(config.DiscordTokenFile)
		if err != nil {
			return nil, err
		}

		s := bufio.NewScanner(f)
		s.Scan()
		if err := s.Err(); err != nil {
			return nil, err
		}
		config.DiscordToken = s.Text()
	}

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
	r.Logger.Info().Msg("starting")

	err := r.DiscordBot.Open()
	if err != nil {
		return err
	}

	r.Logger.Info().Msg("treediagram-bot has started")

	<-r.quit

	return nil
}

func (r *BotRunner) Stop() {
	r.Logger.Info().Msg("stopping")

	r.EncodedConn.Drain()
	r.Conn.Drain()
	r.WaitGroup.Wait()

	close(r.quit)
	r.DiscordBot.Close()
	r.ClientConn.Close()
}
