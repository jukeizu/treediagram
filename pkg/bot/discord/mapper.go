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

func mapToPbServer(guild *discordgo.Guild) *processingpb.Server {
	return &processingpb.Server{
		Id:              guild.ID,
		Name:            guild.Name,
		OwnerId:         guild.OwnerID,
		Description:     guild.Description,
		UserCount:       int32(guild.MemberCount),
		IconUrl:         guild.IconURL(),
		SystemChannelId: guild.SystemChannelID,
	}
}

func mapToPbServers(userId string, guilds []*discordgo.Guild) []*processingpb.Server {
	servers := []*processingpb.Server{}

	for _, guild := range guilds {
		for _, member := range guild.Members {
			if member.User.ID == userId {
				server := mapToPbServer(guild)
				servers = append(servers, server)
				break
			}
		}
	}

	return servers
}

func mapToPbMessageRequest(state *discordgo.State, a *discordgo.Application, m *discordgo.Message) *processingpb.MessageRequest {
	return &processingpb.MessageRequest{
		Id:          m.ID,
		Source:      "discord",
		Bot:         mapToPbUser(state.User),
		Author:      mapToPbUser(m.Author),
		ChannelId:   m.ChannelID,
		ServerId:    m.GuildID,
		IsDirect:    m.GuildID == "",
		Servers:     mapToPbServers(m.Author.ID, state.Guilds),
		Content:     m.Content,
		Mentions:    mapToPbUsers(m.Mentions),
		Application: mapToPbApplication(a),
	}
}

func mapToPbReaction(state *discordgo.State, r *discordgo.MessageReaction, a *discordgo.Application, m *discordgo.Message) *processingpb.Reaction {
	return &processingpb.Reaction{
		UserId:         r.UserID,
		ChannelId:      r.ChannelID,
		ServerId:       r.GuildID,
		Emoji:          mapToPbEmoji(r.Emoji),
		MessageRequest: mapToPbMessageRequest(state, a, m),
	}
}

func mapToPbInteraction(state *discordgo.State, i *discordgo.InteractionCreate) *processingpb.Interaction {
	user := i.User
	if i.User == nil && i.Member != nil && i.Member.User != nil {
		user = i.Member.User
	}

	data := i.MessageComponentData()

	return &processingpb.Interaction{
		Identifier: data.CustomID,
		Values:     data.Values,
		Type:       data.Type().String(),
		Source:     "discord",
		Bot:        mapToPbUser(state.User),
		User:       mapToPbUser(user),
		ChannelId:  i.ChannelID,
		ServerId:   i.GuildID,
		IsDirect:   i.GuildID == "",
	}
}

func mapToPbEmoji(e discordgo.Emoji) *processingpb.Emoji {
	return &processingpb.Emoji{
		Id:            e.ID,
		Name:          e.Name,
		Roles:         e.Roles,
		Managed:       e.Managed,
		RequireColons: e.RequireColons,
		Animated:      e.Animated,
		Available:     e.Available,
	}
}

func mapToPbApplication(a *discordgo.Application) *processingpb.Application {
	return &processingpb.Application{
		Id:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Icon:        a.Icon,
		Owner:       mapToPbUser(a.Owner),
	}
}

func mapToMessageSend(contractMessage *contract.Message) (*discordgo.MessageSend, error) {
	if contractMessage == nil {
		return nil, nil
	}

	messageSend := discordgo.MessageSend{
		Content:    contractMessage.Content,
		Embeds:     mapToMessageEmbeds(contractMessage.Embed),
		Files:      mapToFiles(contractMessage.Files),
		Components: mapToMessageComponents(contractMessage.Compontents),
	}

	return &messageSend, nil
}

func mapToMessageEmbeds(contractEmbed *contract.Embed) []*discordgo.MessageEmbed {
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

	return []*discordgo.MessageEmbed{&embed}
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

func mapToMessageComponents(components *contract.Components) []discordgo.MessageComponent {
	if components == nil {
		return nil
	}

	messageComponents := []discordgo.MessageComponent{}

	for _, actionRow := range mapToActionRows(components.ActionsRows) {
		messageComponents = append(messageComponents, actionRow)
	}

	for _, button := range mapToButtons(components.Buttons) {
		messageComponents = append(messageComponents, button)
	}

	for _, textInput := range mapToTextInputs(components.TextInputs) {
		messageComponents = append(messageComponents, textInput)
	}

	for _, selectMenu := range mapToSelectMenus(components.SelectMenus) {
		messageComponents = append(messageComponents, selectMenu)
	}

	return messageComponents
}

func mapToActionRows(actionsRows []*contract.ActionsRow) []*discordgo.ActionsRow {
	if actionsRows == nil {
		return nil
	}

	dActionsRows := []*discordgo.ActionsRow{}

	for _, actionRow := range actionsRows {
		dActionRow := &discordgo.ActionsRow{}

		for _, button := range mapToButtons(actionRow.Buttons) {
			dActionRow.Components = append(dActionRow.Components, button)
		}

		for _, textInput := range mapToTextInputs(actionRow.TextInputs) {
			dActionRow.Components = append(dActionRow.Components, textInput)
		}

		for _, selectMenu := range mapToSelectMenus(actionRow.SelectMenus) {
			dActionRow.Components = append(dActionRow.Components, selectMenu)
		}

		dActionsRows = append(dActionsRows, dActionRow)
	}

	return dActionsRows
}

func mapToComponentEmoji(emoji *contract.ComponentEmoji) discordgo.ComponentEmoji {
	if emoji == nil {
		return discordgo.ComponentEmoji{}
	}

	return discordgo.ComponentEmoji{
		ID:       emoji.Id,
		Name:     emoji.Name,
		Animated: emoji.Animated,
	}
}

func mapToButtons(buttons []*contract.Button) []*discordgo.Button {
	if buttons == nil {
		return nil
	}

	dButtons := []*discordgo.Button{}

	for _, button := range buttons {
		dButton := &discordgo.Button{
			Label:    button.Label,
			Style:    discordgo.ButtonStyle(button.Style),
			Emoji:    mapToComponentEmoji(&button.Emoji),
			Disabled: button.Disabled,
			URL:      button.Url,
			CustomID: button.CustomId,
		}
		dButtons = append(dButtons, dButton)
	}

	return dButtons
}

func mapToTextInputs(textInputs []*contract.TextInput) []*discordgo.TextInput {
	if textInputs == nil {
		return nil
	}

	dTextInputs := []*discordgo.TextInput{}

	for _, textInput := range textInputs {
		dButton := &discordgo.TextInput{
			CustomID:    textInput.CustomId,
			Label:       textInput.Label,
			Style:       discordgo.TextInputStyle(textInput.Style),
			Placeholder: textInput.Placeholder,
			Value:       textInput.Value,
			Required:    textInput.Required,
			MinLength:   textInput.MinLength,
			MaxLength:   textInput.MaxLength,
		}
		dTextInputs = append(dTextInputs, dButton)
	}

	return dTextInputs
}

func mapToSelectMenuOptions(selectMenuOptions []contract.SelectMenuOption) []discordgo.SelectMenuOption {
	dSelectMenuOptions := []discordgo.SelectMenuOption{}

	for _, selectMenuOption := range selectMenuOptions {
		dSelectMenuOption := discordgo.SelectMenuOption{
			Label:       selectMenuOption.Label,
			Value:       selectMenuOption.Value,
			Description: selectMenuOption.Description,
			Emoji:       mapToComponentEmoji(&selectMenuOption.Emoji),
			Default:     selectMenuOption.Default,
		}
		dSelectMenuOptions = append(dSelectMenuOptions, dSelectMenuOption)
	}

	return dSelectMenuOptions
}

func mapToSelectMenus(selectMenus []*contract.SelectMenu) []*discordgo.SelectMenu {
	if selectMenus == nil {
		return nil
	}

	dSelectMenus := []*discordgo.SelectMenu{}

	for _, selectMenu := range selectMenus {
		dSelectMenu := &discordgo.SelectMenu{
			CustomID:    selectMenu.CustomId,
			Placeholder: selectMenu.Placeholder,
			MinValues:   selectMenu.MinValues,
			MaxValues:   selectMenu.MaxValues,
			Options:     mapToSelectMenuOptions(selectMenu.Options),
			Disabled:    selectMenu.Disabled,
		}
		dSelectMenus = append(dSelectMenus, dSelectMenu)
	}

	return dSelectMenus
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
