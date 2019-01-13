package startup

type Config struct {
	GrpcPort          int
	NatsServers       string
	MdbUrl            string
	DbUrl             string
	DiscordToken      string
	ReceivingEndpoint string
}
