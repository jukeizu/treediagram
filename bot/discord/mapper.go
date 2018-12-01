package discord

import (
	"bytes"

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

func mapToEmbed(embed *pb.Embed) *discordgo.MessageEmbed {
	if embed == nil {
		return nil
	}

	discordEmbed := discordgo.MessageEmbed{
		URL:         embed.Url,
		Type:        embed.Type,
		Title:       embed.Title,
		Description: embed.Description,
		Timestamp:   embed.Timestamp,
		Color:       int(embed.Color),
		Footer:      mapFooter(embed.Footer),
		Image:       mapImage(embed.Image),
		Thumbnail:   mapThumbnail(embed.Thumbnail),
		Video:       mapVideo(embed.Video),
		Provider:    mapProvider(embed.Provider),
		Author:      mapAuthor(embed.Author),
		Fields:      mapFields(embed.Fields),
	}

	return &discordEmbed
}

func mapToMessageSend(message *pb.MessageReply) *discordgo.MessageSend {
	if message == nil {
		return nil
	}

	mappedEmbed := mapToEmbed(message.Embed)

	messageSend := discordgo.MessageSend{
		Content: message.Content,
		Embed:   mappedEmbed,
	}

	if message.Files != nil {
		for _, file := range message.Files {
			discordFile := mapToFile(file)
			messageSend.Files = append(messageSend.Files, discordFile)
		}
	}

	return &messageSend
}

func mapToFile(file *pb.File) *discordgo.File {
	if file == nil {
		return nil
	}
	discordFile := discordgo.File{
		Name:        file.Name,
		ContentType: file.ContentType,
		Reader:      bytes.NewReader(file.Bytes),
	}

	return &discordFile
}

func mapFooter(footer *pb.EmbedFooter) *discordgo.MessageEmbedFooter {
	if footer == nil {
		return nil
	}

	discordFooter := discordgo.MessageEmbedFooter{
		Text:         footer.Text,
		IconURL:      footer.IconUrl,
		ProxyIconURL: footer.ProxyIconUrl,
	}

	return &discordFooter
}

func mapImage(image *pb.EmbedImage) *discordgo.MessageEmbedImage {
	if image == nil {
		return nil
	}

	discordImage := discordgo.MessageEmbedImage{
		URL:      image.Url,
		ProxyURL: image.ProxyUrl,
		Width:    int(image.Width),
		Height:   int(image.Height),
	}

	return &discordImage
}

func mapThumbnail(thumbnail *pb.EmbedThumbnail) *discordgo.MessageEmbedThumbnail {
	if thumbnail == nil {
		return nil
	}

	discordThumbnail := discordgo.MessageEmbedThumbnail{
		URL:      thumbnail.Url,
		ProxyURL: thumbnail.ProxyUrl,
		Width:    int(thumbnail.Width),
		Height:   int(thumbnail.Height),
	}

	return &discordThumbnail
}

func mapVideo(video *pb.EmbedVideo) *discordgo.MessageEmbedVideo {
	if video == nil {
		return nil
	}

	discordVideo := discordgo.MessageEmbedVideo{
		URL:      video.Url,
		ProxyURL: video.ProxyUrl,
		Width:    int(video.Width),
		Height:   int(video.Height),
	}

	return &discordVideo
}

func mapProvider(provider *pb.EmbedProvider) *discordgo.MessageEmbedProvider {
	if provider == nil {
		return nil
	}

	discordProvider := discordgo.MessageEmbedProvider{
		URL:  provider.Url,
		Name: provider.Name,
	}

	return &discordProvider
}

func mapAuthor(author *pb.EmbedAuthor) *discordgo.MessageEmbedAuthor {
	if author == nil {
		return nil
	}

	discordAuthor := discordgo.MessageEmbedAuthor{
		URL:          author.Url,
		Name:         author.Name,
		IconURL:      author.IconUrl,
		ProxyIconURL: author.ProxyIconUrl,
	}

	return &discordAuthor
}

func mapFields(fields []*pb.EmbedField) []*discordgo.MessageEmbedField {
	discordFields := []*discordgo.MessageEmbedField{}

	for _, field := range fields {
		discordField := discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}

		discordFields = append(discordFields, &discordField)
	}

	return discordFields
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
