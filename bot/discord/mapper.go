package discord

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	"github.com/rs/zerolog"
)

func mapToPbUser(discordUser *discordgo.User) *pb.User {
	user := &pb.User{
		Id:   discordUser.ID,
		Name: discordUser.Username,
	}

	return user
}

func mapToPbUsers(discordUsers []*discordgo.User) []*pb.User {
	users := []*pb.User{}

	for _, discordUser := range discordUsers {
		users = append(users, mapToPbUser(discordUser))
	}

	return users
}

func mapToPbServers(userId string, guilds []*discordgo.Guild) []*pb.Server {
	servers := []*pb.Server{}

	for _, guild := range guilds {
		for _, member := range guild.Members {
			if member.User.ID == userId {
				server := pb.Server{
					Id:   guild.ID,
					Name: guild.Name,
				}

				servers = append(servers, &server)
				break
			}
		}
	}

	return servers
}

func mapToMessageSend(message *pb.MessageReply) (*discordgo.MessageSend, error) {
	if message == nil {
		return nil, nil
	}

	messageSend := discordgo.MessageSend{}

	err := json.Unmarshal([]byte(message.Content), &messageSend)
	if err != nil {
		return nil, err
	}

	return &messageSend, nil
}

func mapToLevel(dgoLevel int) zerolog.Level {
	switch dgoLevel {
	case discordgo.LogError:
		return zerolog.ErrorLevel
	case discordgo.LogWarning:
		return zerolog.WarnLevel
	case discordgo.LogInformational:
		return zerolog.InfoLevel
	case discordgo.LogDebug:
		return zerolog.DebugLevel
	}

	return zerolog.InfoLevel
}
