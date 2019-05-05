package discord

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/jukeizu/contract"
	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/rs/zerolog"
)

func mapToPbUser(discordUser *discordgo.User) *processingpb.User {
	user := &processingpb.User{
		Id:            discordUser.ID,
		Name:          discordUser.Username,
		Discriminator: discordUser.Discriminator,
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

func mapToMessageSend(contractMessage *contract.Message) (*discordgo.MessageSend, error) {
	if contractMessage == nil {
		return nil, nil
	}

	messageSend := discordgo.MessageSend{
		Content: contractMessage.Content,
		Embed:   mapToMessageEmbed(contractMessage.Embed),
		Files:   mapToFiles(contractMessage.Files),
	}

	return &messageSend, nil
}

func mapToMessageEmbed(contractEmbed *contract.Embed) *discordgo.MessageEmbed {
	if contractEmbed == nil {
		return nil
	}

	embed := discordgo.MessageEmbed{
		URL:         contractEmbed.Url,
		Title:       contractEmbed.Title,
		Description: contractEmbed.Description,
		Timestamp:   contractEmbed.Timestamp,
		Color:       contractEmbed.Color,
		Footer:      mapToMessageEmbedFooter(contractEmbed.Footer),
		Image:       &discordgo.MessageEmbedImage{URL: contractEmbed.ImageUrl},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: contractEmbed.ThumbnailUrl},
		Video:       &discordgo.MessageEmbedVideo{URL: contractEmbed.VideoUrl},
		Author:      mapToMessageEmbedAuthor(contractEmbed.Author),
		Fields:      mapToMessageEmbedFields(contractEmbed.Fields),
	}

	return &embed
}

func mapToMessageEmbedFooter(contractEmbedFooter *contract.EmbedFooter) *discordgo.MessageEmbedFooter {
	if contractEmbedFooter == nil {
		return nil
	}

	footer := discordgo.MessageEmbedFooter{
		Text:    contractEmbedFooter.Text,
		IconURL: contractEmbedFooter.IconUrl,
	}

	return &footer
}

func mapToMessageEmbedAuthor(contractEmbedAuthor *contract.EmbedAuthor) *discordgo.MessageEmbedAuthor {
	if contractEmbedAuthor == nil {
		return nil
	}

	author := discordgo.MessageEmbedAuthor{
		URL:  contractEmbedAuthor.Url,
		Name: contractEmbedAuthor.Name,
	}

	return &author
}

func mapToMessageEmbedFields(contractFields []*contract.EmbedField) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	if contractFields == nil {
		return fields
	}

	for _, contractField := range contractFields {
		field := discordgo.MessageEmbedField{
			Name:   contractField.Name,
			Value:  contractField.Value,
			Inline: contractField.Inline,
		}

		fields = append(fields, &field)
	}

	return fields
}

func mapToFiles(contractFiles []*contract.File) []*discordgo.File {
	files := []*discordgo.File{}

	if contractFiles == nil {
		return files
	}

	for _, contractFile := range contractFiles {
		file := discordgo.File{
			Name:        contractFile.Name,
			ContentType: contractFile.ContentType,
			Reader:      bytes.NewReader(contractFile.Bytes),
		}

		files = append(files, &file)
	}

	return files
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
