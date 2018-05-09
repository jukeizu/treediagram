package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/treediagram/api"
)

func mapToUser(discordUser *discordgo.User) api.User {
	user := api.User{
		Id:   discordUser.ID,
		Name: discordUser.Username,
	}

	return user
}

func mapToUsers(discordUsers []*discordgo.User) api.Users {
	users := api.Users{}

	for _, discordUser := range discordUsers {
		users = append(users, mapToUser(discordUser))
	}

	return users
}

func getServerId(s *discordgo.Session, channelId string) string {
	channel, err := s.State.Channel(channelId)

	if err != nil || channel == nil {
		return ""
	}

	return channel.GuildID
}
