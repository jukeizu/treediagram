package discord

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/rs/zerolog"
)

func mapToPbUser(discordUser *discordgo.User) *processingpb.User {
	user := &processingpb.User{
		Id:   discordUser.ID,
		Name: discordUser.Username,
	}

	return user
}

func mapToPbUsers(discordUsers []*discordgo.User) []*processingpb.User {
	users := []*processingpb.User{}

	for _, discordUser := range discordUsers {
		users = append(users, mapToPbUser(discordUser))
	}

	return users
}

func mapToPbServers(userId string, guilds []*discordgo.Guild) []*processingpb.Server {
	servers := []*processingpb.Server{}

	for _, guild := range guilds {
		for _, member := range guild.Members {
			if member.User.ID == userId {
				server := processingpb.Server{
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

func mapToMessageSend(message *processingpb.MessageReply) (*discordgo.MessageSend, error) {
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
