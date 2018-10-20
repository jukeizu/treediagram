package discord

import (
	"github.com/bwmarrin/discordgo"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/processing"
)

func mapToUser(discordUser *discordgo.User) *pb.User {
	user := &pb.User{
		Id:   discordUser.ID,
		Name: discordUser.Username,
	}

	return user
}

func mapToUsers(discordUsers []*discordgo.User) []*pb.User {
	users := []*pb.User{}

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
