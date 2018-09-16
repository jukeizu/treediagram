package startup

type Config struct {
	GrpcPort          int
	HttpPort          int
	NatsServers       string
	DbUrl             string
	DiscordToken      string
	ReceivingEndpoint string
}
