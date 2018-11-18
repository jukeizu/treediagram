package receiver

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/receiving"
	"github.com/rs/xid"
)

type Message struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Bot       User   `json:"bot"`
	Author    User   `json:"author"`
	ChannelId string `json:"channelId"`
	ServerId  string `json:"serverId"`
	Mentions  []User `json:"mentions"`
	Content   string `json:"content"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewMessage(req *receiving.Message) Message {
	message := Message{
		Id:        xid.New().String(),
		Source:    req.Source,
		Bot:       NewUser(req.Bot),
		Author:    NewUser(req.Author),
		ChannelId: req.ChannelId,
		ServerId:  req.ServerId,
		Content:   req.Content,
	}

	for _, mention := range req.Mentions {
		message.Mentions = append(message.Mentions, NewUser(mention))
	}

	return message
}

func NewUser(receivingUser *receiving.User) User {
	user := User{
		Id:   receivingUser.Id,
		Name: receivingUser.Name,
	}

	return user
}
