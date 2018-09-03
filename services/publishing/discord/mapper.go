package discord

import (
	"bytes"

	"github.com/bwmarrin/discordgo"

	pb "github.com/jukeizu/treediagram/api/publishing"
)

func MapToEmbed(embed *pb.Embed) *discordgo.MessageEmbed {
	discordEmbed := discordgo.MessageEmbed{}

	if embed == nil {
		return nil
	}

	discordEmbed.URL = embed.Url
	discordEmbed.Type = embed.Type
	discordEmbed.Title = embed.Title
	discordEmbed.Description = embed.Description
	discordEmbed.Timestamp = embed.Timestamp
	discordEmbed.Color = int(embed.Color)
	discordEmbed.Footer = MapFooter(embed.Footer)
	discordEmbed.Image = MapImage(embed.Image)
	discordEmbed.Thumbnail = MapThumbnail(embed.Thumbnail)
	discordEmbed.Video = MapVideo(embed.Video)
	discordEmbed.Provider = MapProvider(embed.Provider)
	discordEmbed.Author = MapAuthor(embed.Author)
	discordEmbed.Fields = MapFields(embed.Fields)

	return &discordEmbed
}

func MapToMessageSend(message *pb.Message) *discordgo.MessageSend {
	if message == nil {
		return nil
	}

	mappedEmbed := MapToEmbed(message.Embed)

	messageSend := discordgo.MessageSend{}
	messageSend.Content = message.Content
	messageSend.Embed = mappedEmbed

	if message.Files != nil {
		for _, file := range message.Files {
			discordFile := MapToFile(file)
			messageSend.Files = append(messageSend.Files, discordFile)
		}
	}

	return &messageSend
}

func MapToFile(file *pb.File) *discordgo.File {
	if file == nil {
		return nil
	}
	discordFile := discordgo.File{}
	discordFile.Name = file.Name
	discordFile.ContentType = file.ContentType
	discordFile.Reader = bytes.NewReader(file.Bytes)

	return &discordFile
}

func MapFooter(footer *pb.EmbedFooter) *discordgo.MessageEmbedFooter {
	if footer == nil {
		return nil
	}

	discordFooter := discordgo.MessageEmbedFooter{}

	discordFooter.Text = footer.Text
	discordFooter.IconURL = footer.IconUrl
	discordFooter.ProxyIconURL = footer.ProxyIconUrl

	return &discordFooter
}

func MapImage(image *pb.EmbedImage) *discordgo.MessageEmbedImage {
	if image == nil {
		return nil
	}

	discordImage := discordgo.MessageEmbedImage{}

	discordImage.URL = image.Url
	discordImage.ProxyURL = image.ProxyUrl
	discordImage.Width = int(image.Width)
	discordImage.Height = int(image.Height)

	return &discordImage
}

func MapThumbnail(thumbnail *pb.EmbedThumbnail) *discordgo.MessageEmbedThumbnail {
	if thumbnail == nil {
		return nil
	}

	discordThumbnail := discordgo.MessageEmbedThumbnail{}

	discordThumbnail.URL = thumbnail.Url
	discordThumbnail.ProxyURL = thumbnail.ProxyUrl
	discordThumbnail.Width = int(thumbnail.Width)
	discordThumbnail.Height = int(thumbnail.Height)

	return &discordThumbnail
}

func MapVideo(video *pb.EmbedVideo) *discordgo.MessageEmbedVideo {
	if video == nil {
		return nil
	}

	discordVideo := discordgo.MessageEmbedVideo{}

	discordVideo.URL = video.Url
	discordVideo.ProxyURL = video.ProxyUrl
	discordVideo.Width = int(video.Width)
	discordVideo.Height = int(video.Height)

	return &discordVideo
}

func MapProvider(provider *pb.EmbedProvider) *discordgo.MessageEmbedProvider {
	if provider == nil {
		return nil
	}

	discordProvider := discordgo.MessageEmbedProvider{}

	discordProvider.URL = provider.Url
	discordProvider.Name = provider.Name

	return &discordProvider
}

func MapAuthor(author *pb.EmbedAuthor) *discordgo.MessageEmbedAuthor {
	if author == nil {
		return nil
	}

	discordAuthor := discordgo.MessageEmbedAuthor{}

	discordAuthor.URL = author.Url
	discordAuthor.Name = author.Name
	discordAuthor.IconURL = author.IconUrl
	discordAuthor.ProxyIconURL = author.ProxyIconUrl

	return &discordAuthor
}

func MapFields(fields []*pb.EmbedField) []*discordgo.MessageEmbedField {
	discordFields := []*discordgo.MessageEmbedField{}

	for _, field := range fields {
		discordField := discordgo.MessageEmbedField{}

		discordField.Name = field.Name
		discordField.Value = field.Value
		discordField.Inline = field.Inline

		discordFields = append(discordFields, &discordField)
	}

	return discordFields
}
