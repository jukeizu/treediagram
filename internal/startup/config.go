package startup

type Config struct {
	GrpcPort          int
	NatsServers       string
	DbUrl             string
	DiscordToken      string
	ReceivingEndpoint string
}
